package dao

import (
	"encoding/json"
)

type ResultData struct {
	Data     []byte
	Checksum uint64
	Time     string
}

type CachedData struct {
	Data []Places `json:"data"`
}

func (p *CachedData) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *CachedData) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}
