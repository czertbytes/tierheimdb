package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func fetchAnimals(shelterId string) ([]*pb.Animal, error) {
	tdbRoot := os.Getenv("TIERHEIMDB_ROOT")
	if len(tdbRoot) == 0 {
		return nil, fmt.Errorf("Environment variable TIERHEIMDB_ROOT not set!")
	}

	return runCatnip(
		fmt.Sprintf("%s/bin/%s", tdbRoot, strings.Replace(shelterId, "-", "", -1)), // TODO: hmmm
		fmt.Sprintf("%s/sources/%s.json", tdbRoot, shelterId),
	)
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
