package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
    tdbRoot string
)

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    tdbRoot = os.Getenv("GOPATH")
    if len(tdbRoot) == 0 {
        log.Fatalf("Environment variable GOPATH not set!")
    }
}

func main() {
    if err := pb.RedisInit(); err != nil {
        log.Fatalf("Redis init failed! Error: %s\n", err)
        return
    }

    http.HandleFunc("/sitemap.xml", GetSitemapHandler)

    serveFile("/favicon.ico", "favicon.ico")
    serveFile("/robots.txt", "robots.txt")
    serveFile("/humans.txt", "humans.txt")
    serveFile("/styles.min.css", "styles.min.css")
    serveFile("/scripts.min.js", "scripts.min.js")

    for _, path := range []string{"/views/", "/fonts/"} {
        http.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/public%s", tdbRoot, path)))))
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/public/index.html", tdbRoot))
    })

    log.Println("Running Parade")
    log.Fatalf("Failed to run webserver: %s", http.ListenAndServe(":8081", nil))
}
