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

func NormalizeName(name string) string {
	return strings.ToLower(name)
}

func NormalizeBreed(breed string) string {
	return breed
}

func NormalizeSex(sex string) string {
	return sex
}
