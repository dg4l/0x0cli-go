package main

import (
    "bytes"
    "fmt"
    "os"
    "net/http"
    "io"
    "mime/multipart"
    "path/filepath"
    "flag"
)

type RequestParams struct {
    path *string
    expires *string
    url *string
    secret *string
}

// this is not reusable in other contexts.
// doesn't matter though.
func create0x0UploadRequest(uri string, params RequestParams) (*http.Request, error){
    requestBody := &bytes.Buffer{}
    writer := multipart.NewWriter(requestBody)
    if *params.path != "" {
        file, err := os.Open(*params.path)
        if err != nil {
            return nil, err
        }
        defer file.Close()
        part, err := writer.CreateFormFile("file", filepath.Base(*params.path))
        if err != nil {
            return nil, err
        }
        _, err = io.Copy(part, file)
        if err != nil {
            return nil, err
        }
    }
    if *params.url != "" {
        writer.WriteField("url", *params.url)
    }
    if *params.secret != "" {
        writer.WriteField("secret", *params.secret)
    }
    if *params.expires != "" {
        writer.WriteField("expires", *params.expires)
    }
    err := writer.Close()
    if err != nil {
        return nil, err
    }
    req, err := http.NewRequest("POST", uri, requestBody)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    return req, nil
}

func main() {
    params := RequestParams{}
    server_url := "https://0x0.st"
    params.path = flag.String("f", "", "path")
    params.expires = flag.String("e", "", "expiry date")
    params.url = flag.String("u", "", "url")
    params.secret = flag.String("s", "", "secret")
    flag.Parse()
    if *params.path == "" && *params.url == "" {
        fmt.Println("must either provide a filepath with -f or a url to a file with -u")
        return   
    }
    request, err := create0x0UploadRequest(server_url, params)
    if err != nil {
        fmt.Println(err)
        return
    }
    request.Header.Set("User-Agent", "0x0cli-go/1.0")
    client := http.Client{}
    response, err := client.Do(request)
    if err != nil{
        fmt.Println(err)
        return
    }
    body := bytes.Buffer{}
    _, err = body.ReadFrom(response.Body)
    if err != nil {
        fmt.Println(err)
        return
    }
    response.Body.Close()
    fmt.Println(response.StatusCode)
    fmt.Printf("%s", body.String())
}
