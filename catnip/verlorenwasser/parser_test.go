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
			"test_data/list-cats.xml",
			[]pb.Animal{
				pb.Animal{
					Id:       "iben",
					Name:     "Iben",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2025-iben.html",
				},
				pb.Animal{
					Id:       "amela",
					Name:     "Amela",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2024-amela.html",
				},
				pb.Animal{
					Id:       "paddy",
					Name:     "Paddy",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2023-paddy.html",
				},
				pb.Animal{
					Id:       "abbie",
					Name:     "Abbie",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2022-abbie.html",
				},
				pb.Animal{
					Id:       "jarko",
					Name:     "Jarko",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2021-jarko.html",
				},
				pb.Animal{
					Id:       "stina",
					Name:     "Stina",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2020-stina.html",
				},
				pb.Animal{
					Id:       "annika",
					Name:     "Annika",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2019-annika.html",
				},
				pb.Animal{
					Id:       "ambra",
					Name:     "Ambra",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2018-ambra.html",
				},
				pb.Animal{
					Id:       "steen",
					Name:     "Steen",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2017-stenn.html",
				},
				pb.Animal{
					Id:       "kiri",
					Name:     "Kiri",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/katzen/2016-kiri.html",
				},
			},
		},
		{
			"test_data/list-dogs.xml",
			[]pb.Animal{
				pb.Animal{
					Id:       "lino",
					Name:     "Lino",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/2013-lino.html",
				},
				pb.Animal{
					Id:       "tamy",
					Name:     "Tamy",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/2012-tamy.html",
				},
				pb.Animal{
					Id:       "anja",
					Name:     "Anja",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/2011-anja.html",
				},
				pb.Animal{
					Id:       "maja",
					Name:     "Maja",
					LongDesc: "",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/2010-maja.html",
				},
				pb.Animal{
					Id:       "marlitt",
					Name:     "Marlitt",
					LongDesc: "Marlitt ist eine ganz liebe kleine Maus. Manchmal hat es den Anschein, als ob sie die Menschen richtig anlächelt. Anfangs zeigt sie sich mitunter zwar etwas schüchtern, aber wenn man sie näher kennenlernt, baut sie Vertrauen auf. Da sie sich in der Gesellschaft anderer Hunde sehr wohl fühlt, wäre sie als Zweithund gut geeignet. Dann könnte sie sich an einem anderen souveränen Hund gut orientieren. Mit ihren Kumpels spielt und tobt sie ausgiebig. Sie fühlt sich in der Gruppe gut aufgehoben und kommt mit allen gut aus. Auch bei Katzen und Kleintieren zeigt sie sich verträglich. Da sie als Welpe wahrscheinlich vieles nicht kennenglernt hat, fehlt ihr das nötige Selbstvertrauen. So sind ihr manche Geräusche nicht so angenehm. Sie reagiert nicht panisch, aber ein ruhiges Zuhause würde sie bevorzugen. Kleinere Kinder - vor allem hektische und schrille - sollten nicht vorhanden sein. Marlitt wurde 2011 geboren. Sie wiegt 12 kg.",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/1798-marlitt.html",
				},
				pb.Animal{
					Id:       "bernice",
					Name:     "Bernice",
					LongDesc: "Bernice ist eine ganz liebe und verschmuste Foxterrier-Mischlingshündin. Bei Menschen ist sie sehr anhänglich. Und auch mit Artgenossen versteht sie sich sehr gut. Sowohl der Umgang mit Rüden wie auch Hündinnen sind unproblematisch für sie. Bei Gassitouren ist sie gern vorneweg - und immer mit der Nase am Boden. Sie liebt Schnüffeltouren über alles. Die ausgesprochen freundliche Hündin ist kinderlieb und rundum verträglich. Bernice wurde 2006 geboren und wiegt nur 13,5 kg.",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/2002-bernice.html",
				},
				pb.Animal{
					Id:       "sonngard",
					Name:     "Sonngard",
					LongDesc: "Sonngard sieht auf dem Foto wie ein normaler Schäferhundmix aus - aber sie ist viel kleiner. Die Fotos täuschen: Sonngard wiegt nur 12 kg. Bitte vergleichen Sie die Fotos mit den Fotos anderer Hunde und deren Gewichtsangaben. Die junge Hündin ist sehr freundlich, aber noch eine Spur unsicher. Größere Kinder stellen wahrscheinlich kein Problem für sie dar - bei kleineren könnte es ihr eventuell zu hektisch werden. Sonngard versteht sich mit Rüden und Hündinnen gleichermaßen - könnte also auch gut als Zweithund gehalten werden. Dann gewinnt sie auch mehr Sicherheit und wird nicht so angespannt und nervös. Leinenführigkeit üben wir zur Zeit, denn sie kennt nichts. Sie ist da eher auf dem Stand eines Welpen, obwohl sie schon zwei Jahre alt ist. Aber da sie sehr aufgeweckt ist, wird sie alles schnell lernen.",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/2001-sonngard.html",
				},
				pb.Animal{
					Id:       "anieke",
					Name:     "Anieke",
					LongDesc: "Anieke ist ebenfalls eine kleine Hündin. Sie sieht viel größer auf den Fotos aus. Aber sie wiegt nur 10 kg und ist damit sogar zarter als der Dackelverschnitt Vin mit seinen 13,5 kg. Anieke versteht sich mit Hündinnen und Rüden. Sie lebt gern in der Gruppe unter ihren Artgenossen. Mit Menschen hat sie es allerdings nicht so sehr. Sie zeigt sich bei uns recht scheu und unsicher. Aber wenn die Richtigen kommen, wird sie Vertrauen fassen und sich über ein neues ruhiges Zuhause sehr freuen. Anieke wurde 2012 geboren.",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/2000-anieke.html",
				},
				pb.Animal{
					Id:       "vin",
					Name:     "Vin",
					LongDesc: "Vin ist ein Dackelmischling. Mit seinen kurzen Beinchen ist er ein lustiges kleines Kerlchen. Wir halten ihn für einen Hund, der auch für Anfänger in der Hundehaltung gut geeignet wäre, denn Vin ist lieb, nett, freundlich, ruhig, ausgeglichen ... Er macht keinem etwas streitig, kämpft nicht um Futter - ist einfach nur zurückhaltend und angenehm. Vin schmust für sein Leben gern und liebt es, gestreichelt zu werden. Er wurde 2005 geboren und wiegt 13,5 kg.",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/1999-vin.html",
				},
				pb.Animal{
					Id:       "romina",
					Name:     "Romina",
					LongDesc: "Romina ist ein Schnauzer-Pudel-Mischling. Die junge aufgeweckte Hündin ist freundlich und angenehm vom Wesen her. Sie versteht sich auch mit ihren Artgenossen. Romina lebte früher mit mehreren Hunden in einem Haushalt (einer davon ist Pauletta). Allerdings ging es ihr dort nicht so gut. Davon zeugt auch ihr Ohr, das seltsam absteht. Es handelt sich um ein unbehandeltes Blutohr (Bluterguss an der Ohrmuschel, der immer wieder einblutet). Wir haben es bereits operieren lassen und hoffen, es wird sich alles richten. Romina bleibt nicht gern allein. Sie muss sich erst einleben und anpassen. So mag sie auch noch keine Leine, denn richtige Gassitouren kennt sie gar nicht. Aber da Romina freundlich und quirlig-neugierig ist, wird sie das bisher Versäumte schnell nachholen. Wir halten sie für eine gute Familienhündin. Sie wäre auch als Zweithund sehr geeignet. Romina wurde 2011 geboren und wiegt 14 kg.",
					URL:      "http://www.tierheim-verlorenwasser.de/unsere-tiere/hunde/1998-romina.html",
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
