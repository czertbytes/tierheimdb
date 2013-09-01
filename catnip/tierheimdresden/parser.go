package main

import (
	"bytes"
	"fmt"
	"io"
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
	var animals []*pb.Animal

	doc, err := html.Parse(content)
	if err != nil {
		return animals, err
	}

	//  parse complete animals
	for _, animalNode := range p.listAnimalNodes(doc) {
		name := p.parseName(animalNode)
		if len(name) > 0 {
			animals = append(animals, &pb.Animal{
				Id:       cp.NormalizeName(name),
				Name:     name,
				Breed:    p.parseBreed(animalNode),
				Sex:      p.parseSex(animalNode),
				LongDesc: p.parseDesc(animalNode),
				Images:   p.parseImages(animalNode),
			})
		}
	}

	return animals, nil
}

func (p *Parser) ParseDetail(content io.Reader) (*pb.Animal, error) {
	panic("Not supported!")
}

func (p *Parser) ParseDetailExtra(r io.Reader) (*pb.Animal, error) {
	panic("Not supported!")
}

func (p *Parser) listAnimalNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, []string{
		"html",
		"body#inhalt",
		"div.contentcontainerborder",
		"div.contentcontainerbg",
		"div#contentcontainer",
		"div#content",
		"div.contentelements",
		"div.block",
	})
}

func (p *Parser) parseName(doc *html.Node) string {
	var name string

	nameNodes := p.nameNodes(doc)
	if len(nameNodes) != 1 {
		return ""
	}

	name = nameNodes[0].LastChild.Data
	//      fix for http://www.dresden.de/de/02/070/111/01/hunde/05_listenhunde.php
	if name == "Hunderassen" {
		return ""
	}

	name = strings.Replace(name, "&quot;", "\"", -1)
	name = strings.Replace(name, "&nbsp;", "", -1)
	name = strings.Replace(name, "„", "\"", -1)
	name = strings.Replace(name, "»", "\"", -1)
	name = strings.Replace(name, "«", "\"", -1)
	nameStart := strings.Index(name, "\"")
	if nameStart > 0 {
		nameEnd := strings.LastIndex(name, "\"")
		name = name[nameStart:nameEnd]
	}
	name = strings.Replace(name, "\"", "", -1)
	name = strings.Trim(name, " ")

	return name
}

func (p *Parser) nameNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, []string{
		"div.rightcol",
		"h3",
	})
}

func (p *Parser) normalizeName(name string) string {
	return strings.ToLower(name)
}

func (p *Parser) parseBreed(doc *html.Node) string {
	var breed string

	detailNodes := p.detailNodes(doc)
	if len(detailNodes) > 0 {
		for c := detailNodes[0].FirstChild; c != nil; c = c.NextSibling {
			if c.DataAtom.String() == "strong" {
				//  iterate one again inside the strong element
				for c2 := c.FirstChild; c2 != nil; c2 = c2.NextSibling {
					if strings.Contains(c2.Data, "Rasse") {
						breed = c.NextSibling.Data
					}
				}
			}
		}
	}

	breed = strings.Replace(breed, "\u00A0", "", -1)
	breed = strings.Trim(breed, " :\n")

	return breed
}

func (p *Parser) detailNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, []string{
		"div.rightcol",
		"p",
	})
}

func (p *Parser) parseSex(doc *html.Node) string {
	var sex string

	detailNodes := p.detailNodes(doc)
	if len(detailNodes) > 0 {
		for c := detailNodes[0].FirstChild; c != nil; c = c.NextSibling {
			if c.DataAtom.String() == "strong" && strings.Contains(c.FirstChild.Data, "Geschlecht") {
				sex = c.NextSibling.Data
			}
		}
	}

	sex = strings.Replace(sex, "\u00A0", "", -1)
	sex = strings.Trim(sex, " /:\n")

	return sex
}

func (p *Parser) parseDesc(doc *html.Node) string {
	var longDescBuffer bytes.Buffer

	detailNodes := p.detailNodes(doc)
	if len(detailNodes) > 0 {
		descNode := detailNodes[len(detailNodes)-1]

		for c := descNode.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				longDescBuffer.WriteString(c.Data)
				longDescBuffer.WriteString(" ")
			}
		}
	}

	longDesc := strings.Trim(longDescBuffer.String(), " ")
	longDesc = strings.Replace(longDesc, "\u000A", "", -1)
	longDesc = strings.Replace(longDesc, "\u00A0", " ", -1)
	longDesc = strings.Replace(longDesc, "\u0009", "", -1)
	longDesc = strings.Replace(longDesc, "  ", " ", -1)
	longDesc = strings.Trim(longDesc, " \n")

	return longDesc
}

func (p *Parser) parseImages(doc *html.Node) []pb.Image {
	imageNodes := cp.NodeSelect(doc, []string{
		"div.leftcol",
		"div.pic",
		"img",
	})

	if len(imageNodes) == 0 {
		imageNodes = cp.NodeSelect(doc, []string{
			"div.leftcol",
			"div.pic",
			"div",
			"img",
		})
	}

	images := []pb.Image{}
	if len(imageNodes) > 0 {
		imageUrl := cp.NodeAttribute(imageNodes[0], "src")
		imageUrl = strings.Trim(imageUrl, " ")

		images = append(images, pb.Image{
			URL: fmt.Sprintf("http://www.dresden.de%s", imageUrl),
		})
	}

	return images
}
