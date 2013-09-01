package main

import (
	"errors"
	"fmt"
	"io"
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
	if len(paginationNodes) != 1 {
		return 0, 0, 0, errors.New("parsing counter failed!")
	}

	for c := paginationNodes[0].FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && strings.HasPrefix(c.Data, "Anzahl: ") {
			counter, err := strconv.ParseInt(c.Data[8:], 10, 0)
			if err != nil {
				return 0, 0, 0, err
			}

			return 1, 8, int(counter), nil
		}
	}

	return 0, 0, 0, errors.New("Counter not found!")
}

func (p *Parser) ParseList(r io.Reader) ([]*pb.Animal, error) {
	var animals []*pb.Animal

	doc, err := html.Parse(r)
	if err != nil {
		return animals, err
	}

	for _, animalNode := range p.listAnimalNodes(doc) {
		link := cp.NodeAttribute(animalNode, "href")
		sex := ""
		shortDesc := ""
		if strings.HasPrefix(link, "?f_mandant=bmt_hamburg_1fb0bd784c03ad8c500a2c224deb22b5&") {
			animal := &pb.Animal{
				URL:       fmt.Sprintf("http://presenter.comedius.de/design/bmt_hamburg_standard_10001.php%s", link),
				Sex:       sex,
				ShortDesc: shortDesc,
			}
			animals = append(animals, animal)
		}
	}

	return animals, nil
}

func (p *Parser) ParseDetail(r io.Reader) (*pb.Animal, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	name := p.parseName(doc)
	return &pb.Animal{
		Id:       cp.NormalizeName(name),
		Name:     name,
		LongDesc: p.parseLongDesc(doc),
		Images:   p.parseImages(doc),
	}, nil
}

func (p *Parser) parseName(doc *html.Node) string {
	nameNodes := p.nameNodes(doc)
	if len(nameNodes) != 1 {
		return ""
	}

	name := nameNodes[0].FirstChild.Data
	if i := strings.Index(name, ","); i > 0 {
		name = name[0:i]
	}
	name = strings.Trim(name, " ")

	return cp.ToUTF8(name)
}

func (p *Parser) parseLongDesc(doc *html.Node) string {
	detailNodes := p.longDescNodes(doc)
	if len(detailNodes) != 1 {
		return ""
	}
	detailNode := detailNodes[0]

	var longDesc string
	for c := detailNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			longDesc += c.Data
		}
	}
	longDesc = strings.Replace(longDesc, "\u0009", "", -1)
	longDesc = strings.Replace(longDesc, "\u000A", "", -1)
	longDesc = strings.Trim(longDesc, " ")

	return cp.ToUTF8(longDesc)
}

func (p *Parser) parseImages(doc *html.Node) []pb.Image {
	var images []pb.Image

	for _, imageNode := range p.imageNodes(doc) {
		var width int
		widthStr := cp.NodeAttribute(imageNode, "width")
		widthInt, err := strconv.ParseInt(widthStr, 10, 0)
		if err == nil {
			width = int(widthInt)
		}

		src := cp.NodeAttribute(imageNode, "src")
		src = fmt.Sprintf("http://presenter.comedius.de%s", src[6:])

		images = append(images, pb.Image{
			URL:   src,
			Width: width,
		})
	}

	return images
}

func (p *Parser) selector(s []string) []string {
	baseSelector := []string{
		"html",
		"body",
		"table",
		"tbody",
	}

	return append(baseSelector, s...)
}

func (p *Parser) listAnimalNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"tr",
		"td",
		"table",
		"tbody",
		"tr",
		"td",
		"span",
		"a",
	}))
}

func (p *Parser) paginationNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"tr",
		"td",
		"div#seitenanzeigen_oben",
		"p#TextSeitenanzeige",
	}))
}

func (p *Parser) nameNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, p.selector([]string{
		"tr",
		"td",
		"table",
		"tbody",
		"tr",
		"td",
		"p",
		"b",
	}))
}

func (p *Parser) longDescNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, p.selector([]string{
		"tr",
		"td",
		"table",
		"tbody",
		"tr",
		"td",
		"p",
	}))
}

func (p *Parser) imageNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, p.selector([]string{
		"tr",
		"td",
		"span",
		"p",
		"img",
	}))
}
