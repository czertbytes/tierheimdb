package main

import (
	"os"
	"testing"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

type ParseDetailListTest struct {
	File    string
	Animals []pb.Animal
}

func TestParseDetailList(t *testing.T) {
	p := NewParser()
	tests := []ParseDetailListTest{
		{
			"test_data/detail_list.html",
			[]pb.Animal{
				pb.Animal{
					Id:       "timo",
					Name:     "Timo",
					Sex:      "männlich",
					Breed:    "europäisch Kurzhaar",
					LongDesc: "Unser Timo ist ein stolzer Kater mit weißem Fell und schwarzen Flecken. Timo ist Menschen gegenüber sehr liebevoll und zugänglich. Er holt sich täglich seine Streicheleinheiten ab. Der Kater akzeptiert keine anderen Vierbeiner in seiner Nähe, denn er möchte gern den Platzhirsch spielen und ganz wichtig: draußen herumstromern dürfen, wie und wann es ihm gefällt.",
					Images: []pb.Image{
						pb.Image{URL: "http://www.dresden.de/media/bilder/tiere/tierheim/336_K_139_13_Timo_CIMG1619.jpg"},
					},
				},
				pb.Animal{
					Id:       "emeli",
					Name:     "Emeli",
					Sex:      "weiblich",
					Breed:    "europäisch Kurzhaar",
					LongDesc: "Emeli wurde aus Zeitgründen in unserem Tierheim abgegeben. Sie ist eine liebevolle Katzendame, sehr verschmust und anhänglich. Da unsere gemütliche Emeli ein total unkompliziertes Wesen hat und sich sehr gut mit Artgenossen verträgt, eignet sie sich als Anfängertier oder auch als Zweittier, mit täglichem Freigang.",
					Images: []pb.Image{
						pb.Image{URL: "http://www.dresden.de/media/bilder/tiere/tierheim/336_K_41_13_Emeli_CIMG1610.jpg"},
					},
				},
			},
		},
		{
			"test_data/detail_list2.html",
			[]pb.Animal{
				pb.Animal{
					Id:       "opa",
					Name:     "Opa",
					Sex:      "männlich",
					Breed:    "Spitz - Mix",
					LongDesc: "Opa wurde von der Feuerwehr auf der Gertud-Caspari-Straße aufgefunden und ins Tierheim gebracht. Der Rüde war in einem sehr schlechten Gesundheitszustand, hatte kaum noch Fell, sehr ängstlich und schüchtern. Man muss sehr behutsam mit ihm umgehen, da der Rüde sehr schreckhaft ist und sich sofort einigelt. Die kahlen Stellen brauchen noch Zeit um nachzuwachsen. Opa braucht jetzt liebevolle Menschen, welche mit Geduld und viel Liebe versuchen aus ihm wieder einen fröhlichen Hund zu machen, auch er hat eine zweite Chance verdient.",
					Images: []pb.Image{
						pb.Image{URL: "http://www.dresden.de/media/bilder/tiere/tierheim/336_H_165_13_Opa_CIMG1886.jpg"},
					},
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
			if a.Id != exp.Id || a.Name != exp.Name || a.LongDesc != exp.LongDesc || a.Sex != exp.Sex || a.Breed != exp.Breed {
				t.Errorf("Parsing animal failed!\nGot: %s\nExp: %s", a, exp)
				return
			}

			for j, img := range a.Images {
				expImg := exp.Images[j]
				if img.URL != expImg.URL {
					t.Errorf("Pargins Animal images failed!\nGot: %s\nExp: %s", img, expImg)
					return
				}
			}
		}
	}
}
