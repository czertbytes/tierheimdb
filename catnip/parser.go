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

func NormalizeSex(s string) string {
	if len(s) == 0 {
		return s
	}

	t := PrepareStringChunk(s)
	t = strings.ToLower(t)
	t = strings.Replace(t, "/", " ", -1)
	t = strings.Replace(t, ",", " ", -1)

	parsedSex := []string{}
	for _, token := range strings.Split(t, " ") {
		if token == "männlich" || token == "rüde" {
			parsedSex = append(parsedSex, "M")
		}
		if token == "weiblich" || token == "hündin" || token == "weibl." {
			parsedSex = append(parsedSex, "F")
		}
	}

	return strings.Join(parsedSex, "/")
}

func PrepareStringChunk(s string) string {
	if len(s) == 0 {
		return s
	}

	t := strings.Trim(ToUTF8(s), " ")
	t = strings.Replace(t, "\u0009", "", -1)
	t = strings.Replace(t, "\u000A", "", -1)
	t = strings.Replace(t, "\u00A0", "", -1)
	t = strings.Trim(t, " ")

	return t
}
