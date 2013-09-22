package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func main() {
	if err := pb.RedisInit(); err != nil {
		log.Fatalf("Redis init failed! Error: %s\n", err)
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/", GetIndexHandler).Methods("GET")
	router.HandleFunc("/{shelterId}", GetShelterHandler).Methods("GET")
	router.HandleFunc("/{shelterId}/{updateId}/{animalId}", GetAnimalHandler).Methods("GET")
	router.HandleFunc("/help", GetHelpHandler).Methods("GET")
	router.HandleFunc("/about", GetAboutHandler).Methods("GET")

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

	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("./s/"))))

	log.Println("Running Parade")
	log.Fatalf("Failed to run webserver: %s", http.ListenAndServe(":8081", nil))
}
