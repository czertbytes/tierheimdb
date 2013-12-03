package main

import (
    "fmt"
    "net/http"
)

func serveFile(pattern string, filename string) {
    http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
        path := fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/public/%s", tdbRoot, filename)
        http.ServeFile(w, r, path)
    })
}
