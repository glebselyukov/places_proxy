package server

import (
	"time"

	"github.com/valyala/fasthttp"

	config "github.com/prospik/places_proxy/pkg/configing"
)

func NewHTTPServer(cfg *config.ServerConfig, serve fasthttp.RequestHandler) *fasthttp.Server {
	rt, wt := time.Duration(cfg.ReadTimeout), time.Duration(cfg.WriteTimeout)
	return &fasthttp.Server{
		Handler:                            serve,
		Name:                               cfg.Name,
		LogAllErrors:                       false,
		DisableHeaderNamesNormalizing:      false,
		SleepWhenConcurrencyLimitsExceeded: 0,
		NoDefaultServerHeader:              true,
		TCPKeepalive:                       true,
		ReadTimeout:                        time.Duration(time.Second * rt),
		WriteTimeout:                       time.Duration(time.Second * wt),
	}
}
