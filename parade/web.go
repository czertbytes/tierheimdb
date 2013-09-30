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

type HomePage struct {
	Title    string
	Shelters Shelters
}

type ShelterPage struct {
	Title    string
	Shelters Shelters
	Shelter  Shelter
}

type AnimalPage struct {
	Title    string
	Shelters Shelters
	Shelter  Shelter
	Animal   pb.Animal
}

type ContactPage struct {
	Title string
}

func init() {
	tdbRoot := os.Getenv("GOPATH")
	if len(tdbRoot) == 0 {
		log.Fatalf("Environment variable GOPATH not set!")
	}

	tmpl = template.Must(
		template.New("all").
			Funcs(template.FuncMap{"dateFormat": DateFormatter}).
			ParseGlob(fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/tmpl/*.html.tmpl", tdbRoot)))
}

func GetHomeHandler(w http.ResponseWriter, r *http.Request) {
	shelters, err := makeShelters()
	if err != nil {
		log.Println(err)
		return
	}

	tmpl.ExecuteTemplate(w, "home", &HomePage{
		"Home",
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
		shelter.PBShelter.Name,
		shelters,
		shelter,
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
		animal.Name,
		shelters,
		shelter,
		animal,
	})
}

func GetContactHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "contact", &ContactPage{
		"Contact",
	})
}
