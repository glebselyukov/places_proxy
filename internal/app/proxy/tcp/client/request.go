package client

import (
	"bytes"
	"time"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/internal/app/proxy/tcp/header"
)

type RequestOpt func(request *fasthttp.Request)
type DoOpt func() time.Duration

var GET = []byte("GET")

func (_ *fastHTTPClient) Request(method []byte, url string, opts ...RequestOpt) (request *fasthttp.Request) {
	request = fasthttp.AcquireRequest()
	request.SetRequestURI(url)
	request.Header.SetMethodBytes(method)

	for _, opt := range opts {
		opt(request)
	}
	return
}

func (f *fastHTTPClient) DoTimeout(request *fasthttp.Request, opts ...DoOpt) (response *fasthttp.Response, err error) {
	response = fasthttp.AcquireResponse()
	t := f.defaultDoTimeout
	for _, opt := range opts {
		t = opt()
	}

	// TODO: Gleb Selyukov - contribute and fix fasthttp package with DoTimeout
	err = f.client.DoTimeout(request, response, t)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func (f *fastHTTPClient) Do(request *fasthttp.Request, opts ...DoOpt) (response *fasthttp.Response, err error) {
	response = fasthttp.AcquireResponse()

	err = f.client.Do(request, response)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func (f *fastHTTPClient) Fetch(method []byte, uri string, in ...[]byte) (
	response *fasthttp.Response, err error) {

	opts := []RequestOpt{
		f.Header(header.ContentType, header.JSONContentType),
	}

	if bytes.Compare(method, GET) != 0 {
		if len(in) < 1 {
			err = errors.New("body is empty")
			return
		}

		opts = append(opts, f.Body(in[0]))
	}

	request := f.Request(method, uri, opts...)
	defer fasthttp.ReleaseRequest(request)

	response, err = f.Do(request)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	if response == nil {
		err = errors.New("body is empty")
		return
	}

	return
}

func (_ *fastHTTPClient) Header(key, value string) RequestOpt {
	return func(request *fasthttp.Request) {
		request.Header.Set(key, value)
	}
}

func (_ *fastHTTPClient) Body(body []byte) RequestOpt {
	return func(request *fasthttp.Request) {
		request.SetBody(body)
	}
}

func (_ *fastHTTPClient) Delay(deadline time.Duration) DoOpt {
	return func() time.Duration {
		return deadline
	}
}
