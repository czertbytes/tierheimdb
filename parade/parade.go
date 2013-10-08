package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
	tdbRoot string
)

func init() {
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

	router := mux.NewRouter()
	router.HandleFunc("/", GetHomeHandler).Methods("GET")
	serveFile("/favicon.ico", "favicon.ico")
	serveFile("/robots.txt", "robots.txt")
	serveFile("/humans.txt", "humans.txt")
	router.HandleFunc("/sitemap.xml", GetSitemapHandler).Methods("GET")

	router.HandleFunc("/about", GetAboutHandler).Methods("GET")
	router.HandleFunc("/{shelterId}", GetShelterHandler).Methods("GET")
	router.HandleFunc("/{shelterId}/{updateId}/{animalId}", GetAnimalHandler).Methods("GET")

	http.Handle("/", router)
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir(fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/s", tdbRoot)))))

	log.Println("Running Parade")
	log.Fatalf("Failed to run webserver: %s", http.ListenAndServe(":8081", nil))
}
