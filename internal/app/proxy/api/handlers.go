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
	"github.com/prospik/places_proxy/internal/app/proxy/dal/store"
	"github.com/prospik/places_proxy/internal/app/proxy/tcp/client"
	"github.com/prospik/places_proxy/internal/app/proxy/tcp/header"
	logging "github.com/prospik/places_proxy/pkg/logger"
)

type placesHandler struct {
	client   client.Interaction
	storage  store.Storage
	endpoint string
}

// NewGraphqlHandler constructor for graphqlHandler
func NewPlacesHandler(client client.Interaction, storage store.Storage) Handler {
	return &placesHandler{
		client:   client,
		storage:  storage,
		endpoint: "https://places.aviasales.ru/v2/places.json?",
	}
}

// Places handler
func (h *placesHandler) Places(ctx *fasthttp.RequestCtx, log logging.Logger) {
	ctx.Response.Header.SetContentType(header.JSONContentType)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	key := prepareKey(ctx)
	{
		exist := h.storage.CheckKey(ctx, key)
		if exist {
			cached, _ := h.storage.GetPlaces(ctx, key)
			if cached != nil && len(cached.Data) > 0 {
				ctx.SetBody(cached.Data)
				return
			}
		}
	}

	response, err := h.client.Fetch(GET, h.endpoint+queryArgs(ctx))
	if err != nil {
		h.errorHandler(ctx, internalError)
		return
	}
	defer fasthttp.ReleaseResponse(response)

	body := response.Body()

	checksum := xxhash.New64()
	r := bytes.NewReader(body)
	_, _ = io.Copy(checksum, r)

	data := make([]dao.Types, 0)
	if err = json.Unmarshal(body, &data); err != nil {
		err = errors.WithStack(err)
		return
	}

	places := dao.ExtractTypes(data)
	h.storage.SavePlaces(context.Background(), key, checksum.Sum64(), places)

	ctx.SetBody(ctx.Request.Body())
}
