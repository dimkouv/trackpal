package main

import (
	"github.com/dimkouv/trackpal/internal/server"
	"github.com/dimkouv/trackpal/internal/services"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	trackingService := services.NewTrackingServiceMock()
	trackpalServer := server.NewTrackpalServer(trackingService)
	routes := trackpalServer.RegisterRoutes()
	trackpalServer.ListenAndServe(":8080", routes)
}
