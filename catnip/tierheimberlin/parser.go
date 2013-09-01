package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	doc, err := html.Parse(r)
	if err != nil {
		return 0, 0, 0, err
	}

	paginationNodes := p.paginationNodes(doc)
	paginationNodesLen := len(paginationNodes)
	if paginationNodesLen > 0 {
		totalAnimals := paginationNodes[paginationNodesLen-1].FirstChild.Data
		totalAnimals = strings.Replace(totalAnimals, "\u00A0", "", -1)
		totalAnimals = strings.Replace(totalAnimals, "Tiere", "", -1)
		totalAnimals = strings.Trim(totalAnimals, " ()")

		total, err := strconv.Atoi(totalAnimals)
		if err != nil {
			return 0, 0, 0, err
		}

		return 0, 10, total, nil
	}

	return 0, 10, 10, nil
}

func (p *Parser) ParseList(r io.Reader) ([]*pb.Animal, error) {
	var animals []*pb.Animal

	doc, err := html.Parse(r)
	if err != nil {
		return animals, err
	}

	for _, animalNode := range p.listAnimalNodes(doc) {
		link := cp.NodeAttribute(animalNode, "href")
		animals = append(animals, &pb.Animal{
			URL: fmt.Sprintf("http://www.tierschutz-berlin.de/%s", link),
		})
	}

	return animals, nil
}

func (p *Parser) ParseDetail(r io.Reader) (*pb.Animal, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	name := p.parseName(doc)
	shortDesc, longDesc := p.parseDescriptions(doc)

	images, err := p.parseImages(doc)
	if err != nil {
		return nil, err
	}

	return &pb.Animal{
		Id:        cp.NormalizeName(name),
		Name:      name,
		Breed:     p.parseBreed(doc),
		Sex:       p.parseSex(doc),
		ShortDesc: shortDesc,
		LongDesc:  longDesc,
		Images:    images,
	}, nil
}

func (p *Parser) parseImages(doc *html.Node) ([]pb.Image, error) {
	galleryLinkNodes := p.galleryNodes(doc)
	if len(galleryLinkNodes) < 1 {
		return nil, errors.New("Gallery link not found!")
	}

	galleryLinkNode := galleryLinkNodes[0]
	galleryLink := cp.NodeAttribute(galleryLinkNode, "onclick")
	galleryLink = strings.Split(galleryLink, "'")[1]

	resp, err := http.Get(galleryLink)
	if err != nil {
		return nil, err
	}

	galleryDoc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	images := []pb.Image{}
	for _, imageNode := range p.imageNodes(galleryDoc) {
		url := cp.NodeAttribute(imageNode, "onclick")
		url = strings.Split(url, "'")[1]

		images = append(images, pb.Image{
			URL: fmt.Sprintf("http://www.tierschutz-berlin.de/%s", url),
		})
	}

	return images, nil
}

func (p *Parser) parseName(doc *html.Node) string {
	nameNodes := p.nameNodes(doc)
	if len(nameNodes) != 1 {
		return ""
	}

	name := nameNodes[0].FirstChild.Data

	name = strings.Replace(name, "\u000A", "", -1)
	name = strings.Replace(name, "\u0009", "", -1)
	name = strings.Trim(name, " ")

	return name
}

func (p *Parser) parseBreed(doc *html.Node) string {
	overviewRows := p.overviewRows(doc)
	if len(overviewRows) < 1 {
		return ""
	}

	breed := ""
	for i := 0; i < len(overviewRows); i++ {
		overviewRow := overviewRows[i]
		for c := overviewRow.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "td" {
				if c.FirstChild.Data == "Rasse:" {
					breed = c.NextSibling.FirstChild.Data
				}
			}
		}
	}

	return breed
}

func (p *Parser) parseSex(doc *html.Node) string {
	var breed string

	overviewRows := p.overviewRows(doc)
	if len(overviewRows) < 1 {
		return ""
	}

	for i := 0; i < len(overviewRows); i++ {
		overviewRow := overviewRows[i]
		for c := overviewRow.FirstChild; c != nil; c = c.NextSibling {
			if c.Data == "td" {
				if c.FirstChild.Data == "Geschlecht:" {
					breed = c.NextSibling.FirstChild.Data
				}
			}
		}
	}

	return breed
}

func (p *Parser) parseDescriptions(doc *html.Node) (string, string) {
	descriptionNodes := p.descriptionNodes(doc)
	if len(descriptionNodes) < 2 {
		return "", ""
	}

	shortDesc := descriptionNodes[0].FirstChild.Data
	shortDesc = strings.Replace(shortDesc, "\u000A", "", -1)
	shortDesc = strings.Replace(shortDesc, "\u00A0", " ", -1)
	shortDesc = strings.Replace(shortDesc, "\u0009", "", -1)
	shortDesc = strings.Replace(shortDesc, "  ", " ", -1)
	shortDesc = strings.Trim(shortDesc, " \n")

	var longDescBuffer bytes.Buffer
	for c := descriptionNodes[0].NextSibling; c != nil; c = c.NextSibling {
		for c2 := c.FirstChild; c2 != nil; c2 = c2.NextSibling {
			if c2.Type == html.TextNode {
				longDescBuffer.WriteString(c2.Data)
				longDescBuffer.WriteString(" ")
			}
		}
	}

	longDesc := strings.Trim(longDescBuffer.String(), " ")
	longDesc = strings.Replace(longDesc, "\u000A", "", -1)
	longDesc = strings.Replace(longDesc, "\u00A0", " ", -1)
	longDesc = strings.Replace(longDesc, "\u0009", "", -1)
	longDesc = strings.Replace(longDesc, "   ", " ", -1)
	longDesc = strings.Replace(longDesc, "  ", " ", -1)
	longDesc = strings.Trim(longDesc, " \n")

	return shortDesc, longDesc
}

func (p *Parser) paginationNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.list-view",
		"form#tx_realty_pi1_list_view",
		"table.result",
		"thead",
		"tr",
		"td",
		"table.pagination",
		"tbody",
		"tr",
		"td",
	}))
}

func (p *Parser) selector(s []string) []string {
	baseSelector := []string{
		"html",
		"body",
		"div#outerWrap",
		"div#contentWrap",
		"div#main",
		"div#mainContent",
		"div",
		"div.tx-realty-pi1",
	}

	return append(baseSelector, s...)
}

func (p *Parser) listAnimalNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.list-view",
		"form#tx_realty_pi1_list_view",
		"table.result",
		"tbody",
		"tr",
		"td.item",
		"table.item",
		"tbody",
		"tr",
		"td.description",
		"h3",
		"a",
	}))
}

func (p *Parser) overviewRows(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.single-view",
		"div",
		"div.blockall",
		"div.blockleft",
		"table.overview",
		"tbody",
		"tr",
	}))
}
func (p *Parser) nameNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.single-view",
		"div",
		"h2",
	}))
}

func (p *Parser) descriptionNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.single-view",
		"div",
		"div.blockall",
		"div.blockleft",
		"div.description",
		"p.bodytext",
	}))
}

func (p *Parser) galleryNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.single-view",
		"div",
		"div.blockall",
		"div.images",
		"table.item",
		"tbody",
		"tr",
		"td.image",
		"a",
	}))
}

func (p *Parser) imageNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.gallery-view",
		"div.thumbs",
		"table#tx_realty_thumbnailTable",
		"tbody",
		"tr",
		"td.image",
		"a",
	}))
}
