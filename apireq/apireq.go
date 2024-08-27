package apireq

import (
    "net/http"
    "io"
)

var (
    APIUrl = "https://api-inference.huggingface.co/models/google/gemma-2b"
)

func LLMRequest(payload io.Reader, requestToken string) (string, error) {
    req, err := http.NewRequest("POST", APIUrl, payload)
    if err != nil {
        return "", err
    }

    req.Header.Set("Authorization", "Bearer " + requestToken)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }

    result, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(result), nil
}
