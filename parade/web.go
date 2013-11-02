package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	t_template "text/template"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
	tmpl   *template.Template
	t_tmpl *t_template.Template
)

type Shelter struct {
	PBShelter pb.Shelter
	Selected  bool
	PBAnimals []pb.Animal
	PBUpdate  pb.Update
}

type Shelters []*Shelter

type HomePage struct {
	Title    string
	Shelters Shelters
}

type ShelterPage struct {
	Title    string
	Shelters Shelters
	Shelter  Shelter
}

type AnimalPage struct {
	Title    string
	Shelters Shelters
	Shelter  Shelter
	Animal   pb.Animal
}

type AboutPage struct {
	Title    string
	Shelters Shelters
}

type Sitemap struct {
	Animals []pb.Animal
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	tdbRoot := os.Getenv("GOPATH")
	if len(tdbRoot) == 0 {
		log.Fatalf("Environment variable GOPATH not set!")
	}

	t_tmpl = t_template.Must(
		t_template.New("sitemap").
			ParseGlob(fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/tmpl/*.xml.tmpl", tdbRoot)))
}

func GetSitemapHandler(w http.ResponseWriter, r *http.Request) {
	shelters, err := pb.GetEnabledShelters()
	if err != nil {
		log.Println(err)
		return
	}

	animals := []pb.Animal{}
	for _, s := range shelters {
		update, err := pb.GetLastUpdate(s.Id)
		if err != nil {
			log.Println(err)
			return
		}

		as, err := pb.GetAnimals(s.Id, update.Id, "", pb.Pagination{0, 999})
		if err != nil {
			log.Println(err)
			return
		}

		animals = append(animals, as...)
	}

	w.Header().Add("Content-Type", "application/xml; charset=utf-8")
	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>"))
	t_tmpl.ExecuteTemplate(w, "sitemap", &Sitemap{
		animals,
	})
}
