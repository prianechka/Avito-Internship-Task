package main

import (
	"Avito-Internship-Task/internal/app/balance_service_app/server"
	"flag"
	"github.com/sirupsen/logrus"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
}

// @title BalanceApp
// @version 1.0
// @description Server for Balance application.

// @BasePath /api/v1

// @x-extension-openapi {"example": "value on a json format"}
func main() {
	logger := logrus.Logger{}
	appServer := server.CreateServer(&logger)
	err := appServer.Start()
	if err != nil {
		panic(err)
	}
}
