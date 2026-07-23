package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/teris-io/shortid"
)

var token = os.Getenv("DROP_TOKEN") // cant put := between imports and main

func main() {
    initDB()
    log.Println("db intialized")
    os.MkdirAll("./data", 0755) // creates folder and sets perms
    http.HandleFunc("/upload", upload) // handlers
    http.HandleFunc("/file/", download)
    http.HandleFunc("/d/", showPage) // the download page
    log.Println("server running! on port 5100")
    log.Fatal(http.ListenAndServe(":5100", nil))
}

func upload(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Authorization") != "Bearer "+token {
        http.Error(w, "unauthorized", 401)
        return
    } // checks headers and token
    r.ParseMultipartForm(500 << 20) // makes sure file size isnt above
    file, info, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "no file", 400) // file header not set
        return
    }
    defer file.Close()
    id, err := shortid.Generate() // gen id
    if err != nil {
        http.Error(w, "id error", 500)
        return
    }
    path := filepath.Join("./data", id)
    out, err := os.Create(path) // create the file
    if err != nil {
        http.Error(w, "save error", 500)
        return
    }
    defer out.Close()
    size, err := io.Copy(out, file)
    if err != nil {
        http.Error(w, "upload error", 500)
        return
    }
    saveFileDB(File{
        ID:   id,
        Name: info.Filename,
        Path: path,
        Size: size,
    }) // feeds to my func the data
    fmt.Fprintf(w, `{
        "id": "%s",
        "file": "%s",
        "url": "/d/%s"
    }`, id, info.Filename, id) // prints
}

func download(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/file/")
    file, err := getFileDB(id)
    if err != nil {
        http.Error(w, "not found", 404)
        return
    }
    f, err := os.Open(file.Path)
    if err != nil {
        http.Error(w, "missing file", 404)
        return
    }
    defer f.Close()
    w.Header().Set(
        "Content-Disposition",
        fmt.Sprintf(`attachment; filename="%s"`, file.Name),
    )
    io.Copy(w, f)
}

func showPage(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/d/")
    file, err := getFileDB(id)
    if err != nil {
        http.Error(w, "not found", 404)
        return
    }
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Downloading %s</title></head>
<body>
    <p>Downloading %s...</p>
    <script>window.location = "/file/%s";</script>
    <noscript><a href="/file/%s">Click here to download</a></noscript>
</body>
</html>`, file.Name, file.Name, id, id)
}
