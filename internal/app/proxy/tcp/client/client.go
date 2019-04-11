package client

import (
	"time"

	"github.com/valyala/fasthttp"

	config "github.com/prospik/places_proxy/pkg/configing"
)

type fastHTTPClient struct {
	client           *fasthttp.Client
	defaultDoTimeout time.Duration
}

func NewHTTPClient(cfg *config.ClientConfig) Interaction {
	return newFastHTTPClient(cfg)
}

func newFastHTTPClient(cfg *config.ClientConfig) *fastHTTPClient {
	rt, wt := time.Duration(cfg.ReadTimeout), time.Duration(cfg.WriteTimeout)
	return &fastHTTPClient{
		client: &fasthttp.Client{
			Name:                          cfg.Name,
			DisableHeaderNamesNormalizing: false,
			ReadTimeout:                   time.Duration(time.Second * rt),
			WriteTimeout:                  time.Duration(time.Second * wt),
		},
		defaultDoTimeout: time.Duration(time.Second * rt),
	}
}
