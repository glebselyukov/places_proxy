package dao

import (
	"encoding/json"
)

type Places struct {
	Slug     string `json:"slug"`
	Subtitle string `json:"subtitle"`
	Title    string `json:"title"`
}

func (p *Places) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *Places) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}
