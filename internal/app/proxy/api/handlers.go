package api

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/internal/app/proxy/dal/dao"
	"github.com/prospik/places_proxy/internal/app/proxy/tcp/client"
	"github.com/prospik/places_proxy/internal/app/proxy/tcp/header"
	logging "github.com/prospik/places_proxy/pkg/logger"
)

type placesHandler struct {
	client   client.Interaction
	endpoint string
}

// NewGraphqlHandler constructor for graphqlHandler
func NewPlacesHandler(client client.Interaction) Handler {
	return &placesHandler{
		client:   client,
		endpoint: "https://places.aviasales.ru/v2/places.json?",
	}
}

// Places handler
func (h *placesHandler) Places(ctx *fasthttp.RequestCtx, log logging.Logger) {
	ctx.Response.Header.SetContentType(header.JSONContentType)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	response, err := h.client.Fetch(GET, h.endpoint+queryArgs(ctx))
	if err != nil {
		h.errorHandler(ctx, internalError)
		return
	}
	defer fasthttp.ReleaseResponse(response)

	data := make([]dao.Types, 0)
	if err = json.Unmarshal(response.Body(), &data); err != nil {
		err = errors.WithStack(err)
		return
	}

	// for developing
	//s, _ := json.MarshalIndent(dao.ExtractTypes(data), " ", "\t")
	//fmt.Println(string(s))
	//
	//_ = prepareKey(ctx)

	ctx.SetBody(ctx.Request.Body())
}
