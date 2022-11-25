package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	apikey "restapi/internal/apiKey"
	apikeyDb "restapi/internal/apiKey/db"
	"restapi/internal/config"
	"restapi/pkg/client/postgresql"
	"restapi/pkg/logging"
)

func main() {

	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.DBConf)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	repository := apikeyDb.NewRepository(postgreSQLClient, cfg, logger)

	logger.Info("register handler")
	authorHandler := apikey.NewHandler(cfg, repository, logger)
	authorHandler.Register(router)
	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("listen tcp")
	listener, listenErr := net.Listen("tcp", fmt.Sprintf("%s", cfg.NetworkConf.BindAddr))
	logger.Infof("server is listening port %s", cfg.NetworkConf.BindAddr)

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10,
		ReadTimeout:  10,
	}
	logger.Fatal(server.Serve(listener))

}
