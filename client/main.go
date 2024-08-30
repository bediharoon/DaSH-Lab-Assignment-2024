package main

import (
    "bufio"
    "bytes"
    "encoding/json"
    "hash/adler32"
    "io"
    "log"
    "os"
    "time"

    "github.com/bediharoon/dashassign/datamodel"

    "github.com/godbus/dbus/v5"
)

func hashGenerator(x string) (uint32) {
    return adler32.Checksum([]byte(x))
}

func retrieveInputs(fd *os.File) ([]string, []uint32, error) {
    var prompts []string
    var promptHashes []uint32

    scanner := bufio.NewScanner(fd)
    for scanner.Scan() {
        prompts = append(prompts, scanner.Text()) // Need for more than 64K chars per prompt ignored

        hash := hashGenerator(scanner.Text()) // Uniqueness Ignored in Academic Context
        
        promptHashes = append(promptHashes, hash)
    }

    if err := scanner.Err(); err != nil {
        return nil, nil, err
    }

    return prompts, promptHashes, nil
}

func writeToFile(data *datamodel.RequestData, file *os.File) (error) {
    buf := bytes.NewBuffer(nil)
    _, err := io.Copy(buf, file)
    if err != nil {
        return err
    }

    var dataIn []datamodel.RequestData
    
    if string(buf.Bytes()) != "" {
        err = json.Unmarshal(buf.Bytes(), &dataIn)
        if err != nil {
            return err
        }
    }

    dataIn = append(dataIn, *data)

    dataJSON, err := json.MarshalIndent(dataIn, "", "")
    if err != nil {
        return err
    }

    file.Truncate(0)
    file.Seek(0, 0)

    _, err = file.Write(dataJSON)
    if err != nil {
        return err
    }

    return nil
}

func main() {
    conn, err := dbus.ConnectSessionBus()
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    obj := conn.Object("com.github.bediharoon.ServLLM", "/com/github/bediharoon/ServLLM")


    if len(os.Args) < 3 {
        log.Fatalln("Correct Usage: $ ./client /path/to/input/file /path/to/output/file")
    }

    inputFilePath := os.Args[1]
    outputFilePath := os.Args[2]


    inputFile, err := os.Open(inputFilePath)
    if err != nil {
        log.Fatal(err)
    }
    defer inputFile.Close()

    inputPrompts, promptHashes, err := retrieveInputs(inputFile)
    if err != nil {
        log.Fatal(err)
    }

    var requiredPromptHashes []uint32
    var uid int64
    err = obj.Call("com.github.bediharoon.ServLLM.RegisterHashes", 0, promptHashes).Store(&requiredPromptHashes, &uid)
    if (err != nil) {
        log.Fatal(err)
    }

    var inputFormatted []datamodel.RequestData

    var required bool
    for _, e := range inputPrompts {
        required = false
        hash := hashGenerator(e)
        if err != nil {
            log.Println("Error Generating Hash, Omitting: ", err)
            hash = 0
        }

        for _, f := range requiredPromptHashes {
            if hash == f {
                required = true
            }
        }

        if !required {
            continue
        }

        inputFormatted = append(inputFormatted, datamodel.RequestData {
            Prompt: e,
            TimeSent: time.Now().Unix(),
            InitiatorUID: uid,
        })
    }

    go func() {
        err = obj.Call("com.github.bediharoon.ServLLM.PromptRequest", dbus.FlagNoReplyExpected, inputFormatted).Err
        if err != nil {
            log.Fatal(err)
        }
    }()

    if err = conn.AddMatchSignal(
        dbus.WithMatchObjectPath("/com/github/bediharoon/ServLLM"),
        dbus.WithMatchInterface("com.github.bediharoon.ServLLM"),
        dbus.WithMatchSender("com.github.bediharoon.ServLLM"),
    ); err != nil {
        log.Fatal(err)
    }


    c := make(chan *dbus.Signal, 12)
    conn.Signal(c)

    for x := range c {
        var resp datamodel.RequestData
        err = dbus.Store(x.Body, &resp)
        if err != nil {
            log.Fatal(err)
        }

        if resp.ClientID != uid {
            continue
        }

        resp.TimeRecvd = time.Now().Unix()

        if resp.ClientID != resp.InitiatorUID {
            resp.Source = "user"
        }
 
        outputFile, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_SYNC|os.O_CREATE, 0600)
        if err != nil {
            obj.Call("com.github.bediharoon.ServLLM.WriteFail", 0, resp.ClientID)
            log.Fatal(err)
        }
        defer outputFile.Close()

        err = writeToFile(&resp, outputFile)
        if err != nil {
            obj.Call("com.github.bediharoon.ServLLM.WriteFail", 0, resp.ClientID)
            log.Fatal(err)
        }

        err = obj.Call("com.github.bediharoon.ServLLM.WriteCheck", 0, resp.ClientID).Err
        if err != nil {
            log.Fatal(err)
        }
    }
}
