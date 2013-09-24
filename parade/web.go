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
	tdbRoot := os.Getenv("GOPATH")
	if len(tdbRoot) == 0 {
		log.Fatalf("Environment variable GOPATH not set!")
	}

	files := []string{}
	for _, f := range []string{"shelter.html", "animal.html"} {
		files = append(files, fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/tmpl/%s", tdbRoot, f))
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
			allAnimals, err := pb.GetAnimals(s.PBShelter.Id, s.PBUpdate.Id)
			if err != nil {
				log.Println(err)
				return
			}

			animals := []pb.Animal{}
			for _, a := range allAnimals {
				if len(a.Images) == 0 {
					a.Images = []pb.Image{pb.Image{URL: fmt.Sprintf("http://placehold.it/200x200&text=%s", a.Name)}}
				}

				if len(animalType) > 0 {
					if a.Type == animalType {
						animals = append(animals, a)
					}
				} else {
					animals = append(animals, a)
				}
			}
			s.PBAnimals = animals

			s.Selected = true

			shelter = *s
		}
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

	return Shelter{
		shelter,
		false,
		[]pb.Animal{},
		update,
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
