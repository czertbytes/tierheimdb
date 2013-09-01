package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"code.google.com/p/go.net/html"

	cp "github.com/czertbytes/tierheimdb/catnip"
	pb "github.com/czertbytes/tierheimdb/piggybank"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParsePagination(r io.Reader) (int, int, int, error) {
	panic("Not supported!")
}

func (p *Parser) ParseList(content io.Reader) ([]*pb.Animal, error) {
	doc, err := html.Parse(content)
	if err != nil {
		return nil, err
	}

	//  parse animals from list page
	animals := []*pb.Animal{}
	for _, animalNode := range p.listAnimalNodes(doc) {
		link := p.parseListAnimalLink(animalNode)
		if len(link) > 0 {
			animals = append(animals, &pb.Animal{
				URL: link,
			})
		}
	}

	return animals, nil
}

func (p *Parser) ParseDetail(content io.Reader) (*pb.Animal, error) {
	s := ""
	if b, err := ioutil.ReadAll(content); err == nil {
		s = string(b)
	}

	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return nil, err
	}

	//  parse animals from list page
	animalNodes := p.detailNodes(doc)
	if len(animalNodes) != 1 {
		return nil, errors.New("Animal detail not found!")
	}

	animalNode := animalNodes[0]
	name := p.parseName(animalNode)
	return &pb.Animal{
		Id:       cp.NormalizeName(name),
		Name:     name,
		Sex:      p.parseSex(animalNode),
		Breed:    p.parseBreed(animalNode),
		LongDesc: p.parseLongDesc(animalNode),
		Images:   p.parseImages(s),
	}, nil
}

func (p *Parser) ParseDetailExtra(r io.Reader) (*pb.Animal, error) {
	panic("Not supported!")
}

func (p *Parser) detailNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, p.selector([]string{
		"div.details",
	}))
}

func (p *Parser) parseName(node *html.Node) string {
	nameNodes := p.nameNodes(node)
	if len(nameNodes) != 1 {
		return ""
	}

	name := nameNodes[0].FirstChild.Data
	name = strings.Split(name, " ")[1]
	name = strings.ToLower(name)
	name = strings.ToUpper(name[0:1]) + name[1:]

	return name
}

func (p *Parser) nameNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.box1b",
		"ul",
		"li.name",
		"h1",
	})
}

func (p *Parser) parseSex(node *html.Node) string {
	sexNodes := p.sexNodes(node)
	if len(sexNodes) != 1 {
		return ""
	}

	sexNode := sexNodes[0]

	sex := ""
	if len(sexNode.FirstChild.Data) > 12 {
		sex = sexNode.FirstChild.Data[12:]
	}

	return sex
}

func (p *Parser) sexNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.box1b",
		"ul",
		"li.geschlecht",
	})
}

func (p *Parser) parseBreed(node *html.Node) string {
	breedNodes := p.breedNodes(node)
	if len(breedNodes) != 1 {
		return ""
	}

	breedNode := breedNodes[0]

	breed := ""
	if len(breedNode.FirstChild.Data) > 5 {
		breed = breedNode.FirstChild.Data[5:]
	}

	return breed
}

func (p *Parser) breedNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.box1b",
		"ul",
		"li.art",
	})
}

func (p *Parser) parseLongDesc(node *html.Node) string {
	longDescNodes := p.longDescNodes(node)
	if len(longDescNodes) != 1 {
		return ""
	}

	longDesc := longDescNodes[0].FirstChild.Data
	longDesc = strings.Trim(longDesc, " \n")

	return longDesc
}

func (p *Parser) longDescNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.bemerkung",
	})
}

func (p *Parser) listAnimalNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.liste",
		"div.galeriezelle",
	}))
}

func (p *Parser) parseListAnimalLink(node *html.Node) string {
	for _, linkNode := range p.detailLinkNodes(node) {
		return cp.NodeAttribute(linkNode, "href")
	}

	return ""
}

func (p *Parser) detailLinkNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.box1b",
		"ul",
		"li.link",
		"a",
	})
}

func (p *Parser) selector(s []string) []string {
	baseSelector := []string{
		"html",
		"body",
		"div#rahmen",
		"div#spalte2",
		"div",
		"div",
		"div.tierdb",
		"div.modultierdb",
	}

	return append(baseSelector, s...)
}
func (p *Parser) parseImages(s string) []pb.Image {
	images := []pb.Image{}

	for _, img := range regexp.MustCompile(`/thm/daten/tierdb/bilder/([^ ;]+);`).FindAllStringSubmatch(s, -1) {
		images = append(images, pb.Image{
			URL: fmt.Sprintf("http://tierschutzverein-muenchen.de/thm/daten/tierdb/bilder/%s", img[1]),
		})
	}

	return images
}
