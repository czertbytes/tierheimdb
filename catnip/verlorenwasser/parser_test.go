package main

import (
	"os"
	"testing"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

type ParseListTest struct {
	File    string
	Animals []pb.Animal
}

type ParseDetailTest struct {
	File   string
	Animal pb.Animal
}

func TestParseList(t *testing.T) {
	p := NewParser()
	tests := []ParseListTest{
		{
			"test_data/list-cats.html",
			[]pb.Animal{
				pb.Animal{
					Id:       "iben",
					Name:     "Iben",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2025-iben.html",
				},
				pb.Animal{
					Id:       "ingfried",
					Name:     "Ingfried",
					LongDesc: "Ingfried verzieht sich in die äußerste Ecke, am liebsten unter den Ofen. Er faucht nicht, man kann ihn hochheben - aber er mag es nicht so. Bei ihm sind wir sicher, dass es wird und er dankbar für ein eigenes Plätzchen wäre, auf dem man ihm die Chance auf ein liebevolles Leben gibt. Ingfried wurde 2008 geboren.",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/1227-ingfried.html",
				},
				pb.Animal{
					Id:       "panther",
					Name:     "Panther",
					LongDesc: "Panther wurde 2005 geboren. Dieser wunderschöne Kater ist sehr gelehrig und unternehmungslustig. Er wäre mit einer reinen Wohnungshaltung nicht zufrieden, sondern möchte auch draußen alles erforschen. Wie man sieht mag er Höhlen, Unterschlüpfe und Verstecke.",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/88-panther.html",
				},
			},
		},
	}

	for _, test := range tests {
		var file *os.File
		var err error

		if file, err = os.Open(test.File); err != nil {
			t.Errorf("Can't open input file!")
			return
		}

		animals, err := p.ParseList(file)
		file.Close()
		if err != nil {
			t.Errorf("Error in parsing file! Error: %s", err)
			return
		}

		if len(animals) != len(test.Animals) {
			t.Errorf("Animals size does not match! Got: %d Exp: %d", len(animals), len(test.Animals))
			return
		}

		for i, a := range animals {
			exp := test.Animals[i]
			if a.Id != exp.Id || a.Name != exp.Name || a.LongDesc != exp.LongDesc || a.URL != exp.URL {
				t.Errorf("Parsing DetailSources failed! Got: %s Exp: %s", a, exp)
				return
			}
		}
	}
}

func TestParseDetail(t *testing.T) {
	p := NewParser()
	tests := []ParseDetailTest{
		{
			"test_data/detail-cat.html",
			pb.Animal{
				Images: []pb.Image{
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/katzen/350_iben/iben_01.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/katzen/350_iben/iben_02.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/katzen/350_iben/iben_03.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/katzen/350_iben/iben_04.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/katzen/350_iben/iben_05.jpg"},
				},
			},
		},
		{
			"test_data/detail-dog.html",
			pb.Animal{
				Images: []pb.Image{
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/hunde/1599_bernice/bernice_01.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/hunde/1599_bernice/bernice_02.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/hunde/1599_bernice/bernice_03.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/hunde/1599_bernice/bernice_04.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/hunde/1599_bernice/bernice_05.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/hunde/1599_bernice/bernice_06.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/hunde/1599_bernice/bernice_07.jpg"},
					pb.Image{URL: "http://www.tierheim-verlorenwasser.de/images/stories/tiere/hunde/1599_bernice/bernice_08.jpg"},
				},
			},
		},
	}

	for _, test := range tests {
		var file *os.File
		var err error

		if file, err = os.Open(test.File); err != nil {
			t.Errorf("Can't open input file!")
			return
		}

		a, err := p.ParseDetail(file)
		file.Close()
		if err != nil {
			t.Errorf("Error in parsing file! Error: %s", err)
			return
		}

		for i, img := range a.Images {
			exp := test.Animal.Images[i]
			if img.URL != exp.URL {
				t.Errorf("Pargins Animal images failed! Got: %s Exp: %s", img, exp)
				return
			}
		}
	}
}
