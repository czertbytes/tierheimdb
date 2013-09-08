package main

import (
	"bytes"
	"io"
	"net/url"
	"regexp"
	"strings"

	"code.google.com/p/go.net/html"

	cp "github.com/czertbytes/tierheimdb/catnip"
	pb "github.com/czertbytes/tierheimdb/piggybank"
)

var (
	detailNameRE = regexp.MustCompile(`\x28.*?\x29`)
)

type Parser struct {
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

	//  parse animals from list page
	for _, animalNode := range p.listAnimalNodes(doc) {
		link := p.parseListAnimalLink(animalNode)
		if len(link) > 0 {
			animal := &pb.Animal{
				URL: link,
			}

			animals = append(animals, animal)
		}
	}

	return animals, nil
}

func (p *Parser) ParseDetail(r io.Reader) (*pb.Animal, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return &pb.Animal{}, err
	}

	name := p.parseName(doc)
	shortDesc, longDesc := p.parseDescriptions(doc)

	return &pb.Animal{
		Id:        cp.NormalizeId(name),
		Name:      name,
		ShortDesc: shortDesc,
		LongDesc:  longDesc,
		Images:    p.parseImages(doc),
	}, nil
}

func (p *Parser) parseName(doc *html.Node) string {
	var name string

	nameNodes := p.nameNodes(doc)
	if len(nameNodes) != 1 {
		return ""
	}

	name = nameNodes[0].FirstChild.Data
	name = detailNameRE.ReplaceAllString(name, "")
	name = strings.Replace(name, "  ", " ", -1)
	name = strings.Replace(name, "\u000A", "", -1)
	name = strings.Replace(name, "\u0009", "", -1)
	name = strings.Trim(name, " ")

	return name
}

func (p *Parser) parseDescriptions(doc *html.Node) (string, string) {
	var shortDesc string
	var longDescBuffer bytes.Buffer

	descPNodes := p.descriptionNodes(doc)
	if len(descPNodes) < 2 {
		return "", ""
	}

	for i, pNode := range descPNodes {
		for c := pNode.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				if i == 0 && len(shortDesc) == 0 {
					shortDesc = c.Data
				} else {
					longDescBuffer.WriteString(c.Data)
					longDescBuffer.WriteString(" ")
				}
			}
		}
	}
	longDesc := strings.Trim(longDescBuffer.String(), " ")
	longDesc = strings.Replace(longDesc, "\u000A", "", -1)
	longDesc = strings.Replace(longDesc, "\u00A0", " ", -1)
	longDesc = strings.Replace(longDesc, "\u0009", "", -1)
	longDesc = strings.Replace(longDesc, "  ", " ", -1)

	return shortDesc, longDesc
}

func (p *Parser) parseImages(doc *html.Node) []pb.Image {
	var images []pb.Image

	for _, imageNode := range p.imageNodes(doc) {
		images = append(images, pb.Image{
			URL: cp.NodeAttribute(imageNode, "href"),
		})
	}

	if len(images) == 0 {
		for _, imageNode := range p.singleImageNodes(doc) {
			fullPath := cp.NodeAttribute(imageNode, "src")
			i := strings.Index(fullPath, "?")
			vals, err := url.ParseQuery(fullPath[i+1:])
			if err != nil {
				return images
			}

			images = append(images, pb.Image{URL: vals.Get("src")})
		}
	}

	return images
}

func (p *Parser) listAnimalNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.content_section",
		"ul#portfolio_style2",
		"li",
	}))
}

func (p *Parser) parseListAnimalLink(node *html.Node) string {
	for _, linkNode := range p.detailLinkNodes(node) {
		l := cp.NodeAttribute(linkNode, "href")
		if p.isDetailLink(l) {
			return l
		}
	}

	return ""
}

func (p *Parser) isDetailLink(link string) bool {
	notDetailLinks := []string{
		"http://www.samtpfoten-neukoelln.com/portfolio/quisque",
		"http://www.samtpfoten-neukoelln.com/portfolio/nullam-volutpat",
	}

	for _, notDetailLink := range notDetailLinks {
		if notDetailLink == link {
			return false
		}
	}

	return true
}

func (p *Parser) detailLinkNodes(node *html.Node) []*html.Node {
	return cp.NodeSelect(node, []string{
		"div.one_half_last",
		"div.portfolio_item_content",
		"h4",
		"a",
	})
}

func (p *Parser) selector(s []string) []string {
	baseSelector := []string{
		"html",
		"body",
		"div#page_wrap",
		"div.container_wrapper",
		"div.container",
	}

	return append(baseSelector, s...)
}

func (p *Parser) nameNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, []string{
		"html",
		"body",
		"div#page_wrap",
		"div#sub_header_wrapper",
		"div.container",
		"div#sub_header",
		"h1",
	})
}

func (p *Parser) descriptionNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.content_section",
		"div.portfolio",
		"p",
	}))
}

func (p *Parser) imageNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.content_section",
		"div.portfolio",
		"div.gallery",
		"dl",
		"dt",
		"a",
	}))
}

func (p *Parser) singleImageNodes(doc *html.Node) []*html.Node {
	return cp.NodeSelect(doc, p.selector([]string{
		"div.content_section",
		"div.portfolio",
		"div.blog_single_img",
		"div.sudo_slider_wrapper",
		"div.flexslider",
		"ul.slides",
		"li",
		"img",
	}))
}
