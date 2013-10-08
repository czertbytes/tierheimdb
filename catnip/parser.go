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
	name = strings.ToLower(name)
	name = strings.Replace(name, "/ reserviert", "", -1)
	name = strings.Trim(name, " ")

	return name
}

func NormalizeBreed(breed string) string {
	breed = strings.Replace(breed, "Katze", "", -1)
	breed = strings.Replace(breed, "Hund", "", -1)
	breed = strings.Replace(breed, " -Mix", "-Mix", -1)
	breed = strings.Replace(breed, "EKH", "Europäisch Kurzhaar", -1)
	breed = strings.Trim(breed, " /")

	if len(breed) > 1 {
		breed = strings.ToUpper(breed[0:1]) + breed[1:]
	}

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
