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

type ParsePaginationTest struct {
	File                  string
	Start, PerPage, Total int
}

func TestParsePagination(t *testing.T) {
	p := NewParser()
	tests := []ParsePaginationTest{
		ParsePaginationTest{
			"test_data/list.html",
			0,
			10,
			10,
		},
		ParsePaginationTest{
			"test_data/list2.html",
			0,
			10,
			143,
		},
	}

	for _, test := range tests {
		var file *os.File
		var err error

		if file, err = os.Open(test.File); err != nil {
			t.Errorf("Can't open input file!")
			return
		}

		start, perPage, total, err := p.ParsePagination(file)
		file.Close()
		if err != nil {
			t.Errorf("Error in parsing file! Error: %s", err)
			return
		}

		if start != test.Start || perPage != test.PerPage || total != test.Total {
			t.Errorf("Parsing pagination failed! Got: %d %d %d, Exp: %d %d %d", start, perPage, total, test.Start, test.PerPage, test.Total)
			return
		}
	}
}

func TestParseList(t *testing.T) {
	p := NewParser()
	tests := []ParseListTest{
		{
			"test_data/list.html",
			[]pb.Animal{
				pb.Animal{URL: "http://www.tierschutz-berlin.de/tierheim/tiervermittlung/katzen-sorgenkinder/ks-datenblatt.html?tx_realty_pi1%5BshowUid%5D=9692&cHash=b85409472cdd95571c9a4be3b7d2701e"},
				pb.Animal{URL: "http://www.tierschutz-berlin.de/tierheim/tiervermittlung/katzen-sorgenkinder/ks-datenblatt.html?tx_realty_pi1%5BshowUid%5D=6999&cHash=ae465bfec93cc52eb0534cd735bf539b"},
				pb.Animal{URL: "http://www.tierschutz-berlin.de/tierheim/tiervermittlung/katzen-sorgenkinder/ks-datenblatt.html?tx_realty_pi1%5BshowUid%5D=10158&cHash=55ff193807a9dc869140086ef6e3f090"},
				pb.Animal{URL: "http://www.tierschutz-berlin.de/tierheim/tiervermittlung/katzen-sorgenkinder/ks-datenblatt.html?tx_realty_pi1%5BshowUid%5D=10677&cHash=055ca0208faf8f36f2083c34632bf809"},
				pb.Animal{URL: "http://www.tierschutz-berlin.de/tierheim/tiervermittlung/katzen-sorgenkinder/ks-datenblatt.html?tx_realty_pi1%5BshowUid%5D=9880&cHash=46253aa126b04608f41560538feea440"},
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
			if a.URL != exp.URL {
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
			"test_data/detail.html",
			pb.Animal{
				Id:        "darryl",
				Name:      "Darryl",
				Sex:       "M",
				Breed:     "Cocker Spaniel",
				ShortDesc: "Darryl wurde am 22.04.09 geboren und ist etwa 46 cm groß.",
				LongDesc:  "Er ist ein eigensinniger Kerl, der eine Weile braucht um Vertrauen zu fassen. Wenn er sich bedrängt fühlt, hat er in der Vergangenheit leider gelernt, dass er sich mit den Zähnen verteidigen kann. Unsere Hundetrainerin arbeitet mit ihm und er macht sehr gute Fortschritte. Darryl liebt die Gesellschaft von Hündinnen und orientiert sich auch an diesen, so dass eine Vermittlung zu einer Hündin für ihn wünschenswert wäre. Ideal ist für ihn eine ruhige Umgebung und ein Garten bei geduldigen und hundeerfahrenen Menschen ohne Kinder. Vor der Vermittlung sollten Sie mehrere Besuche bei uns einplanen, damit er Vertrauen fassen kann. Unsere Hundetrainerin steht Ihnen dabei gerne zur Seite. Wenn Sie diesem anspruchsvollen Hund eine Chance geben möchten, dann melden Sie sich bitte bei den Tierpflegern unter 030 / 76 888 204.",
			},
		},
		{
			"test_data/detail_cat.html",
			pb.Animal{
				Id:        "pipusch",
				Name:      "Pipusch",
				Sex:       "M",
				Breed:     "Europäisch Kurzhaar",
				ShortDesc: "Pipusch musste leider ins Tierheim Berlin, weil sein Besitzer in ein Seniorenheim kam und den hübschen Kater nicht mitnehmen konnte. Der einstige Einzelkater ist hier angekommen noch etwas erschrocken von der neuen Situation.",
				LongDesc:  "Pipusch musste leider ins Tierheim Berlin, weil sein Besitzer in ein Seniorenheim kam und den hübschen Kater nicht mitnehmen konnte. Der einstige Einzelkater ist hier angekommen noch etwas erschrocken von der neuen Situation. Der Kater ist sehr gutmütig und möchte wieder in einen ruhigen Haushalt. Der Weg in den Garten sollte Pipusch offen stehen. Besuchen Sie den Kater Pipusch im Samtpfotenhaus 1. Bei Interesse können Sie auch anrufen und sich weitere Informationen über den schwarzweißen Kater einholen: 030 / 76 888 121.",
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

		exp := test.Animal
		if a.Id != exp.Id || a.Name != exp.Name || a.ShortDesc != exp.ShortDesc || a.LongDesc != exp.LongDesc {
			t.Errorf("Parsing animal failed! Got: %s Exp: %s", a, exp)
			return
		}
	}
}

/*
func TestParseDetailExtra(t *testing.T) {
	p := NewParser()
	tests := []ParseDetailTest{
		{
			"test_data/detail_extra.html",
			pb.Animal{
				Images: []pb.Image{
					pb.Image{URL: "http://www.tierschutz-berlin.de/typo3temp/pics/9fdfb7bbdb.jpg"},
					pb.Image{URL: "http://www.tierschutz-berlin.de/typo3temp/pics/162a910ffe.jpg"},
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

		a, err := p.ParseDetailExtra(file)
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
*/
