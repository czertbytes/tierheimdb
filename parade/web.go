package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "text/template"

    pb "github.com/czertbytes/tierheimdb/piggybank"
)

var tmpl *template.Template

type Sitemap struct {
    Animals []pb.Animal
}

func init() {
    log.SetFlags(log.LstdFlags | log.Lshortfile)

    tdbRoot := os.Getenv("GOPATH")
    if len(tdbRoot) == 0 {
        log.Fatalf("Environment variable GOPATH not set!")
    }

    tmpl = template.Must(
        template.New("sitemap").
            ParseGlob(fmt.Sprintf("%s/src/github.com/czertbytes/tierheimdb/parade/templates/*.tmpl.xml", tdbRoot)))
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
    tmpl.ExecuteTemplate(w, "sitemap", &Sitemap{animals})
}
