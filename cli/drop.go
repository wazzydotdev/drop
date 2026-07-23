package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/briandowns/spinner"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: drop <file>")
        return
    }
    fileToDrop := os.Args[1]

    os.MkdirAll(".config/drop", 0755)

    keyUntrimmed, err := os.ReadFile(".config/drop/key.txt")
    if err != nil {
        fmt.Println("no key.txt found in .config/drop, add one")
        return
    }
    key := strings.TrimSpace(string(keyUntrimmed))

    serverUntrimmed, err := os.ReadFile(".config/drop/server.txt")
    if err != nil {
        fmt.Println("no server.txt found in .config/drop, add one")
        return
    }
    srv := strings.TrimSpace(string(serverUntrimmed))

    s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
    s.Suffix = "  uploading " + filepath.Base(fileToDrop)
    s.Start()

    f, _ := os.Open(fileToDrop)
    defer f.Close()

    var buf bytes.Buffer
    w := multipart.NewWriter(&buf)
    part, _ := w.CreateFormFile("file", filepath.Base(fileToDrop))
    io.Copy(part, f)
    w.Close()

    req, _ := http.NewRequest("POST", srv+"/upload", &buf)
    req.Header.Set("Content-Type", w.FormDataContentType())
    req.Header.Set("Authorization", "Bearer "+key)

    resp, err := http.DefaultClient.Do(req)
    s.Stop()
    if err != nil {
        fmt.Println("upload failed:", err)
        return
    }
    defer resp.Body.Close()
    out, _ := io.ReadAll(resp.Body)

    res := map[string]string{}
    json.Unmarshal(out, &res)
    fmt.Printf("\033[32mdone\033[0m -> %s\n\033[36murl:\033[0m %s%s\n", res["file"], srv, res["url"])
}
