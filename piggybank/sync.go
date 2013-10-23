package piggybank

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	tdbRoot = ""
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	tdbRoot = os.Getenv("GOPATH")
	if len(tdbRoot) == 0 {
		log.Fatalf("Environment variable GOPATH not set!")
	}
}

func RunShelterSync(shelterId string) ([]*Animal, error) {
	animals, err := runCatnip(shelterId)
	if err != nil {
		return nil, fmt.Errorf("Crawling %s failed! Error: %s", shelterId, err)
	}

	if err := persist(shelterId, animals); err != nil {
		return nil, fmt.Errorf("Persisting failed! Error: %s", err)
	}

	if err := backup(shelterId, animals); err != nil {
		return nil, fmt.Errorf("Backing up failed! Error: %s", err)
	}

	return animals, nil
}

func runCatnip(shelterId string) ([]*Animal, error) {
	catnipPath := fmt.Sprintf("%s/bin/%s", tdbRoot, strings.Replace(shelterId, "-", "", -1))
	sourcesPath := fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/catnip/sources/%s.json", tdbRoot, shelterId)

	cmd := exec.Command(catnipPath, sourcesPath)
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Creating stdoutPipe failed! Error: %s", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Running command failed! Error: %s", err)
	}

	var animals []*Animal
	if err := json.NewDecoder(stdout).Decode(&animals); err != nil {
		return nil, fmt.Errorf("Decoding animal response failed! Error: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Waiting for animal response failed! Error: %s", err)
	}

	return animals, nil
}

func persist(shelterId string, animals []*Animal) error {
	var catsCounter int
	var dogsCounter int

	u := NewUpdate(shelterId)
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

	_, err := PutAnimals(animals)
	if err != nil {
		return err
	}

	u.Cats = catsCounter
	u.Dogs = dogsCounter

	return PutUpdate(u)
}

func backup(catnipName string, animals []*Animal) error {
	b, err := json.Marshal(animals)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/backup/%s.json", tdbRoot, catnipName), b, 0644); err != nil {
		return err
	}

	return commitChange(catnipName)
}

func commitChange(catnipName string) error {
	backupFile := fmt.Sprintf("%s.json", catnipName)
	backupDir := fmt.Sprintf("%s/backup", tdbRoot)
	backupMessage := fmt.Sprintf("\"Sync at %s\"", time.Now().Format("Jan 2, 2006 at 3:04pm"))

	cmd := exec.Command("git", "commit", "-i", backupFile, "-m", backupMessage)
	cmd.Dir = backupDir
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Running backup commit failed! Error: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("Waiting for backup commit failed! Error: %s", err)
	}

	return nil
}
