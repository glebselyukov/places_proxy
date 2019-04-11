package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/prospik/places_proxy/internal/app/proxy/api"
	"github.com/prospik/places_proxy/internal/app/proxy/dal"
	"github.com/prospik/places_proxy/internal/app/proxy/tcp/client"
	"github.com/prospik/places_proxy/internal/app/proxy/tcp/server"
	config "github.com/prospik/places_proxy/pkg/configing"
	logging "github.com/prospik/places_proxy/pkg/logger"
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

	dbCfg := config.NewDatabaseConfig()
	uri := dbCfg.URI
	db, err := dal.New(uri)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	scheme, err := dal.ParseDBScheme(uri)
	if err != nil {
		log.Fatal(err)
	}

	log.Copy(logging.Any("scheme", scheme)).Info("connection to DB is opened")
	defer func() {
		err = db.Disconnect()
		if err != nil {
			log.Error(err)
		}
		log.Copy(logging.Any("scheme", scheme)).Info("connection to DB is closed")
	}()

	clientCfg := config.NewClientConfig()
	httpClient := client.NewHTTPClient(clientCfg)

	routerLog := log.Copy(logging.Any("module", "proxy_api"))
	router := api.NewRouter(routerLog, httpClient)
	router.RegisterPlacesRoutes()

	serverCfg := config.NewServerConfig()
	httpServer := server.NewHTTPServer(serverCfg, router.ServeHTTP)

	go func(addr string) {
		if err := httpServer.ListenAndServe(addr); err != nil {
			log.Fatal(err)
		}
	}(serverCfg.Addr)

	log.Infow("http server started", "address", serverCfg.Addr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sig
	if err := httpServer.Shutdown(); err != nil {
		log.Fatal(err)
	}
	log.Infow("http server stopped")
	log.Infow("proxy shutting down")
}
