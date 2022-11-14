package main

import (
	"Avito-Internship-Task/configs"
	"Avito-Internship-Task/internal/app/balance_service_app/server"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var configPath string = "configs/config.toml"

func main() {
	config := configs.CreateConfigForServer()
	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		logrus.Fatal(err)
	}
	contextLogger := logrus.WithFields(logrus.Fields{})
	logrus.SetReportCaller(false)
	logrus.SetFormatter(&logrus.TextFormatter{PadLevelText: false, DisableLevelTruncation: false})
	appServer := server.CreateServer(config, contextLogger)

	err = appServer.Start()
	if err != nil {
		panic(err)
	}
}
