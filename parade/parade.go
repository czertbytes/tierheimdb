package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/gorilla/mux"

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

	//router := mux.NewRouter()
	serveFile("/favicon.ico", "favicon.ico")
	serveFile("/robots.txt", "robots.txt")
	serveFile("/humans.txt", "humans.txt")
	http.HandleFunc("/sitemap.xml", GetSitemapHandler)

	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir(fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/s", tdbRoot)))))
	http.Handle("/partials/", http.StripPrefix("/partials/", http.FileServer(http.Dir(fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/partials", tdbRoot)))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/partials/index.html", tdbRoot))
	})

	log.Println("Running Parade")
	log.Fatalf("Failed to run webserver: %s", http.ListenAndServe(":8081", nil))
}
