package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"code.google.com/p/go.net/html"

	cp "github.com/czertbytes/tierheimdb/catnip"
	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
	imageRE = regexp.MustCompile(`/thm/daten/tierdb/bilder/([^ ;]+);`)
)

type Parser struct {
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParsePagination(r io.Reader) (int, int, int, error) {
	panic("Not supported!")
}

func (p *Parser) ParseList(r io.Reader) ([]*pb.Animal, error) {
	var animals []*pb.Animal

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	for _, animalNode := range p.listAnimalNodes(doc) {
		name := cp.NormalizeName(p.parseName(animalNode))
		link := p.parseListAnimalLink(animalNode)

		if len(name) > 0 && len(link) > 0 {
			animals = append(animals, &pb.Animal{
				Id:    cp.NormalizeId(name),
				Name:  name,
				Sex:   cp.NormalizeSex(p.parseSex(animalNode)),
				Breed: cp.NormalizeBreed(p.parseBreed(animalNode)),
				URL:   link,
			})
		}
	}

	return animals, nil
}

func (p *Parser) ParseDetail(r io.Reader) (*pb.Animal, error) {
	s, doc, err := p.parseContent(r)
	if err != nil {
		return nil, err
	}

	animalNodes := p.detailNodes(doc)
	if len(animalNodes) != 1 {
		return nil, errors.New("Animal detail not found!")
	}

	animalNode := animalNodes[0]
	return &pb.Animal{
		LongDesc: p.parseLongDesc(animalNode),
		Images:   p.parseImages(s),
	}, nil
}

func (p *Parser) ParseDetailExtra(r io.Reader) (*pb.Animal, error) {
	panic("Not supported!")
}

func (p *Parser) parseContent(r io.Reader) (string, *html.Node, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil, err
	}

	s := string(b)
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return "", nil, err
	}

	return s, doc, nil
}

func (p *Parser) parseName(node *html.Node) string {
	nameNodes := p.nameNodes(node)
	if len(nameNodes) != 1 {
		log.Println("Parse error! NameNode not found!")
		return ""
	}

	if nameNodes[0].FirstChild != nil {
		name := nameNodes[0].FirstChild.Data
		if len(name) > 8 {
			splittedName := strings.Split(name, " ")
			if len(splittedName) == 2 && len(splittedName[1]) > 0 {
				name = strings.ToLower(splittedName[1])
				return strings.ToUpper(name[0:1]) + name[1:]
			}
		}
	}

	log.Println("Parse error! Name not found!")
	return ""
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

func (p *Parser) parseLongDesc(node *html.Node) string {
	longDescNodes := p.longDescNodes(node)
	if len(longDescNodes) != 1 {
		return ""
	}

	longDesc := longDescNodes[0].FirstChild.Data
	longDesc = strings.Trim(longDesc, " \n")

	return longDesc
}

func (p *Parser) parseListAnimalLink(node *html.Node) string {
	for _, linkNode := range p.detailLinkNodes(node) {
		return cp.NodeAttribute(linkNode, "href")
	}

	return ""
}

func (p *Parser) parseImages(s string) []pb.Image {
	images := []pb.Image{}

	for _, img := range imageRE.FindAllStringSubmatch(s, -1) {
		images = append(images, pb.Image{
			URL: fmt.Sprintf("http://tierschutzverein-muenchen.de/thm/daten/tierdb/bilder/%s", img[1]),
		})
	}

	return images
}

func (p *Parser) detailNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, p.selector([]string{
		"div.details",
	}))
}

func (p *Parser) nameNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.box1b",
		"ul",
		"li.name",
		"h1",
	})
}

func (p *Parser) sexNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.box1b",
		"ul",
		"li.geschlecht",
	})
}

func (p *Parser) breedNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.box1b",
		"ul",
		"li.art",
	})
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
