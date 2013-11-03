package main

import (
	"encoding/xml"
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
	var rss cp.RSS

	d := xml.NewDecoder(r)
	if err := d.Decode(&rss); err != nil {
		return nil, err
	}

	for _, i := range rss.Channel.Items {
		name := cp.NormalizeName(i.Title)
		animals = append(animals, &pb.Animal{
			Id:       cp.NormalizeId(name),
			Name:     name,
			LongDesc: p.parseLongDesc(i.Description),
			URL:      i.Link,
		})
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

func (p *Parser) parseLongDesc(text string) string {
	var longDesc string

	text = cp.PrepareStringChunk(text)
	text = strings.Replace(text, "<p>", "", -1)
	for _, b := range strings.SplitAfter(text, ">") {
		if len(b) > 0 && b[0] != '<' {
			longDescLine := strings.Replace(b, "</p>", "", -1)
			if len(longDescLine) > 0 && (strings.HasPrefix(longDescLine, "Weitere Infos folgen") == false) {
				longDesc += longDescLine
				longDesc += " "
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
		"div#containerwrap2",
		"div#ja-container",
		"div#ja-container2",
		"div#ja-mainbody-fr",
		"div#ja-contentwrap",
		"div#ja-content",
		"div.article-content",
	}

	return append(baseSelector, s...)
}

func (p *Parser) imageNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, p.selector([]string{
		"ul.sig-container",
		"li.sig-block",
		"span.sig-link-wrapper",
		"span.sig-link-innerwrapper",
		"a.sig-link",
	}))
}
