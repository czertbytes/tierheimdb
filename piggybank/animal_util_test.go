package piggybank

import (
	"testing"
)

type UniqueAnimalsTest struct {
	Animals []*Animal
	Unique  []*Animal
}

func TestUniqueAnimals(t *testing.T) {
	tests := []UniqueAnimalsTest{
		{
			[]*Animal{
				&Animal{
					Id:        "charly",
					Name:      "Charly",
					Sex:       "M",
					Breed:     "breed1",
					ShortDesc: "short1",
					LongDesc:  "long1",
					Images: []Image{
						Image{
							URL: "url1",
						},
					},
				},
				&Animal{
					Id:        "charly",
					Name:      "Charly",
					Sex:       "M",
					Breed:     "breed1",
					ShortDesc: "short1",
					LongDesc:  "long1",
					Images: []Image{
						Image{
							URL: "url1",
						},
					},
				},
				&Animal{
					Id:        "charly",
					Name:      "Charly",
					Sex:       "W",
					Breed:     "breed2",
					ShortDesc: "short2",
					LongDesc:  "long2",
					Images: []Image{
						Image{
							URL: "url2",
						},
					},
				},
			},
			[]*Animal{
				&Animal{
					Id:        "charly",
					Name:      "Charly",
					Sex:       "M",
					Breed:     "breed1",
					ShortDesc: "short1",
					LongDesc:  "long1",
					Images: []Image{
						Image{
							URL: "url1",
						},
					},
				},
				&Animal{
					Id:        "charly-1234",
					Name:      "Charly",
					Sex:       "W",
					Breed:     "breed2",
					ShortDesc: "short2",
					LongDesc:  "long2",
					Images: []Image{
						Image{
							URL: "url2",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		uAnimals := uniqueAnimals(test.Animals)
		if len(test.Unique) != len(uAnimals) {
			t.Errorf("Unique Animals does not match! Got: %d Exp: %d", len(test.Unique), len(uAnimals))
			return
		}

		for i, a := range uAnimals {
			exp := test.Unique[i]
			if a.Name != exp.Name {
				t.Errorf("Unique Name is not as expected! Got: %s Exp: %s", a.Name, exp.Name)
				return
			}
			if a.Sex != exp.Sex {
				t.Errorf("Unique Sex is not as expected! Got: %s Exp: %s", a.Sex, exp.Sex)
				return
			}
			if a.Breed != exp.Breed {
				t.Errorf("Unique Breed is not as expected! Got: %s Exp: %s", a.Breed, exp.Breed)
				return
			}
			if a.ShortDesc != exp.ShortDesc {
				t.Errorf("Unique ShortDesc is not as expected! Got: %s Exp: %s", a.ShortDesc, exp.ShortDesc)
				return
			}
			if a.LongDesc != exp.LongDesc {
				t.Errorf("Unique LongDesc is not as expected! Got: %s Exp: %s", a.LongDesc, exp.LongDesc)
				return
			}

			for j, img := range a.Images {
				expImg := exp.Images[j].URL
				if img.URL != expImg {
					t.Errorf("Unique Image is not as expected! Got: %s Exp: %s", img.URL, expImg)
					return
				}
			}
		}
	}
}
