package main

import (
	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/server"
	"github.com/dimkouv/trackpal/internal/services"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2 Jan 2006 15:04:05",
	})

	trackingService, err := services.NewTrackingServicePostgres(postgresDSN)
	if err != nil {
		panic(err)
	}

	uaService, err := services.NewUserAccountServicePostgres(postgresDSN)
	if err != nil {
		panic(err)
	}

	trackpalServer := server.NewTrackpalServer(*trackingService, *uaService)
	routes := trackpalServer.RegisterRoutes()
	trackpalServer.ListenAndServe(serverAddr, routes)
}
