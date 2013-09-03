package catnip

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

type Catnip struct {
	p Parser
}

func NewCatnip(p Parser) *Catnip {
	return &Catnip{p}
}

func (c *Catnip) Run(args []string) {
	if len(args) != 2 {
		log.Fatalf("Invalid arguments!")
	}

	sources, err := c.loadSources(args[1])
	if err != nil {
		log.Fatalf("Loading sources from %s failed! Error: %s", args[1], err)
	}

	animals, err := c.run(sources)
	if err != nil {
		log.Fatalf("Fetching animals failed! Error: %s", err)
	}

	json.NewEncoder(os.Stdout).Encode(animals)
}

func (c *Catnip) loadSources(path string) ([]*Source, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sources []*Source
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		return nil, err
	}

	return sources, nil
}

func (c *Catnip) run(sources []*Source) ([]*pb.Animal, error) {
	animalsChan := make(chan []*pb.Animal)
	errChan := make(chan error)

	for _, source := range sources {
		go func(s *Source) {
			fetched, err := c.Fetch(s)
			if err != nil {
				errChan <- err
			}
			animalsChan <- fetched
		}(source)
	}

	counter := len(sources)
	animals := []*pb.Animal{}
	for {
		select {
		case fetched := <-animalsChan:
			animals = append(animals, fetched...)

			counter = counter - 1
			if counter <= 0 {
				return animals, nil
			}
		case err := <-errChan:
			return nil, err
		}
	}
}

func (c *Catnip) Fetch(s *Source) ([]*pb.Animal, error) {
	sources := []*Source{s}
	if len(s.Pagination) > 0 {
		var err error
		sources, err = PaginatedSources(c.p, s)
		if err != nil {
			return nil, err
		}
	}

	animalChan := make(chan *pb.Animal)
	errorChan := make(chan error)
	counter := 0
	animals := []*pb.Animal{}
	for _, ps := range sources {
		la, err := ParseList(c.p, ps.URL)
		if err != nil {
			return nil, err
		}

		counter = counter + len(la)
		for _, animal := range la {
			go func(a *pb.Animal) {
				a.Type = s.Animal
				a.Priority = s.Priority

				switch s.Type {
				case "none":
					a.URL = s.URL
				case "detail":
					da, err := ParseDetail(c.p, a.URL)
					if err != nil {
						errorChan <- err
					}

					MergeAnimals(a, da)
				default:
					errorChan <- errors.New(fmt.Sprintf("Unknown Detail type: %s!", s.Type))
					return
				}

				animalChan <- a
			}(animal)
		}
	}

	for {
		select {
		case animal := <-animalChan:
			animals = append(animals, animal)

			counter = counter - 1
			if counter <= 0 {
				return animals, nil
			}
		case err := <-errorChan:
			return nil, err
		}
	}
}
