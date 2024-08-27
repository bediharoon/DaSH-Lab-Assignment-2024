package main

import (
    "bytes"
    "net/http"
    "encoding/json"
    "io"
    "log"
    "os"
)

func main() {
    API_URL := "https://api-inference.huggingface.co/models/google/gemma-2b"
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
    payload["inputs"] = string(inputBuffer[:n-1])

    payloadJSON, err := json.Marshal(payload)
    if err != nil {
        log.Fatal(err)
    }
    
    req, err := http.NewRequest("POST", API_URL, bytes.NewReader(payloadJSON))
    log.Println(string(payloadJSON))
    if err != nil {
        log.Fatal(err)
    }

    token, exists := os.LookupEnv("GEMMA_TOKEN")
    if !exists {
        log.Fatal(`Please provide a User Access Token (May be obtained from Hugging Face)
        $ export GEMMA_TOKEN="<TOKEN HERE>"`)
    }
    req.Header.Set("Authorization", "Bearer " + token)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}

    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    respData, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

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

    _, err = fo.Write(respData)
        if err != nil {
            log.Fatal(err)
        }
    
}
