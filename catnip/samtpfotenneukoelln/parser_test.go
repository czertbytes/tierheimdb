package main

import (
	"os"
	"testing"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
	p = NewParser()
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
	tests := []ParseListTest{
		{
			"test_data/list.html",
			[]pb.Animal{
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/stupsi-brauni-beide-2011"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/isis-2007-und-diabolini-2009"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/peter-purzel-beide-2010"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/luzifer-und-balthazar"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/lilly"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/lilly-clemens-2004"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/butterfly-10-september-2007"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/jolie-2009"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/kalle-kalinka-kiana-kira"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/paul"},
				pb.Animal{URL: "http://www.samtpfoten-neukoelln.com/portfolio/barchen-1-august-2000-fritzchen-30-september-2001"},
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
	tests := []ParseDetailTest{
		{
			"test_data/detail.html",
			pb.Animal{
				Id:        "stupsi & brauni",
				Name:      "Stupsi & Brauni",
				ShortDesc: "Stupsi und Brauni sind ein Brüderpärchen",
				LongDesc:  "Stupsi und Brauni ergänzen sich hervorragend. Mit seinem rot-weißen Fell ist Brauni nicht nur eine echte Schönheit, sondern auch charakterlich ein Traum-Kater. Anfänglich etwas schüchtern, entwickelt er sich zum sanften Dauerschmuser ohne aufdringlich zu sein. Ein Platz am Fenster oder auf dem Balkon, ein weiches Körbchen, regelmäßige Streichel- und Spieleinheiten machen aus unserem Brauni einen sehr glücklichen Kater. Unser Stupsi ist ein richtiger Schatz, freundlich, neugierig verspielt und immer zu kleinen Späßen bereit. Er liebt hohe Kratzbäume und einen Ausguck von oben. Das Allerbeste ist aber, unseren Bürotisch zu annektieren. Nichts und niemand kommt dann zum Arbeiten, für zukünftige Halter eine prima Ausrede für deren Arbeitgeber. Und für wirklich entspannte Stunden. Stupsi und Brauni brauchen beide noch Paten!",
				Images: []pb.Image{
					pb.Image{URL: "http://www.samtpfoten-neukoelln.com/wp-content/uploads/2013/08/Brauni-5.jpg"},
					pb.Image{URL: "http://www.samtpfoten-neukoelln.com/wp-content/uploads/2013/08/Stupsi-3.jpg"},
					pb.Image{URL: "http://www.samtpfoten-neukoelln.com/wp-content/uploads/2013/08/Stupsi-2.jpg"},
					pb.Image{URL: "http://www.samtpfoten-neukoelln.com/wp-content/uploads/2013/08/Brauni-4.jpg"},
					pb.Image{URL: "http://www.samtpfoten-neukoelln.com/wp-content/uploads/2013/08/Brauni-3.jpg"},
					pb.Image{URL: "http://www.samtpfoten-neukoelln.com/wp-content/uploads/2013/08/Brauni-2.jpg"},
					pb.Image{URL: "http://www.samtpfoten-neukoelln.com/wp-content/uploads/2013/08/Stupsi.jpg"},
				},
			},
		},
		{
			"test_data/detail2.html",
			pb.Animal{
				Id:        "butterfly",
				Name:      "Butterfly",
				ShortDesc: "Butterfly sucht einen Menschen mit großem Herz zum Teilen",
				LongDesc:  "Butterfly ist eine neugierige und anhängliche Maus mit wunderschönem seidigen Fell. Wie mit Sternenstaub gepudert glänzt ihr Fell, sitzt sie in der Sonne. Im Moment teilt sie sich einen Raum mit Paul und kommt gut mit ihm zurecht. Andere Katzen braucht sie aber nicht unbedingt, dafür umso mehr Zuneigung und Beschäftigung. Tagsüber hält sich die Kleine stets in der Nähe ihrer Bezugspersonen auf, beobachtet sie und steht sofort bereit, um Streicheleinheiten oder ein “Leckerli” aufzuschnappen. Sie könnte ja etwas verpassen… Sie ist dabei aber nicht aufdringlich, sondern ruhig und angenehm. Sie ist aus etwas falsch verstandener Tierliebe, dank zuviel Leckerlies und falschen Futter übergewichtig geworden. Ansonsten ist sie kerngesund. Ihre Ernährung haben wir umgestellt, sie nimmt jetzt langsam und in vernünftigem Maße ab. Das sollte in ihrem neuen Zuhause fortgesetzt werden, sie wir dadurch an Agilität und Schönheit noch gewinnen. Butterfly sucht noch einen Paten.",
				Images: []pb.Image{
					pb.Image{URL: "http://www.samtpfoten-neukoelln.com/wp-content/uploads/2013/06/Butterfly-e1372308322902.jpg"},
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
