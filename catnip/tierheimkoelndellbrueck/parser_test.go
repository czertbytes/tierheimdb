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
			10,
			11,
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
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=55273&f_aktueller_ds_select=0&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54883&f_aktueller_ds_select=1&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54284&f_aktueller_ds_select=2&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54062&f_aktueller_ds_select=3&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=53824&f_aktueller_ds_select=4&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=53822&f_aktueller_ds_select=5&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=53702&f_aktueller_ds_select=6&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=51802&f_aktueller_ds_select=7&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=45582&f_aktueller_ds_select=8&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_koeln_standard_10001.php?f_mandant=bmt_koeln_d620d9faeeb43f717c893b5c818f1287&f_bereich=Vermittlung+kleine+Hunde+&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=34410&f_aktueller_ds_select=9&f_e_suche=&f_funktion=Detailansicht"},
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
				Id:       "nicki",
				Name:     "Nicki",
				Sex:      "M",
				Breed:    "Yorkshire Terrier",
				LongDesc: "Der kleine Nicki ist ein eher unsicherer Hund der bei uns eine Zeit brauchte, um Vertauen zu fassen. Er ist hundeverträglich und möchte auch wie einer behandelt werden, auch wenn er körperlich ein ganz kleiner ist. Nicki ist gerne aktiv unterwegs und genießt ausgiebige Spaziergänge, aber auch die Nähe zu seinen Menschen denen er vertraut.",
				Images: []pb.Image{
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/bmt_koeln_d620d9faeeb43f717c893b5c818f1287/file_1381489886_6fc2d6d0326b4701284c"},
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/bmt_koeln_d620d9faeeb43f717c893b5c818f1287/file_1381489898_46ed0ade15120f281de0"},
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
