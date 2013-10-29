package piggybank

import (
	"fmt"
	"testing"
)

func TestShelterPagination(t *testing.T) {
	shelters := Shelters{}
	for i := 1; i < 6; i++ {
		shelters = append(shelters, Shelter{Id: fmt.Sprintf("Shelter %d", i)})
	}

	inputs := []Pagination{
		Pagination{
			0,
			3,
		},
		Pagination{
			3,
			1,
		},
		Pagination{
			3,
			0,
		},
		Pagination{
			10,
			3,
		},
		Pagination{
			4,
			10,
		},
	}

	results := []Shelters{
		Shelters{
			Shelter{Id: "Shelter 1"},
			Shelter{Id: "Shelter 2"},
			Shelter{Id: "Shelter 3"},
		},
		Shelters{
			Shelter{Id: "Shelter 4"},
		},
		Shelters{},
		Shelters{},
		Shelters{
			Shelter{Id: "Shelter 5"},
		},
	}

	for i := 0; i < len(inputs); i++ {
		res := shelters.Paginate(inputs[i])
		shelterRes := results[i]

		for i, res := range res {
			exp := shelterRes[i]
			if exp.Id != res.Id {
				t.Errorf("shelter pagination failed! Exp: '%v' Got '%v'\n", exp, res)
				return
			}
		}
	}
}
