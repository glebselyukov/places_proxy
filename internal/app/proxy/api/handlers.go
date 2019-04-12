package api

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/internal/app/proxy/dal/store"
	"github.com/prospik/places_proxy/internal/app/proxy/tcp/client"
	"github.com/prospik/places_proxy/internal/app/proxy/tcp/header"
	config "github.com/prospik/places_proxy/pkg/configing"
	logging "github.com/prospik/places_proxy/pkg/logger"
)

type placesHandler struct {
	client    client.Interaction
	storage   store.Storage
	clientCfg *config.ClientConfig
}

// NewGraphqlHandler constructor for graphqlHandler
func NewPlacesHandler(clientCfg *config.ClientConfig, client client.Interaction, storage store.Storage) Handler {
	return &placesHandler{
		client:    client,
		storage:   storage,
		clientCfg: clientCfg,
	}
}

// Places handler
func (h *placesHandler) Places(ctx *fasthttp.RequestCtx, log logging.Logger) {
	ctx.Response.Header.SetContentType(header.JSONContentType)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	fetch := h.client.Fetch
	url := fmt.Sprintf("%s?%s", h.clientCfg.Endpoint, queryArgs(ctx))

	key := prepareKey(ctx)
	{
		exist := h.storage.CheckKey(ctx, key)
		if exist {
			cached, _ := h.storage.GetPlaces(ctx, key)
			if cached != nil && len(cached.Data) > 0 {
				ctx.SetBody(cached.Data)

				t, _ := time.Parse(layoutRFC3339, cached.Time)
				fmt.Println(time.Since(t).Minutes())
				if time.Since(t).Minutes() > h.clientCfg.UpdateMinutes {

					go func(m []byte, u string) {
						response, err := fetch(m, u)
						if err != nil {
							return
						}

						_ = h.placesPolisher(response, key)

					}(GET, url)

				}

				return
			}
		}
	}

	response, err := fetch(GET, url)
	if err != nil {
		h.errorHandler(ctx, internalError)
		return
	}
	defer fasthttp.ReleaseResponse(response)

	err = h.placesPolisher(response, key)
	if err != nil {
		h.errorHandler(ctx, internalError)
		return
	}

	ctx.SetBody(ctx.Request.Body())
}
