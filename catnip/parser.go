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

var (
	MaleSexKeywords = []string{
		"männlich",
		"männl",
		"rüde",
	}

	FemaleSexKeywords = []string{
		"weiblich",
		"weibl",
		"hündin",
	}

	SexKeywords = []string{}
)

func init() {
	SexKeywords = append(SexKeywords, MaleSexKeywords...)
	SexKeywords = append(SexKeywords, FemaleSexKeywords...)
}

func NormalizeId(s string) string {
	if len(s) == 0 {
		return s
	}

	t := PrepareStringChunk(s)
	t = strings.ToLower(t)
	t = strings.Replace(t, "/ reserviert", "", -1)
	t = strings.Trim(t, " ")

	return t
}

func NormalizeBreed(s string) string {
	if len(s) == 0 {
		return s
	}

	t := PrepareStringChunk(s)
	for _, s := range []string{"Katze", "Hund"} {
		t = strings.Replace(t, s, "", -1)
	}
	t = strings.Replace(t, " -Mix", "-Mix", -1)
	t = strings.Replace(t, "EKH", "Europäisch Kurzhaar", -1)
	t = strings.Trim(t, " /")

	if len(t) > 1 {
		t = strings.ToUpper(t[0:1]) + t[1:]
	}

	return t
}

func NormalizeSex(s string) string {
	if len(s) == 0 {
		return s
	}

	t := PrepareStringChunk(s)
	t = strings.ToLower(t)
	for _, s := range []string{"/", ",", ":", "."} {
		t = strings.Replace(t, s, " ", -1)
	}

	parsedSex := []string{}
	for _, token := range strings.Split(t, " ") {
		for _, s := range MaleSexKeywords {
			if token == s {
				parsedSex = append(parsedSex, "M")
			}
		}

		for _, s := range FemaleSexKeywords {
			if token == s {
				parsedSex = append(parsedSex, "F")
			}
		}
	}

	return strings.Join(parsedSex, "/")
}

func PrepareStringChunk(s string) string {
	if len(s) == 0 {
		return s
	}

	t := strings.Trim(ToUTF8(s), " ")
	for _, s := range []string{"\u0009", "\u000A", "\u00A0"} {
		t = strings.Replace(t, s, "", -1)
	}
	t = strings.Trim(t, " ")

	return t
}
