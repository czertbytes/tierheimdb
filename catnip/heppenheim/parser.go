package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"code.google.com/p/go.net/html"

	cp "github.com/czertbytes/tierheimdb/catnip"
	pb "github.com/czertbytes/tierheimdb/piggybank"
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
	doc, err := html.Parse(r)
	if err != nil {
		log.Println("Parse error! Pagination not found!")
		return 0, 0, 0, err
	}

	paginationNodes := p.paginationNodes(doc)
	if len(paginationNodes) != 1 {
		log.Println("Parse error! Pagination counter not found!")
		return 0, 0, 0, errors.New("parsing counter failed!")
	}

	for c := paginationNodes[0].FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && strings.HasPrefix(c.Data, "Anzahl: ") {
			counter, err := strconv.ParseInt(c.Data[8:], 10, 0)
			if err != nil {
				log.Println("Parse error! Parse pagination counter failed!")
				return 0, 0, 0, err
			}

			return 1, 10, int(counter), nil
		}
	}

	log.Println("Parse error! Pagination counter not found!")
	return 0, 0, 0, errors.New("Counter not found!")
}

func (p *Parser) ParseList(r io.Reader) ([]*pb.Animal, error) {
	var animals []*pb.Animal

	doc, err := html.Parse(r)
	if err != nil {
		return animals, err
	}

	for _, animalNode := range p.listAnimalNodes(doc) {
		name := cp.NormalizeName(p.parseName(animalNode))
		linkNodes := p.detailLinkNodes(animalNode)

		if len(name) > 0 {
			link := cp.NodeAttribute(linkNodes[0], "href")
			if strings.HasPrefix(link, "?f_mandant=tierheim_heppenheim") {
				animals = append(animals, &pb.Animal{
					Id:    cp.NormalizeId(name),
					Name:  name,
					Sex:   cp.NormalizeSex(p.parseSex(animalNode)),
					Breed: cp.NormalizeBreed(p.parseBreed(animalNode)),
					URL:   fmt.Sprintf("http://presenter.comedius.de/design/tierheim_heppenheim_standard_10001.php%s", link),
				})
			}
		}
	}

	return animals, nil
}

func (p *Parser) ParseDetail(r io.Reader) (*pb.Animal, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return &pb.Animal{
		LongDesc: p.parseLongDesc(doc),
		Images:   p.parseImages(doc),
	}, nil
}

func (p *Parser) parseName(doc *html.Node) string {
	nameNodes := p.nameNodes(doc)
	if len(nameNodes) != 1 {
		log.Println("Parse error! NameNode not found!")
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

	log.Println("Parse error! Name not found!")
	return ""
}

func (p *Parser) parseSex(doc *html.Node) string {
	var parsedStrings []string

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && len(c.Data) > 0 {
			if ps := cp.PrepareStringChunk(c.Data); len(ps) > 0 {
				parsedStrings = append(parsedStrings, ps)
			}
		}
	}

	if len(parsedStrings) > 0 {
		for _, ps := range parsedStrings {
			s := strings.ToLower(ps)
			for _, sk := range cp.SexKeywords {
				if strings.Index(s, sk) >= 0 {
					return s
				}
			}
		}
	}

	log.Println("Parse error! Sex not found!")
	return ""
}

func (p *Parser) parseBreed(doc *html.Node) string {
	var parsedStrings []string

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && len(c.Data) > 0 {
			if ps := cp.PrepareStringChunk(c.Data); len(ps) > 0 {
				parsedStrings = append(parsedStrings, ps)
			}
		}
	}

	notBreedPrefixes := []string{
		"geboren",
		"alt",
		"gesprochen",
	}

	if len(parsedStrings) > 0 {
		breed := strings.ToLower(parsedStrings[0])

		for _, nbp := range notBreedPrefixes {
			if strings.HasPrefix(breed, nbp) {
				log.Println("Parse error! Breed not found!")
				return ""
			}
		}

		return parsedStrings[0]
	}

	log.Println("Parse error! Breed not found!")
	return ""
}

func (p *Parser) parseLongDesc(doc *html.Node) string {
	detailNodes := p.longDescNodes(doc)
	if len(detailNodes) != 1 {
		log.Println("Parse error! LongDescription not found!")
		return ""
	}

	detailNode := detailNodes[0]

	longDesc := ""
	for c := detailNode.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && len(c.Data) > 0 {
			if t := cp.PrepareStringChunk(c.Data); len(t) > 0 {
				longDesc += t
				longDesc += " "
			}
		}
	}

	if len(longDesc) > 0 {
		longDesc = strings.Replace(longDesc, "  ", " ", -1)
		longDesc = strings.Replace(longDesc, " .", ".", -1)
		longDesc = strings.Trim(longDesc, " ")

		return longDesc
	}

	log.Println("Parse error! LongDescription not found!")
	return ""
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
	}))
}

func (p *Parser) nameNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"b",
	})
}

func (p *Parser) detailLinkNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, []string{
		"a",
	})
}

func (p *Parser) paginationNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"tr",
		"td",
		"div#seitenanzeigen_oben",
		"p#TextSeitenanzeige",
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
