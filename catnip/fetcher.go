package catnip

import (
	"fmt"
	"net/http"

	pb "github.com/czertbytes/tierheimdb/piggybank"
)

func PaginatedSources(p Parser, s *Source) ([]*Source, error) {
	if len(s.Pagination) == 0 {
		return []*Source{s}, nil
	}

	start, perPage, total, err := ParsePagination(p, s.URL)
	if err != nil {
		return nil, err
	}

	sources := []*Source{}
	pages := (total / perPage) + start
	for pageId := start; pageId <= pages; pageId++ {
		sources = append(sources, &Source{
			URL:    fmt.Sprintf(s.Pagination, pageId),
			Type:   s.Type,
			Animal: s.Animal,
		})
	}

	return sources, nil
}

func ParsePagination(p Parser, url string) (int, int, int, error) {
	response, err := http.Get(url)
	if err != nil {
		return 0, 0, 0, err
	}
	defer response.Body.Close()

	return p.ParsePagination(response.Body)
}

func ParseList(p Parser, url string) ([]*pb.Animal, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return p.ParseList(response.Body)
}

func ParseDetail(p Parser, url string) (*pb.Animal, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return p.ParseDetail(response.Body)
}
