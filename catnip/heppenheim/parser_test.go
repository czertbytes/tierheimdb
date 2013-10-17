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
			"test_data/list-dogs.html",
			1,
			10,
			19,
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
			"test_data/list-dogs.html",
			[]pb.Animal{
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54828&f_aktueller_ds_select=0&f_e_suche=&f_funktion=Detailansicht",
					Id:    "ruscha",
					Name:  "Ruscha",
					Sex:   "F",
					Breed: "Labrador-Mischling",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54827&f_aktueller_ds_select=1&f_e_suche=&f_funktion=Detailansicht",
					Id:    "donny und kolga",
					Name:  "Donny und Kolga",
					Sex:   "M/F",
					Breed: "Weimaraner",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54810&f_aktueller_ds_select=2&f_e_suche=&f_funktion=Detailansicht",
					Id:    "ates",
					Name:  "Ates",
					Sex:   "M",
					Breed: "",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54809&f_aktueller_ds_select=3&f_e_suche=&f_funktion=Detailansicht",
					Id:    "spike",
					Name:  "Spike",
					Sex:   "M",
					Breed: "Husky-Schäferhund-Mischling",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54807&f_aktueller_ds_select=4&f_e_suche=&f_funktion=Detailansicht",
					Id:    "gina",
					Name:  "Gina",
					Sex:   "F",
					Breed: "Schäferhund",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=53941&f_aktueller_ds_select=5&f_e_suche=&f_funktion=Detailansicht",
					Id:    "shirley",
					Name:  "Shirley",
					Sex:   "F",
					Breed: "Dackel",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=52274&f_aktueller_ds_select=6&f_e_suche=&f_funktion=Detailansicht",
					Id:    "jackson",
					Name:  "Jackson",
					Sex:   "M",
					Breed: "Bordercolliemischling",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=51134&f_aktueller_ds_select=7&f_e_suche=&f_funktion=Detailansicht",
					Id:    "brownie",
					Name:  "Brownie",
					Sex:   "M",
					Breed: "Beagle",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=49652&f_aktueller_ds_select=8&f_e_suche=&f_funktion=Detailansicht",
					Id:    "jayjo",
					Name:  "Jayjo",
					Sex:   "M",
					Breed: "Staffordshire Terrier Mischling",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Hunde&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=47608&f_aktueller_ds_select=9&f_e_suche=&f_funktion=Detailansicht",
					Id:    "rocky",
					Name:  "Rocky",
					Sex:   "M",
					Breed: "Rottweiler",
				},
			},
		},
		{
			"test_data/list-cats.html",
			[]pb.Animal{
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=55164&f_aktueller_ds_select=0&f_e_suche=&f_funktion=Detailansicht",
					Id:    "dörte",
					Name:  "Dörte",
					Sex:   "",
					Breed: "",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=55026&f_aktueller_ds_select=1&f_e_suche=&f_funktion=Detailansicht",
					Id:    "fee",
					Name:  "Fee",
					Sex:   "",
					Breed: "Persermix",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54842&f_aktueller_ds_select=2&f_e_suche=&f_funktion=Detailansicht",
					Id:    "nico",
					Name:  "Nico",
					Sex:   "M",
					Breed: "",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54841&f_aktueller_ds_select=3&f_e_suche=&f_funktion=Detailansicht",
					Id:    "benny",
					Name:  "Benny",
					Sex:   "",
					Breed: "",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54840&f_aktueller_ds_select=4&f_e_suche=&f_funktion=Detailansicht",
					Id:    "wuschel",
					Name:  "Wuschel",
					Sex:   "",
					Breed: "",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54835&f_aktueller_ds_select=5&f_e_suche=&f_funktion=Detailansicht",
					Id:    "basti",
					Name:  "Basti",
					Sex:   "",
					Breed: "",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54833&f_aktueller_ds_select=6&f_e_suche=&f_funktion=Detailansicht",
					Id:    "sepp",
					Name:  "Sepp",
					Sex:   "",
					Breed: "",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54832&f_aktueller_ds_select=7&f_e_suche=&f_funktion=Detailansicht",
					Id:    "blue",
					Name:  "Blue",
					Sex:   "",
					Breed: "",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54831&f_aktueller_ds_select=8&f_e_suche=&f_funktion=Detailansicht",
					Id:    "mufasa",
					Name:  "Mufasa",
					Sex:   "M",
					Breed: "",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php?f_mandant=tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e&f_bereich=Katzen&f_seite_max_ds=10&f_aktuelle_seite=1&f_aktueller_ds=54830&f_aktueller_ds_select=9&f_e_suche=&f_funktion=Detailansicht",
					Id:    "romeo",
					Name:  "Romeo",
					Sex:   "F",
					Breed: "",
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
			t.Errorf("pb.Animals size does not match! Got: %d Exp: %d", len(animals), len(test.Animals))
			return
		}

		for i, a := range animals {
			exp := test.Animals[i]
			if a.Id != exp.Id || a.Name != exp.Name || a.LongDesc != exp.LongDesc || a.Sex != exp.Sex || a.Breed != exp.Breed || a.URL != exp.URL {
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
			"test_data/detail-dog.html",
			pb.Animal{
				LongDesc: "Rusha musste leider ins Tierheim, da ihr Frauchen verstorben ist. Bei uns zeigt sie sich als sehr gut verträgliche und nette Hündin. Allerdings kann sie auch gut ihre Ressourcen verteidigen und schon mal einen Besucher zwicken, der ihr nicht passt. Daher möchten wir Rusha als Einzeltier in hundeerfahrene Hände vermitteln. Rusha hat einen guten Grundgehorsam und mit ihr wurde auch fleißig auf dem Hundeplatz Bensheim geübt. Sie bleibt gut ein paar Stunden alleine. Bei den richtigen Menschen ist Rusha garantiert ein toller Hund.",
				Images: []pb.Image{
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e/file_1379606554_fcfaa75512e48c8c5f75"},
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e/file_1379606571_fb20e83a2f0af251f529"},
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e/file_1379606579_392a4d0a9490c6efb31a"},
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/tierheim_heppenheim_dba261758ae4da4d90d508d1cbd11d1e/file_1379606585_fd3f9d2bef0d1720d1bf"},
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
		if a.LongDesc != exp.LongDesc {
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
