package catnip

import (
	"unicode/utf8"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func MergeAnimals(a, o *pb.Animal) {
	if len(o.Id) > 0 {
		a.Id = o.Id
	}

	if len(o.URL) > 0 {
		a.URL = o.URL
	}

	if len(o.Name) > 0 {
		a.Name = o.Name
	}

	if len(o.Type) > 0 {
		a.Type = o.Type
	}

	if len(o.Breed) > 0 {
		a.Breed = o.Breed
	}

	if len(o.Sex) > 0 {
		a.Sex = o.Sex
	}

	if len(o.ShortDesc) > 0 {
		a.ShortDesc = o.ShortDesc
	}

	if len(o.LongDesc) > 0 {
		a.LongDesc = o.LongDesc
	}

	if o.Images != nil && len(o.Images) > 0 {
		a.Images = o.Images
	}
}

func ToUTF8(s string) string {
	if !utf8.ValidString(s) {
		bytes := []byte(s)
		buf := make([]rune, len(bytes))
		for i, b := range bytes {
			buf[i] = rune(b)
		}
		return string(buf)
	}

	return s
}
