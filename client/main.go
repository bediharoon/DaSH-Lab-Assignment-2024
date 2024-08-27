package main

import (
    "encoding/json"
    "io"
    "log"
    "os"

    "github.com/godbus/dbus/v5"
)

func main() {
    conn, err := dbus.ConnectSessionBus()
	if err != nil {
        log.Fatal(err)
	}
	defer conn.Close()

    base := conn.Object("com.github.bediharoon.ServLLM", "/com/github/bediharoon/servllm")

    // Read from the input file
    fi, err := os.Open("input.txt")
    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        err := fi.Close()
        if err != nil {
            log.Fatal(err)
        }
    }()

    inputBuffer := make([]byte, 4096)

    n, err := fi.Read(inputBuffer)
    if err != nil && err != io.EOF {
        log.Fatal(err)
    }


    payload := make(map[string]string)
    payload["inputs"] =  string(inputBuffer[:n-1])

    payloadJSON, err := json.Marshal(payload)
    if err != nil {
        log.Fatal(err)
    }

    output := base.Call("com.github.bediharoon.ServLLM.GetResp", 0, payloadJSON, "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").Body
    outJSON, err := json.Marshal(output[0].(string))

    // write output to a file
    fo, err := os.Create("output.json")
    if err != nil {
        log.Fatal(err)
    }

    defer func() {
        err := fo.Close()
        if err != nil {
            log.Fatal(err)
        }
    }()

    _, err = fo.Write(outJSON)
        if err != nil {
            log.Fatal(err)
        }
}
