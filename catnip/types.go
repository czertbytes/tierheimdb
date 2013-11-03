package catnip

type Source struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Pagination string `json:"pagination"`
	Type       string `json:"type"`
	Animal     string `json:"animal"`
	Priority   int    `json:"priority"`
}

type RSS struct {
	Channel *Channel `xml:"channel"`
}

type Channel struct {
	Items []*Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
}
