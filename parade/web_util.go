package main

import (
	"fmt"
	"log"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func makeShelters() (Shelters, error) {
	pbShelters, err := pb.GetEnabledShelters()
	if err != nil {
		log.Println("Getting enabled shelters failed!")
		return nil, err
	}

	shelters := Shelters{}
	for _, pbShelter := range pbShelters {
		shelter, err := makeShelter(pbShelter)
		if err != nil {
			log.Printf("Make shelter '%s' failed!", pbShelter.Id)
			return nil, err
		}

		shelters = append(shelters, &shelter)
	}

	return shelters, nil
}

func makeSheltersShelter(shelterId, animalType string) (Shelters, Shelter, error) {
	shelters, err := makeShelters()
	if err != nil {
		return nil, Shelter{}, err
	}

	shelter, err := selectedShelter(shelters, shelterId, animalType)
	if err != nil {
		return nil, Shelter{}, err
	}

	return shelters, shelter, nil
}

func makeSheltersShelterAnimal(shelterId, updateId, animalId string) (Shelters, Shelter, pb.Animal, error) {
	shelters, err := makeShelters()
	if err != nil {
		return nil, Shelter{}, pb.Animal{}, err
	}

	shelter, err := selectedShelter(shelters, shelterId, "")
	if err != nil {
		return nil, Shelter{}, pb.Animal{}, err
	}

	animal, err := pb.GetAnimal(shelterId, updateId, animalId)
	if err != nil {
		return nil, Shelter{}, pb.Animal{}, err
	}

	return shelters, shelter, animal, nil
}

func makeShelter(shelter pb.Shelter) (Shelter, error) {
	update, err := pb.GetLastUpdate(shelter.Id)
	if err != nil {
		log.Println("Getting last update for '%s' failed!", shelter.Id)
		return Shelter{}, err
	}

	return Shelter{
		shelter,
		false,
		[]pb.Animal{},
		update,
	}, nil
}

func selectedShelter(shelters Shelters, shelterId, animalType string) (Shelter, error) {
	var shelter Shelter
	for _, s := range shelters {
		if s.PBShelter.Id == shelterId {

			animals := []pb.Animal{}
			animals, err := pb.GetAnimals(s.PBShelter.Id, s.PBUpdate.Id, animalType)
			if err != nil {
				return Shelter{}, err
			}

			s.PBAnimals = addImagePlaceholders(animals)
			s.Selected = true

			shelter = *s
		}
	}

	return shelter, nil
}

func addImagePlaceholders(animals []pb.Animal) []pb.Animal {
	animalsIP := []pb.Animal{}
	for _, a := range animals {
		if len(a.Images) == 0 {
			a.Images = []pb.Image{
				pb.Image{
					URL: fmt.Sprintf("http://placehold.it/200x200&text=%s", a.Name),
				},
			}
		}
		animalsIP = append(animalsIP, a)
	}

	return animalsIP
}
