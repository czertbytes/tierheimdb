package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
	tmpl *template.Template
)

type IndexPage struct {
	Shelters      []pb.Shelter
	SheltersTotal int
}

type ShelterPage struct {
	Shelters []pb.Shelter
	Shelter  pb.Shelter
	Animals  []pb.Animal
}

type AnimalPage struct {
	Shelters []pb.Shelter
	Shelter  pb.Shelter
	Animal   pb.Animal
}

func init() {
	files := []string{
		"tmpl/index.html",
		"tmpl/shelter.html",
		"tmpl/animal.html",
	}

	var err error
	tmpl, err = tmpl.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
}

func GetIndexHandler(w http.ResponseWriter, r *http.Request) {
	shelters, err := pb.GetEnabledShelters()
	if err != nil {
		return
	}
	sheltersTotal := len(shelters)

	tmpl.ExecuteTemplate(w, "index", &IndexPage{
		shelters,
		sheltersTotal,
	})
}

func GetShelterHandler(w http.ResponseWriter, r *http.Request) {
	shelterId := mux.Vars(r)["shelterId"]

	shelters, err := pb.GetEnabledShelters()
	if err != nil {
		return
	}

	shelter, err := pb.GetShelter(shelterId)
	if err != nil {
		return
	}

	update, err := pb.GetLastUpdate(shelterId)
	if err != nil {
		return
	}

	animals, err := pb.GetAnimals(shelterId, update.Id)
	if err != nil {
		return
	}

	tmpl.ExecuteTemplate(w, "shelter", &ShelterPage{
		shelters,
		shelter,
		animals,
	})
}

func GetShelterUpdateHandler(w http.ResponseWriter, r *http.Request) {
}

func GetShelterUpdateAnimalHandler(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)
	shelterId := routeParams["shelterId"]
	updateId := routeParams["updateId"]
	animalId := routeParams["animalId"]

	shelters, err := pb.GetEnabledShelters()
	if err != nil {
		return
	}

	shelter, err := pb.GetShelter(shelterId)
	if err != nil {
		return
	}

	animal, err := pb.GetAnimal(shelterId, updateId, animalId)
	if err != nil {
		return
	}

	tmpl.ExecuteTemplate(w, "animal", &AnimalPage{
		shelters,
		shelter,
		animal,
	})
}

func GetHelpHandler(w http.ResponseWriter, r *http.Request) {
}

func GetAboutHandler(w http.ResponseWriter, r *http.Request) {
}
