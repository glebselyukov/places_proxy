package api

import (
	"bytes"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/pkg/conv"
)

var accessControlAllowHeaders = strings.Join([]string{
	acceptHeader,
	"Content-Type",
	"Content-Length",
	"Accept-Encoding",
	authorizationHeader,
	"X-CSRF-Token",
	requestHeader,
}, ", ")

func (r *Router) corsHandler(ctx *fasthttp.RequestCtx) {
	origin := ctx.Request.Header.Peek(originHeader)
	if origin == nil {
		return
	}
	allowedMethods := "HEAD, OPTIONS"
	route, ok := r.routes[conv.B2S(ctx.Path())]
	if ok {
		allowedMethods = conv.B2S(bytes.Join(route.methods, []byte(", "))) + ", " + allowedMethods
	}
	ctx.Response.Header.SetBytesV(accessControlAllowOriginHeader, origin)
	ctx.Response.Header.Set(accessControlAllowMethods, allowedMethods)
	ctx.Response.Header.Set(accessControlAllowHeadersHeader, accessControlAllowHeaders)
}
