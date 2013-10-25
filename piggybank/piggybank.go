package piggybank

type Keys []string
type Ids []string

type Shelter struct {
	Id             string   `json:"id" redis:"id"`
	Created        string   `json:"created" redis:"created"`
	Enabled        bool     `json:"enabled" redis:"enabled"`
	Name           string   `json:"name" redis:"name"`
	FullName       string   `json:"fullName" redis:"fullName"`
	URL            string   `json:"url" redis:"url"`
	LogoURL        string   `json:"logoUrl" redis:"logoUrl"`
	Phone          string   `json:"phone" redis:"phone"`
	Email          string   `json:"email" redis:"email"`
	ShortDesc      string   `json:"shortDesc" redis:"shortDesc"`
	LongDesc       string   `json:"longDesc" redis:"longDesc"`
	Street         string   `json:"street" redis:"street"`
	StreetNumber   string   `json:"streetNumber" redis:"streetNumber"`
	PostalCode     string   `json:"postalCode" redis:"postalCode"`
	City           string   `json:"city" redis:"city"`
	LatLon         string   `json:"latLon" redis:"latLon"`
	Note           string   `json:"note" redis:"note"`
	AnimalTypes    []string `json:"animalTypes" redis:"-"`
	IntAnimalTypes int      `json:"-" redis:"animalTypes"`
}

func (s *Shelter) SetIntAnimalTypes() {
	s.IntAnimalTypes = intAnimalTypes(s.AnimalTypes)
}

func (s *Shelter) SetAnimalTypes() {
	s.AnimalTypes = animalTypes(s.IntAnimalTypes)
}

func (s *Shelter) HasAnimalType(animalType string) bool {
	for _, aType := range s.AnimalTypes {
		if aType == animalType {
			return true
		}
	}

	return false
}

type Shelters []Shelter

func (ss Shelters) Len() int {
	return len(ss)
}

func (ss Shelters) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

func (ss Shelters) Paginate(p Pagination) Shelters {
	offset := p.Offset
	limit := p.Limit
	resLength := len(ss)

	if offset > resLength {
		offset = resLength
	}

	if limit > resLength {
		limit = resLength
	}

	return ss[offset:(offset + limit)]
}

type SheltersByName struct {
	Shelters
}

func (s SheltersByName) Less(i, j int) bool {
	return s.Shelters[i].Name < s.Shelters[j].Name
}

type Animal struct {
	Id        string `json:"id" redis:"id"`
	Created   string `json:"created" redis:"created"`
	Name      string `json:"name" redis:"name"`
	URL       string `json:"url" redis:"url"`
	Priority  int    `json:"priority" redis:"priority"`
	Type      string `json:"type" redis:"type"`
	Breed     string `json:"breed" redis:"breed"`
	Sex       string `json:"sex" redis:"sex"`
	ShortDesc string `json:"shortDesc" redis:"shortDesc"`
	LongDesc  string `json:"longDesc" redis:"longDesc"`
	Images    Images `json:"images" redis:"-"`
	ShelterId string `json:"shelterId" redis:"shelterId"`
	UpdateId  string `json:"updateId" redis:"updateId"`
}

type Animals []Animal

func (as Animals) Len() int {
	return len(as)
}

func (as Animals) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

func (as Animals) Paginate(p Pagination) Animals {
	offset := p.Offset
	limit := p.Limit
	resLength := len(as)

	if offset > resLength {
		offset = resLength
	}

	if limit > resLength {
		limit = resLength
	}

	return as[offset:(offset + limit)]
}

type AnimalsByName struct {
	Animals
}

func (s AnimalsByName) Less(i, j int) bool {
	return s.Animals[i].Name < s.Animals[j].Name
}

type Image struct {
	Width   int    `json:"width" redis:"width"`
	Height  int    `json:"height" redis:"height"`
	URL     string `json:"url" redis:"url"`
	Comment string `json:"comment" redis:"comment"`
}

type Images []Image

type Update struct {
	Id        string `json:"id" redis:"id"`
	Created   string `json:"created" redis:"created"`
	ShelterId string `json:"shelterId" redis:"shelterId"`
	Cats      int    `json:"cats" redis:"cats"`
	Dogs      int    `json:"dogs" redis:"dogs"`
}

type Updates []Update

func (us Updates) Len() int {
	return len(us)
}

func (us Updates) Swap(i, j int) {
	us[i], us[j] = us[j], us[i]
}

func (us Updates) Paginate(p Pagination) Updates {
	offset := p.Offset
	limit := p.Limit
	resLength := len(us)

	if offset > resLength {
		offset = resLength
	}

	if limit > resLength {
		limit = resLength
	}

	return us[offset:(offset + limit)]
}

type ByDate struct {
	Updates
}

func (s ByDate) Less(i, j int) bool {
	return s.Updates[i].Created > s.Updates[j].Created
}

type Pagination struct {
	Offset int
	Limit  int
}

var (
	maxPagination = Pagination{
		0,
		999,
	}
)
