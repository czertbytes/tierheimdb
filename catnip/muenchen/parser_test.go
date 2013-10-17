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
			"test_data/list.html",
			[]pb.Animal{
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1211160445332797746:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1307300912214594207:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1304130445026303051:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1303070447477558117:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1306150445049542382:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1304250445029688036:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1308160446114733374:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1307260445084344016:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1307260445086308939:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1307260445084462513:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1308070445382144233:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1307270445124506247:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1308030445308926053:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1308200446175273305:und::0:15::1::0:::::"},
				pb.Animal{URL: "http://tierschutzverein-muenchen.de/de/hundeklein.htm?q=0&acltierdb=1000:2000:1308200446175934494:und::0:15::1::0:::::"},
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
				Id:       "gina",
				Name:     "Gina",
				Sex:      "Weiblich",
				LongDesc: "Dackelmischling Gina ist etwa im Jahr 2009 geboren, kastriert und kam zu uns, weil ihr Besitzer mit ihr überfordert war. Da Gina einen starken Willen hat und zur Leinenaggression gegen Artgenossen neigt, sollten ihren neuen Besitzer Hundeerfahrung haben und sie sehr konsequent führen. Läuft sie frei, ist sie mit Rüden verträglich und zum Spielen aufgelegt. Die sportliche Hündin könnte gut an eine Familie vermittelt werden, da sie Kinder aller Altersklassen kennt. Die kleine Hündin hat leider auch einen ausgeprägten Jagdtrieb, Kleintiere und Katzen sollten deshalb nicht in ihrem zu Hause leben. Wenn sie entsprechend beschäftigt wurde, kann sie aber auch mal ganz gut ein paar Stunden alleine bleiben. Der Besuch einer Hundeschule wäre von Vorteil. Wer sich für Gina interessiert, meldet sich in den Hundetrakten zu unseren Besuchszeiten. Telefonische Auskünfte über sie bekommt man unter der Rufnummer 089-921 000-26.",
				Images: []pb.Image{
					pb.Image{URL: "http://tierschutzverein-muenchen.de/thm/daten/tierdb/bilder/130815(1).JPG"},
					pb.Image{URL: "http://tierschutzverein-muenchen.de/thm/daten/tierdb/bilder/130815-2.jpg"},
					pb.Image{URL: "http://tierschutzverein-muenchen.de/thm/daten/tierdb/bilder/130815-3.JPG"},
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
		if a.Id != exp.Id || a.Name != exp.Name || a.ShortDesc != exp.ShortDesc || a.LongDesc != exp.LongDesc {
			t.Errorf("Parsing animal failed! Got: %s Exp: %s", a, exp)
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
