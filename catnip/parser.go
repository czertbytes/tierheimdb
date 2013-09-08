package catnip

import (
	"io"
	"strings"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

type Parser interface {
	ParsePagination(r io.Reader) (int, int, int, error)
	ParseList(r io.Reader) ([]*pb.Animal, error)
	ParseDetail(r io.Reader) (*pb.Animal, error)
}

func NormalizeId(name string) string {
	return strings.ToLower(name)
}

func NormalizeBreed(breed string) string {
	return breed
}

func NormalizeSex(sex string) string {
	parsedSex := []string{}
	for _, token := range strings.Split(sex, " ,/") {
		if token == "männlich" {
			parsedSex = append(parsedSex, "M")
		}
		if token == "weiblich" {
			parsedSex = append(parsedSex, "F")
		}
	}

	return strings.Join(parsedSex, "/")
}
