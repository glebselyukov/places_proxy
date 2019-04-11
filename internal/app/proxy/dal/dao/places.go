package dao

import (
	"encoding/json"
)

type Places struct {
	ID string `json:"id"`
}

func (p *Places) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *Places) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}
