package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var v1Router *mux.Router

func main() {
	if err := pb.RedisInit(); err != nil {
		log.Fatalf("Redis init failed! Error: %s\n", err)
		return
	}

	filterRouter := new(FilterRouter)
	router := mux.NewRouter()
	v1Router = router.PathPrefix("/v1").Subrouter()

	v1Router.HandleFunc("/animals", APIv1GetAnimalsHandler).Methods("GET")
	v1Router.HandleFunc("/shelters", APIv1GetSheltersHandler).Methods("GET")
	v1Router.HandleFunc("/shelters", APIv1PostSheltersHandler).Methods("POST")
	v1Router.HandleFunc("/shelters", APIv1DeleteSheltersHandler).Methods("DELETE")
	v1Router.HandleFunc("/shelter/{shelterId}", APIv1GetShelterHandler).Methods("GET")
	v1Router.HandleFunc("/shelter/{shelterId}", APIv1DeleteShelterHandler).Methods("DELETE")
	v1Router.HandleFunc("/shelter/{shelterId}/animals", APIv1GetShelterAnimalsHandler).Methods("GET")
	v1Router.HandleFunc("/shelter/{shelterId}/animals", APIv1PostShelterAnimalsHandler).Methods("POST")
	v1Router.HandleFunc("/shelter/{shelterId}/sources", APIv1SyncShelterSourcesHandler).Methods("SYNC")
	v1Router.HandleFunc("/shelter/{shelterId}/updates", APIv1GetShelterUpdatesHandler).Methods("GET")
	v1Router.HandleFunc("/shelter/{shelterId}/updates", APIv1DeleteShelterUpdatesHandler).Methods("DELETE")
	v1Router.HandleFunc("/shelter/{shelterId}/update/{updateId}", APIv1GetShelterUpdateHandler).Methods("GET")
	v1Router.HandleFunc("/shelter/{shelterId}/update/{updateId}", APIv1DeleteShelterUpdateHandler).Methods("DELETE")
	v1Router.HandleFunc("/shelter/{shelterId}/update/{updateId}/animals", APIv1GetShelterUpdateAnimalsHandler).Methods("GET")
	v1Router.HandleFunc("/shelter/{shelterId}/update/{updateId}/animals", APIv1DeleteShelterUpdateAnimalsHandler).Methods("DELETE")
	v1Router.HandleFunc("/shelter/{shelterId}/update/{updateId}/animals/diff/{otherUpdateId}", APIv1GetShelterUpdatesDiffHandler).Methods("GET")
	v1Router.HandleFunc("/shelter/{shelterId}/update/{updateId}/animal/{animalId}", APIv1GetShelterUpdateAnimalHandler).Methods("GET")
	v1Router.HandleFunc("/shelter/{shelterId}/animal/{updateId}/animal/{animalId}", APIv1DeleteShelterUpdateAnimalHandler).Methods("DELETE")

	http.Handle("/", router)

	log.Println("Running Kennel")
	log.Fatalf("Failed to run webserver: %s", http.ListenAndServe(":8080", filterRouter))
}
