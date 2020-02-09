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
	trackpalServer.RegisterRoutes()
	trackpalServer.Listen(":8080")
}
