package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/OneOfOne/xxhash"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/internal/app/proxy/dal/dao"
)

const layoutRFC3339 = "2006-01-02T15:04:05.999999999Z07:00"

func (h *placesHandler) placesPolisher(res *fasthttp.Response, key string) error {
	body := res.Body()

	checksum := xxhash.New64()
	r := bytes.NewReader(body)
	_, _ = io.Copy(checksum, r)

	data := make([]dao.Types, 0)
	if err := json.Unmarshal(body, &data); err != nil {
		return errors.WithStack(err)
	}

	places := dao.ExtractTypes(data)
	h.storage.SavePlaces(context.Background(), key, checksum.Sum64(), places)

	fasthttp.ReleaseResponse(res)
	return nil
}
