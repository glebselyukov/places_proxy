package api

import (
	"bytes"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/internal/app/proxy/tcp/header"
	"github.com/prospik/places_proxy/pkg/conv"
)

var accessControlAllowHeaders = strings.Join([]string{
	header.AcceptHeader,
	header.ContentType,
	header.ContentLength,
	header.AcceptEncoding,
	header.AuthorizationHeader,
	header.XCSRFToken,
	header.RequestHeader,
}, ", ")

func (r *Router) corsHandler(ctx *fasthttp.RequestCtx) {
	origin := ctx.Request.Header.Peek(header.OriginHeader)
	if origin == nil {
		return
	}
	allowedMethods := "HEAD, OPTIONS"
	route, ok := r.routes[conv.B2S(ctx.Path())]
	if ok {
		allowedMethods = conv.B2S(bytes.Join(route.methods, []byte(", "))) + ", " + allowedMethods
	}
	ctx.Response.Header.SetBytesV(header.AccessControlAllowOriginHeader, origin)
	ctx.Response.Header.Set(header.AccessControlAllowMethods, allowedMethods)
	ctx.Response.Header.Set(header.AccessControlAllowHeadersHeader, accessControlAllowHeaders)
}
