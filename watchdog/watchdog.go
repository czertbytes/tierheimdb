package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
	validShelters = []string{
		"arche-noah",
		"berlin",
		"dellbrueck",
		"dresden",
		"frankfurtmain",
		"franziskus-hamburg",
		"heppenheim",
		"muenchen",
		"samtpfoten-neukoelln",
	}

	tdbRoot = ""
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

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

	log.Printf("Starting sync for '%s'\n", args[1])

	if err := pb.RedisInit(); err != nil {
		log.Fatalf("Redis init failed! Error: %s", err)
	}

	shelterId, err := parseArgs(args)
	if err != nil {
		log.Fatalf("Parsing arguments failed! Error: %s", err)
	}

	animals, err := pb.RunShelterSync(shelterId)
	if err != nil {
		log.Fatalf("Crawling %s failed! Error: %s", shelterId, err)
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
