package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
	tmpl *template.Template
)

type Shelter struct {
	PBShelter pb.Shelter
	Selected  bool
	PBAnimals []pb.Animal
	PBUpdate  pb.Update
}

type Shelters []*Shelter

type IndexPage struct {
	Shelters Shelters
}

type ShelterPage struct {
	Shelters      Shelters
	Shelter       Shelter
	SheltersTotal int
}

type AnimalPage struct {
	Shelters      Shelters
	Shelter       Shelter
	SheltersTotal int
	Animal        pb.Animal
}

type ContactPage struct {
}

func init() {
	tdbRoot := os.Getenv("GOPATH")
	if len(tdbRoot) == 0 {
		log.Fatalf("Environment variable GOPATH not set!")
	}

	files := []string{}
	for _, f := range []string{"index.html", "shelter.html", "animal.html", "contact.html"} {
		files = append(files, fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/tmpl/%s", tdbRoot, f))
	}

	var err error
	tmpl, err = tmpl.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
}

func GetIndexHandler(w http.ResponseWriter, r *http.Request) {
	shelters, err := makeShelters()
	if err != nil {
		log.Println(err)
		return
	}

	tmpl.ExecuteTemplate(w, "index", &IndexPage{
		shelters,
	})
}

func GetShelterHandler(w http.ResponseWriter, r *http.Request) {
	shelterId := mux.Vars(r)["shelterId"]
	animalType := r.URL.Query().Get("type")

	shelters, shelter, err := makeSheltersShelter(shelterId, animalType)
	if err != nil {
		log.Println(err)
		return
	}

	tmpl.ExecuteTemplate(w, "shelter", &ShelterPage{
		shelters,
		shelter,
		len(shelters),
	})
}

func GetAnimalHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shelterId := params["shelterId"]
	updateId := params["updateId"]
	animalId := params["animalId"]

	shelters, shelter, animal, err := makeSheltersShelterAnimal(shelterId, updateId, animalId)
	if err != nil {
		log.Println(err)
		return
	}

	tmpl.ExecuteTemplate(w, "animal", &AnimalPage{
		shelters,
		shelter,
		len(shelters),
		animal,
	})
}

func GetContactHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "contact", &ContactPage{})
}
