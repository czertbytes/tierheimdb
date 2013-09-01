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
			8,
			17,
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
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_hamburg_standard_10001.php?f_mandant=bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5&f_bereich=Vermittlungstiere+Katzen&f_seite_max_ds=8&f_aktuelle_seite=1&f_aktueller_ds=53663&f_aktueller_ds_select=0&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_hamburg_standard_10001.php?f_mandant=bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5&f_bereich=Vermittlungstiere+Katzen&f_seite_max_ds=8&f_aktuelle_seite=1&f_aktueller_ds=53634&f_aktueller_ds_select=1&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_hamburg_standard_10001.php?f_mandant=bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5&f_bereich=Vermittlungstiere+Katzen&f_seite_max_ds=8&f_aktuelle_seite=1&f_aktueller_ds=53633&f_aktueller_ds_select=2&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_hamburg_standard_10001.php?f_mandant=bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5&f_bereich=Vermittlungstiere+Katzen&f_seite_max_ds=8&f_aktuelle_seite=1&f_aktueller_ds=53632&f_aktueller_ds_select=3&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_hamburg_standard_10001.php?f_mandant=bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5&f_bereich=Vermittlungstiere+Katzen&f_seite_max_ds=8&f_aktuelle_seite=1&f_aktueller_ds=53628&f_aktueller_ds_select=4&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_hamburg_standard_10001.php?f_mandant=bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5&f_bereich=Vermittlungstiere+Katzen&f_seite_max_ds=8&f_aktuelle_seite=1&f_aktueller_ds=53625&f_aktueller_ds_select=5&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_hamburg_standard_10001.php?f_mandant=bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5&f_bereich=Vermittlungstiere+Katzen&f_seite_max_ds=8&f_aktuelle_seite=1&f_aktueller_ds=52915&f_aktueller_ds_select=6&f_e_suche=&f_funktion=Detailansicht"},
				pb.Animal{URL: "http://presenter.comedius.de/design/bmt_hamburg_standard_10001.php?f_mandant=bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5&f_bereich=Vermittlungstiere+Katzen&f_seite_max_ds=8&f_aktuelle_seite=1&f_aktueller_ds=52914&f_aktueller_ds_select=7&f_e_suche=&f_funktion=Detailansicht"},
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
				Id:       "blacky und stromerle",
				Name:     "Blacky und Stromerle",
				LongDesc: "sind Fremden gegenüber misstrauisch, sie warten auf Menschen, die mit Geduld und Einfühlungsvermögen ihre Herzen erobern. Keine Familienkatzen",
				Images: []pb.Image{
					pb.Image{URL: "http://presenter.comedius.de/pic/bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5/file_1372691020_672101c82a2001d24f97"},
					pb.Image{URL: "http://presenter.comedius.de/pic/bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5/file_1372691035_1c12c0fd12b42d8a40fd"},
					pb.Image{URL: "http://presenter.comedius.de/pic/bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5/file_1372691046_aa7b7d0be7eba182a30a"},
					pb.Image{URL: "http://presenter.comedius.de/pic/bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5/file_1372691072_fe162c9d52f82292c046"},
					pb.Image{URL: "http://presenter.comedius.de/pic/bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5/file_1372691089_bb6a743abe7258a6c8b6"},
					pb.Image{URL: "http://presenter.comedius.de/pic/bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5/file_1372691111_90bf47a256e9df30be1a"},
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
		if a.Id != exp.Id || a.Name != exp.Name || a.LongDesc != exp.LongDesc {
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
