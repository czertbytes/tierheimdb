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
			1,
			5,
			33,
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
				pb.Animal{URL: "http://presenter.comedius.de/design/tsv_frankfurt_standard_10001.php?f_mandant=tsv_ffm_f8bc3055161419854ea9ddb99936b98e&f_bereich=Katzen&f_seite_max_ds=5&f_aktuelle_seite=1&f_aktueller_ds=54049&f_aktueller_ds_select=0&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/tsv_frankfurt_standard_10001.php?f_mandant=tsv_ffm_f8bc3055161419854ea9ddb99936b98e&f_bereich=Katzen&f_seite_max_ds=5&f_aktuelle_seite=1&f_aktueller_ds=54023&f_aktueller_ds_select=1&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/tsv_frankfurt_standard_10001.php?f_mandant=tsv_ffm_f8bc3055161419854ea9ddb99936b98e&f_bereich=Katzen&f_seite_max_ds=5&f_aktuelle_seite=1&f_aktueller_ds=54021&f_aktueller_ds_select=2&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/tsv_frankfurt_standard_10001.php?f_mandant=tsv_ffm_f8bc3055161419854ea9ddb99936b98e&f_bereich=Katzen&f_seite_max_ds=5&f_aktuelle_seite=1&f_aktueller_ds=54020&f_aktueller_ds_select=3&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/tsv_frankfurt_standard_10001.php?f_mandant=tsv_ffm_f8bc3055161419854ea9ddb99936b98e&f_bereich=Katzen&f_seite_max_ds=5&f_aktuelle_seite=1&f_aktueller_ds=54014&f_aktueller_ds_select=4&f_e_suche=&f_funktion=Detailansicht"},
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
			t.Errorf("pb.Animals size does not match! Got: %d Exp: %d", len(animals), len(test.Animals))
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
				Id:       "lisa - freigängerin",
				Name:     "Lisa - Freigängerin",
				Sex:      "F",
				Breed:    "Europäisch Kurzhaar",
				LongDesc: "Wirbelwind sucht schönes Zuhause! Hallo ich bin die süsse und freche Lisa! Ich möchte gerne als Zweitkatze vermittelt werden. Ich glaube sonst wirds mir einfach zu langweilig. Ich bin sehr aktiv und immer auf Achse. Am liebsten wäre mir daher eine grooooooooooße Wohnung und Freigang. Schön wäre es auch, wenn ich viel Beschäftigung (Spass & Spiel) bekommen könnte. Meldet euch, wenn ihr Lust und Zeit für einen kleinen Wirbelwind, wie ich es bin, habt. Eure Lisa",
				Images: []pb.Image{
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/tsv_ffm_f8bc3055161419854ea9ddb99936b98e/file_1376642063_8d50581ea44dab5b6fa9"},
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/tsv_ffm_f8bc3055161419854ea9ddb99936b98e/file_1376642076_400dd3d6d670c98497d3"},
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

		exp := test.Animal
		if a.Id != exp.Id ||
			a.Name != exp.Name ||
			a.Sex != exp.Sex ||
			a.Breed != exp.Breed ||
			a.LongDesc != exp.LongDesc {
			t.Errorf("Parsing animal failed! Got: %s Exp: %s", a, exp)
			return
		}

		for i, img := range a.Images {
			exp := test.Animal.Images[i]
			if img.URL != exp.URL {
				t.Errorf("Pargins pb.Animal images failed! Got: %s Exp: %s", img, exp)
				return
			}
		}
	}
}
