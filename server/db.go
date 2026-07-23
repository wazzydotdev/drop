package main

import (
    "database/sql"

    _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type File struct {
    ID   string
    Name string
    Path string
    Size int64
}

func initDB() error {
    var err error

    db, err = sql.Open("sqlite3", "files.db")
    if err != nil {
        return err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS files (
            id TEXT PRIMARY KEY,
            name TEXT,
            path TEXT,
            size INTEGER
        )
    `)

    return err
}


func saveFileDB(file File) error {
    _, err := db.Exec(
        `
        INSERT INTO files (
            id,
            name,
            path,
            size
        )
        VALUES (?, ?, ?, ?)
        `,
        file.ID,
        file.Name,
        file.Path,
        file.Size,
    )

    return err
}


func getFileDB(id string) (File, error) {
    var file File

    err := db.QueryRow(
        `
        SELECT id, name, path, size
        FROM files
        WHERE id = ?
        `,
        id,
    ).Scan(
        &file.ID,
        &file.Name,
        &file.Path,
        &file.Size,
    )

    return file, err
}
