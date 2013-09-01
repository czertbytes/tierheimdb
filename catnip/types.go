package catnip

type Source struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Pagination string `json:"pagination"`
	Type       string `json:"type"`
	Animal     string `json:"animal"`
	Priority   int    `json:"priority"`
}
