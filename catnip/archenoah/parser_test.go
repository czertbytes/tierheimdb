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
			"test_data/list-cats.html",
			1,
			6,
			20,
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
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Kleine+Hunde&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=54694&f_aktueller_ds_select=0&f_e_suche=&f_funktion=Detailansicht",
					Id:    "struppi",
					Name:  "Struppi",
					Sex:   "M",
					Breed: "Terrier Mix",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Kleine+Hunde&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=54687&f_aktueller_ds_select=1&f_e_suche=&f_funktion=Detailansicht",
					Id:    "emma",
					Name:  "Emma",
					Sex:   "F",
					Breed: "Terrier-Mix",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Kleine+Hunde&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=53272&f_aktueller_ds_select=2&f_e_suche=&f_funktion=Detailansicht",
					Id:    "simba",
					Name:  "Simba",
					Sex:   "F",
					Breed: "Mix",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Kleine+Hunde&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=53267&f_aktueller_ds_select=3&f_e_suche=&f_funktion=Detailansicht",
					Id:    "gria",
					Name:  "Gria",
					Sex:   "F",
					Breed: "Mix",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Kleine+Hunde&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=52702&f_aktueller_ds_select=4&f_e_suche=&f_funktion=Detailansicht",
					Id:    "rambo",
					Name:  "Rambo",
					Sex:   "M",
					Breed: "Terrier Mix",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Kleine+Hunde&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=50931&f_aktueller_ds_select=5&f_e_suche=&f_funktion=Detailansicht",
					Id:    "tara",
					Name:  "Tara",
					Sex:   "F",
					Breed: "Mix",
				},
			},
		},
		{
			"test_data/list-cats.html",
			[]pb.Animal{
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Tiervermittlung+Katzen&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=53738&f_aktueller_ds_select=0&f_e_suche=&f_funktion=Detailansicht",
					Id:    "linus",
					Name:  "Linus",
					Sex:   "M",
					Breed: "Europäisch Kurzhaar",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Tiervermittlung+Katzen&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=53736&f_aktueller_ds_select=1&f_e_suche=&f_funktion=Detailansicht",
					Id:    "lorry",
					Name:  "Lorry",
					Sex:   "F",
					Breed: "Europäisch Kurzhaar",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Tiervermittlung+Katzen&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=53735&f_aktueller_ds_select=2&f_e_suche=&f_funktion=Detailansicht",
					Id:    "lotta",
					Name:  "Lotta",
					Sex:   "F",
					Breed: "Europäisch Kurzhaar",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Tiervermittlung+Katzen&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=53732&f_aktueller_ds_select=3&f_e_suche=&f_funktion=Detailansicht",
					Id:    "mora",
					Name:  "Mora",
					Sex:   "F",
					Breed: "Europäisch Kurzhaar",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Tiervermittlung+Katzen&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=53730&f_aktueller_ds_select=4&f_e_suche=&f_funktion=Detailansicht",
					Id:    "maili",
					Name:  "Maili",
					Sex:   "F",
					Breed: "Europäisch Kurzhaar",
				},
				pb.Animal{
					URL:   "http://presenter.comedius.de/design/bmt_brinkum2_standard_10001.php?f_mandant=bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce&f_bereich=Tiervermittlung+Katzen&f_seite_max_ds=6&f_aktuelle_seite=1&f_aktueller_ds=52903&f_aktueller_ds_select=5&f_e_suche=&f_funktion=Detailansicht",
					Id:    "olli",
					Name:  "Olli",
					Sex:   "M",
					Breed: "Europäisch Kurzhaar",
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
			"test_data/detail-dog1.html",
			pb.Animal{
				LongDesc: "Struppi ist ein netter, mit Artgenossen meist verträglicher Rüde. Er ist Fremden gegenüber erst unsicher, baut aber recht schnell Vertrauen auf. Bislang ist er noch recht Handscheu und schreckhaft. Wenn es zu stressig wird schnappt er auch mal, wenn man aber ruhig mit ihm umgeht lässt er alles zu.",
				Images: []pb.Image{
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce/file_1379079642_03968b8711f30d6c6bb2"},
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce/file_1379079679_c9c826eb828528291819"},
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce/file_1379079688_e41fe3e2449fec48cca9"},
				},
			},
		},
		{
			"test_data/detail-dog2.html",
			pb.Animal{
				LongDesc: "Pauls Lebensgeschichte ist leider nicht lückenlos nachzuvollziehen. Sicher ist, dass er durch sehr viele Hände gereicht wurde, so dass er bislang nicht die Chance hatte, ein geregeltes Leben kennen zu lernen, und somit natürlich viele Unsicherheiten in sich trägt. Bei uns im Tierheimalltag ist Paul komplett unauffällig. Er begegnet Menschen und auch fast allen Hunden freundlich. Durch 2 Vorbesitzer, die sich angefunden haben, konnten wir aber erfahren, daß er sich in Wohnbereichen sehr schnell einlebt und die Probleme erst dann ans Tageslicht kommen. Paul neigt dazu, weibliche Bezugspersonen zu bevorzugen und auch zu vereinnahmen. Diese möchte er beschützen und wird dann anwesenden Männern gegenüber sehr kritisch. Es soll auch schon Beißvorfälle gegeben haben. Wir suchen für Paul hundeerfahrene Menschen die bereit sind die Arbeit mit ihm aufzunehmen und die Weichen nochmal ganz neu zu stellen, um ihm endlich ein geregeltes, ausgeglichenes Leben zu ermöglichen.",
				Images: []pb.Image{
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce/file_1375909181_6ee5c1485a4cccff76b2"},
					pb.Image{Width: 300, URL: "http://presenter.comedius.de/pic/bmt_brinkum_131a3dc5d9a36a62409a455f08fa9dce/file_1375909196_8fbf0563eb2675cbafe5"},
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
