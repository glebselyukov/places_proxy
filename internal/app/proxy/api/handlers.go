package api

import (
	"github.com/valyala/fasthttp"

	logging "github.com/prospik/places_proxy/pkg/logger"
)

type placesHandler struct{}

// NewGraphqlHandler constructor for graphqlHandler
func NewPlacesHandler() Handler {
	return &placesHandler{}
}

func (h *placesHandler) Places(ctx *fasthttp.RequestCtx, log logging.Logger) {
	ctx.Response.Header.SetContentType(jsonContentType)
	ctx.Response.SetStatusCode(fasthttp.StatusOK)

	ctx.SetBody(ctx.Request.Body())
}
