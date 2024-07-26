package main

import (
	"go-service-template/config"
	"go-service-template/db"
	"go-service-template/http"
	"go-service-template/services"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v10"
	"github.com/sirupsen/logrus"
)

const (
	AppVersion = "1.0"
)

func main() {
	logrus.Info("Application init")
	logrus.Infof("Application version: %s", AppVersion)

	cfg := &config.MainConfig{}
	if err := env.Parse(cfg); err != nil {
		logrus.Fatal("failed to parse config:", err)
	}

	lvl, parseErr := logrus.ParseLevel(cfg.AppConfig.LogLvl)
	if parseErr != nil {
		lvl = logrus.InfoLevel
	}
	logrus.SetLevel(lvl)
	if lvl >= logrus.DebugLevel {
		logrus.Infof("Config: %+v", cfg)
		logrus.SetReportCaller(true) // enable file name and line number in log
	}

	logrus.Infof("App log level: %s", lvl)

	logrus.SetFormatter(&logrus.TextFormatter{DisableQuote: true})

	db, err := db.NewDBConnection(cfg.DbConfig)
	if err != nil {
		logrus.Fatal("failed to init db connection:", err)
	}

	someService, err := services.NewSomeService(cfg, db)
	if err != nil {
		logrus.Fatal("failed to init iot service:", err)
	}

	server, err := http.NewServer(&cfg.ServerConfig, someService)
	if err != nil {
		logrus.Fatal("failed to init http server:", err)
	}

	go func() {
		if err := server.Start(); err != nil {
			logrus.Fatalf("Error on starting http server: %s", err)
		}
	}()

	logrus.Info("Wait for SIGINT or SIGTERM to stop the program...")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	logrus.Info("Stopping the app...")

	server.Shutdown()

	logrus.Info("Stopped.")
}
