package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
	validShelters = []string{
		"samtpfoten-neukoelln",
		"franziskus-hamburg",
		"tierheim-dresden",
		"tierheim-berlin",
		"tierheim-muenchen",
		"tierheim-dellbrueck",
	}
)

var (
	tdbRoot = ""
)

func init() {
	tdbRoot = os.Getenv("GOPATH")
	if len(tdbRoot) == 0 {
		log.Fatalf("Environment variable GOPATH not set!")
	}
}

//  run me as
//  ./watchdog <shelterId>
//
//  example
//  ./watchdog samtpfoten-neukoelln
func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatalf("Invalid arguments!")
	}

	if err := pb.RedisInit(); err != nil {
		log.Fatalf("Redis init failed! Error: %s", err)
	}

	catnipName, err := parseArgs(args)
	if err != nil {
		log.Fatalf("Parsing arguments failed! Error: %s", err)
	}

	animals, err := fetchAnimals(catnipName)
	if err != nil {
		log.Fatalf("Crawling %s failed! Error: %s", catnipName, err)
	}

	if err := persist(catnipName, animals); err != nil {
		log.Fatalf("Persisting failed! Error: %s", err)
	}

	if err := backup(catnipName, animals); err != nil {
		log.Fatalf("Backing up failed! Error: %s", err)
	}
}

func parseArgs(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Invalid argument count!")
	}

	catnip := args[1]
	for _, s := range validShelters {
		if s == catnip {
			return catnip, nil
		}
	}

	return "", fmt.Errorf("Catnip not found!")
}

func fetchAnimals(catnipName string) ([]*pb.Animal, error) {
	return runCatnip(
		fmt.Sprintf("%s/bin/%s", tdbRoot, strings.Replace(catnipName, "-", "", -1)),
		fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/catnip/sources/%s.json", tdbRoot, catnipName),
	)
}

func persist(shelterId string, animals []*pb.Animal) error {
	var catsCounter int
	var dogsCounter int

	u := pb.NewUpdate(shelterId)
	for _, a := range animals {
		if len(a.Type) > 0 {
			switch a.Type {
			case "cat":
				catsCounter = catsCounter + 1
			case "dog":
				dogsCounter = dogsCounter + 1
			}
		}

		a.UpdateId = u.Id
		a.ShelterId = shelterId
	}

	_, err := pb.PutAnimals(animals)
	if err != nil {
		return err
	}

	u.Cats = catsCounter
	u.Dogs = dogsCounter

	return pb.PutUpdate(u)
}

func backup(catnipName string, animals []*pb.Animal) error {
	b, err := json.Marshal(animals)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fmt.Sprintf("%s/backup/%s.json", tdbRoot, catnipName), b, 0644)
}

func runCatnip(catnipPath, sourcesPath string) ([]*pb.Animal, error) {
	cmd := exec.Command(catnipPath, sourcesPath)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Creating stdoutPipe failed! Error: %s", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Running command failed! Error: %s", err)
	}

	var animals []*pb.Animal
	if err := json.NewDecoder(stdout).Decode(&animals); err != nil {
		return nil, fmt.Errorf("Decoding animal response failed! Error: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Waiting for animal response failed! Error: %s", err)
	}

	return animals, nil
}
