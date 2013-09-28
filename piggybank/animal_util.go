package piggybank

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

func uniqueAnimals(animals []*Animal) []*Animal {
	unique := []*Animal{}

	hashes := make(map[string]bool)
	names := make(map[string]bool)
	for _, a := range animals {
		ah := hashAnimal(a)
		if _, found := hashes[ah]; found == true {
			continue
		}

		if _, found := names[a.Id]; found == true {
			a.Id = fmt.Sprintf("%s-%d", a.Id, time.Now().UnixNano())
		}

		unique = append(unique, a)

		hashes[ah] = true
		names[a.Id] = true
	}

	return unique
}

func hashAnimal(a *Animal) string {
	h := md5.New()
	io.WriteString(h, a.Name)
	io.WriteString(h, a.Sex)
	io.WriteString(h, a.Breed)
	io.WriteString(h, a.ShortDesc)
	io.WriteString(h, a.LongDesc)
	for _, i := range a.Images {
		io.WriteString(h, i.URL)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
