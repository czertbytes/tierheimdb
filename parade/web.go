package main

import (
	"html/template"
	"log"
	"net/http"

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
	Cats      int
	Dogs      int
}

type ShelterPage struct {
	Shelters      []*Shelter
	Shelter       Shelter
	SheltersTotal int
}

type AnimalPage struct {
	Shelters      []*Shelter
	Shelter       Shelter
	SheltersTotal int
	Animal        pb.Animal
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
}

func GetShelterHandler(w http.ResponseWriter, r *http.Request) {
	shelterId := mux.Vars(r)["shelterId"]
	animalType := r.URL.Query().Get("type")

	shelters, err := enabledShelters(shelterId)
	if err != nil {
		log.Println(err)
		return
	}

	var shelter Shelter
	for _, s := range shelters {
		if s.PBShelter.Id == shelterId {
			s.Selected = true
			shelter = *s
		}
	}

	if len(animalType) > 0 {
		animals := []pb.Animal{}
		for _, a := range shelter.PBAnimals {
			if a.Type == animalType {
				animals = append(animals, a)
			}
		}
		shelter.PBAnimals = animals
	}

	for _, s := range shelters {
		log.Println(s.Selected)
	}

	tmpl.ExecuteTemplate(w, "shelter", &ShelterPage{
		shelters,
		shelter,
		len(shelters),
	})
}

func enabledShelters(shelterId string) ([]*Shelter, error) {
	pbShelters, err := pb.GetEnabledShelters()
	if err != nil {
		return nil, err
	}

	shelters := []*Shelter{}
	for _, pbShelter := range pbShelters {
		shelter, err := makeShelter(pbShelter)
		if err != nil {
			return nil, err
		}

		shelters = append(shelters, &shelter)
	}

	return shelters, nil
}

func makeShelter(shelter pb.Shelter) (Shelter, error) {
	update, err := pb.GetLastUpdate(shelter.Id)
	if err != nil {
		return Shelter{}, err
	}

	animals, err := pb.GetAnimals(shelter.Id, update.Id)
	if err != nil {
		return Shelter{}, err
	}

	cats, dogs := 0, 0
	for _, animal := range animals {
		if animal.Type == "cat" {
			cats = cats + 1
		}
		if animal.Type == "dog" {
			dogs = dogs + 1
		}
	}

	return Shelter{
		shelter,
		false,
		animals,
		update,
		cats,
		dogs,
	}, nil
}

func GetAnimalHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	shelterId := params["shelterId"]
	updateId := params["updateId"]
	animalId := params["animalId"]

	shelters, err := enabledShelters(shelterId)
	if err != nil {
		log.Println(err)
		return
	}

	var shelter Shelter
	for _, s := range shelters {
		if s.PBShelter.Id == shelterId {
			s.Selected = true
			shelter = *s
		}
	}

	animal, err := pb.GetAnimal(shelterId, updateId, animalId)
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

func GetHelpHandler(w http.ResponseWriter, r *http.Request) {
}

func GetAboutHandler(w http.ResponseWriter, r *http.Request) {
}
