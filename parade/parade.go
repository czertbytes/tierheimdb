package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func main() {
	if err := pb.RedisInit(); err != nil {
		log.Fatalf("Redis init failed! Error: %s\n", err)
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", GetHomeHandler).Methods("GET")
	router.HandleFunc("/contact", GetContactHandler).Methods("GET")
	router.HandleFunc("/{shelterId}", GetShelterHandler).Methods("GET")
	router.HandleFunc("/{shelterId}/{updateId}/{animalId}", GetAnimalHandler).Methods("GET")

	http.Handle("/", router)

	serveFile("/sitemap.xml", "./sitemap.xml")
	serveFile("/favicon.ico", "./favicon.ico")
	serveFile("/robots.txt", "./robots.txt")
	serveFile("/humans.txt", "./humans.txt")
	serveFile("/pure-min.css", "./pure-min.css")
	serveFile("/tierheimdb.css", "./tierheimdb.css")
	serveFile("/bmt-logo.png", "./bmt-logo.png")
	serveFile("/dresden-logo.gif", "./dresden-logo.gif")
	serveFile("/samtpfoten-neukoelln-logo.png", "./samtpfoten-neukoelln-logo.png")
	serveFile("/tierheim-muenchen-logo.png", "./tierheim-muenchen-logo.png")
	serveFile("/tierheim-berlin-logo.jpg", "./tierheim-berlin-logo.jpg")

	tdbRoot := os.Getenv("GOPATH")
	if len(tdbRoot) == 0 {
		log.Fatalf("Environment variable GOPATH not set!")
	}
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir(fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/s", tdbRoot)))))

	log.Println("Running Parade")
	log.Fatalf("Failed to run webserver: %s", http.ListenAndServe(":8081", nil))
}
