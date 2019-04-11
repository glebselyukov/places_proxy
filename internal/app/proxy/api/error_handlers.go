package api

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/pkg/conv"
)

type httpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *httpError) MarshalJSON() ([]byte, error) {
	return conv.S2B(fmt.Sprintf("{ \"code\": %d, \"message\": \"%s\" }", e.Code, e.Message)), nil
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

func errorRelease(ctx *fasthttp.RequestCtx, e *httpError) (err error) {
	ctx.Response.SetStatusCode(e.Code)
	data, err := json.Marshal(e)
	if err != nil {
		err = errors.New("can't marshal json")
		return
	}
	ctx.SetBody(data)
	return
}

func (_ *placesHandler) errorHandler(ctx *fasthttp.RequestCtx, e *httpError) {
	_ = errorRelease(ctx, e)
}

func (r *Router) errorHandler(ctx *fasthttp.RequestCtx, e *httpError) {
	err := errorRelease(ctx, e)
	if err != nil {
		r.log.Error(err.Error())
	}
}

func (r *Router) internalError(ctx *fasthttp.RequestCtx, err error) {
	r.errorHandler(ctx, internalError)
}
