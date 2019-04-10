package api

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/pkg/conv"
)

type httpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *httpError) MarshalJSON() ([]byte, error) {
	response, err := json.Marshal(e)
	if err != nil {
		return conv.S2B(fmt.Sprintf("{ \"code\": %d, \"message\": \"%s\" }", e.Code, e.Message)), nil
	}
	return response, nil
}

var (
	notFoundError = &httpError{
		Code:    fasthttp.StatusNotFound,
		Message: "not found",
	}
	internalError = &httpError{
		Code:    fasthttp.StatusInternalServerError,
		Message: "internal error",
	}
	methodNotAllowedError = &httpError{
		Code:    fasthttp.StatusMethodNotAllowed,
		Message: "method not allowed",
	}
)

func (r *Router) errorHandler(ctx *fasthttp.RequestCtx, e *httpError) {
	ctx.Response.SetStatusCode(e.Code)
	data, err := json.Marshal(e)
	if err != nil {
		r.log.Error("can't marshal json")
		return
	}
	ctx.SetBody(data)
}

func (r *Router) internalError(ctx *fasthttp.RequestCtx, err error) {
	r.errorHandler(ctx, internalError)
}
