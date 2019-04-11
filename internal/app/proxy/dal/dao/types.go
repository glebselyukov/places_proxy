package dao

type Types struct {
	Type        string `json:"type"`         // detect type
	Code        string `json:"code"`         // slug -> types: country, city, airport
	Name        string `json:"name"`         // title -> types: country, city, airport
	CountryName string `json:"country_name"` // subtitle -> types: city
	CityName    string `json:"city_name"`    // subtitle -> types: airport
}

func ExtractTypes(types []Types) []Places {
	places := make([]Places, 0, len(types))
	for _, item := range types {
		place := &Places{}
		place.Slug = item.Code
		place.Title = item.Name
		switch item.Type {
		case "airport":
			place.Subtitle = item.CityName
		case "city":
			place.Subtitle = item.CountryName
		}
		places = append(places, *place)
	}

	return places
}
