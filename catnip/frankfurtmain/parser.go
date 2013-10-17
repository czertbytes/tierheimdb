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

			return 1, 5, int(counter), nil
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
		if strings.HasPrefix(link, "?f_mandant=tsv") {
			animals = append(animals, &pb.Animal{
				URL: fmt.Sprintf("http://presenter.comedius.de/design/tsv_frankfurt_standard_10001.php%s", link),
			})
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
		Id:       cp.NormalizeId(name),
		Name:     name,
		Sex:      cp.NormalizeSex(p.parseSex(doc)),
		Breed:    cp.NormalizeBreed(p.parseBreed(doc)),
		LongDesc: p.parseLongDesc(doc),
		Images:   p.parseImages(doc),
	}, nil
}

func (p *Parser) parseName(doc *html.Node) string {
	nameNodes := p.nameNodes(doc)
	if len(nameNodes) != 1 {
		return ""
	}

	if nameNodes[0].FirstChild != nil {
		name := nameNodes[0].FirstChild.Data
		if i := strings.Index(name, ","); i > 0 {
			name = name[0:i]
		}
		name = strings.Trim(name, " ")

		return cp.ToUTF8(name)
	}

	return ""
}

func (p *Parser) parseSex(doc *html.Node) string {
	detailNodes := p.longDescNodes(doc)
	if len(detailNodes) != 1 {
		return ""
	}
	detailNode := detailNodes[0]

	counter := 0
	for c := detailNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			if len(c.Data) > 0 {
				sex := cp.PrepareStringChunk(c.Data)
				if counter == 2 {
					return sex
				}

				counter += 1
			}
		}
	}

	return ""
}

func (p *Parser) parseBreed(doc *html.Node) string {
	detailNodes := p.longDescNodes(doc)
	if len(detailNodes) != 1 {
		return ""
	}
	detailNode := detailNodes[0]

	for c := detailNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			breed := cp.PrepareStringChunk(c.Data)
			if len(breed) > 0 {
				return breed
			}
		}
	}

	return ""
}

func (p *Parser) parseLongDesc(doc *html.Node) string {
	detailNodes := p.longDescNodes(doc)
	if len(detailNodes) != 1 {
		return ""
	}
	detailNode := detailNodes[0]

	var longDesc string
	counter := 0
	for c := detailNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			if counter > 3 {
				longDesc += cp.PrepareStringChunk(c.Data)
				longDesc += " "
			}
			counter += 1
		}
	}
	longDesc = strings.Replace(longDesc, "  ", " ", -1)
	longDesc = strings.Replace(longDesc, " .", ".", -1)
	longDesc = strings.Trim(longDesc, " ")

	return longDesc
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