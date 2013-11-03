package main

import (
	"fmt"
	"io"
	"log"
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
	panic("Not supported!")
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

		if len(name) > 0 && len(linkNodes) == 1 {
			link := cp.NodeAttribute(linkNodes[0], "href")
			animals = append(animals, &pb.Animal{
				Id:       cp.NormalizeId(name),
				Name:     name,
				LongDesc: p.parseLongDesc(animalNode),
				URL:      fmt.Sprintf("http://www.tierheim-verlorenwasser.de%s", link),
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

	return &pb.Animal{
		Images: p.parseImages(doc),
	}, nil
}

func (p *Parser) parseName(doc *html.Node) string {
	nameNodes := p.nameNodes(doc)
	if len(nameNodes) != 1 {
		return ""
	}

	if nameNodes[0].FirstChild != nil {
		name := nameNodes[0].FirstChild.Data

		return cp.PrepareStringChunk(name)
	}

	return ""
}

func (p *Parser) parseLongDesc(doc *html.Node) string {
	var longDesc string

	for _, longDescNode := range p.longDescNodes(doc) {

		for c := longDescNode.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				if strings.HasPrefix(c.Data, "Weitere Infos folgen") == false {
					longDesc += cp.PrepareStringChunk(c.Data)
					longDesc += " "
				}
			}
		}
	}

	return cp.PrepareStringChunk(longDesc)
}

func (p *Parser) parseImages(doc *html.Node) []pb.Image {
	var images []pb.Image

	for _, imageNode := range p.imageNodes(doc) {
		src := cp.NodeAttribute(imageNode, "href")

		images = append(images, pb.Image{
			URL: src,
		})
	}

	return images
}

func (p *Parser) selector(s []string) []string {
	baseSelector := []string{
		"html",
		"body",
		"div#ja-wrapper",
		"div#ja-containerwrap-fr",
		"div#ja-containerwrap2",
		"div#ja-container",
		"div#ja-container2",
		"div#ja-mainbody-fr",
		"div#ja-contentwrap",
		"div#ja-content",
	}

	return append(baseSelector, s...)
}

func (p *Parser) listAnimalNodes(node *html.Node) []*html.Node {
	var detailNodes []*html.Node

	top := cp.NodeSelect(node, p.selector([]string{
		"table.blog",
		"tbody",
		"tr",
		"td",
		"div",
		"div.contentpaneopen",
	}))

	rest := cp.NodeSelect(node, p.selector([]string{
		"table.blog",
		"tbody",
		"tr",
		"td",
		"table",
		"tbody",
		"tr",
		"td",
		"div.contentpaneopen",
	}))

	detailNodes = append(detailNodes, top...)

	return append(detailNodes, rest...)
}

func (p *Parser) nameNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"h2.contentheading",
	})
}

func (p *Parser) detailLinkNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.article-content",
		"p",
		"a",
	})
}

func (p *Parser) longDescNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.article-content",
		"p",
	})
}

func (p *Parser) imageNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, p.selector([]string{
		"div.article-content",
		"ul.sig-container",
		"li.sig-block",
		"span.sig-link-wrapper",
		"span.sig-link-innerwrapper",
		"a.sig-link",
	}))
}
