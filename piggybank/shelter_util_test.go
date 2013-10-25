package piggybank

import (
	"testing"
)

func TestIntAnimalTypes(t *testing.T) {
	inputs := [][]string{
		{},
		{"cat"},
		{"dog"},
		{"cat", "dog"},
		{"cat", "dog", "XXX"},
		{"XXX"},
	}
	results := []int{0, 1, 2, 3, 3, 0}

	for i := 0; i < len(inputs); i++ {
		res := intAnimalTypes(inputs[i])
		if res != results[i] {
			t.Errorf("intAnimalTypes failed! Exp: '%d' Got: '%d'\n", res, results[i])
			return
		}
	}
}

func TestAnimalTypes(t *testing.T) {
	results := [][]string{
		{},
		{"cat"},
		{"dog"},
		{"cat", "dog"},
		{},
	}

	for i := 0; i < 4; i++ {
		res := animalTypes(i)
		if len(res) != len(results[i]) {
			t.Errorf("animalTypes failed! Exp: '%v' Got: '%v'\n", res, results[i])
			return
		}

		for j, r := range res {
			if r != results[i][j] {
				t.Errorf("animalTypes failed! Exp: '%v' Got: '%v'\n", r, results[i][j])
				return
			}
		}
	}
}

func TestParseLatLon(t *testing.T) {
	inputs := []string{
		"",
		"abcdef",
		"12.22;54.12",
		"12.22;abcde",
		"abcde;fghij",
	}
	results := [][]float64{
		{0, 0},
		{0, 0},
		{12.22, 54.12},
		{0, 0},
		{0, 0},
	}

	for i := 0; i < len(inputs); i++ {
		lat, lon, _ := parseLatLon(inputs[i])
		if lat != results[i][0] || lon != results[i][1] {
			t.Errorf("parseLatLon failed! Exp: '%f;%f' Got: '%v'\n", lat, lon, results[i])
			return
		}
	}
}
