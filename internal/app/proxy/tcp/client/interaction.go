package client

import (
	"time"

	"github.com/valyala/fasthttp"
)

type Interaction interface {
	Interacting
	Requester
	Fetcher
	Delayer
}

type Interacting interface {
	Request(method []byte, url string, opts ...RequestOpt) (request *fasthttp.Request)
	Do(request *fasthttp.Request, opts ...DoOpt) (response *fasthttp.Response, err error)
}

type Requester interface {
	Header(key, value string) RequestOpt
	Body(body []byte) RequestOpt
}

type Fetcher interface {
	Fetch(method []byte, uri string, in ...[]byte) (response *fasthttp.Response, err error)
}

type Delayer interface {
	Delay(deadline time.Duration) DoOpt
}
