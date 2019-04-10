package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/prospik/places_proxy/internal/app/proxy/api"
	"github.com/prospik/places_proxy/pkg/configing"
	"github.com/prospik/places_proxy/pkg/logger"
)

func main() {
	loggerCfg := config.NewLoggerConfig()
	log := logging.New(
		logging.Level(loggerCfg.Level),
		logging.StackTraceLevel(loggerCfg.StackTraceLevel),
		logging.SentryLevel(loggerCfg.SentryLevel),
		logging.SentryDSN(loggerCfg.SentryDSN),
		logging.SentryStacktraceEnabled(loggerCfg.SentryStacktraceEnabled),
	)

	router := api.NewRouter(log.Copy(logging.Any("module", "proxy_api")))
	router.RegisterPlacesRoutes()

	server := &fasthttp.Server{
		Handler:                            router.ServeHTTP,
		Name:                               "places_proxy",
		LogAllErrors:                       false,
		DisableHeaderNamesNormalizing:      true,
		SleepWhenConcurrencyLimitsExceeded: 0,
		NoDefaultServerHeader:              true,
		TCPKeepalive:                       true,
		ReadTimeout:                        time.Duration(time.Second * 30),
		WriteTimeout:                       time.Duration(time.Second * 30),
	}

	serverCfg := config.NewServerConfig()
	go func(url string) {
		if err := server.ListenAndServe(url); err != nil {
			log.Fatal(err)
		}

	}(serverCfg.Addr)

	log.Infow("http server started", "address", serverCfg.Addr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sig
	if err := server.Shutdown(); err != nil {
		log.Fatal(err)
	}
	log.Infow("http server stopped")
	log.Infow("proxy shutting down...")
}
