package main

import (
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/server"
	"github.com/dimkouv/trackpal/internal/services"
	"github.com/dimkouv/trackpal/pkg/mailutils"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2 Jan 2006 15:04:05",
	})

	mailSender := mailutils.NewPlainMailSender(mailutils.SMTPSettings{
		Port: func() int {
			p, err := strconv.Atoi(smtpPort)
			if err != nil {
				panic(err)
			}
			return p
		}(),
		Host:     smtpHost,
		User:     smtpUser,
		Password: smtpPassword,
	})

	trackingService, err := services.NewTrackingServicePostgres(postgresDSN, mailSender)
	if err != nil {
		panic(err)
	}

	uaService, err := services.NewUserAccountServicePostgres(postgresDSN, mailSender)
	if err != nil {
		panic(err)
	}

	trackpalServer := server.NewTrackpalServer(trackingService, uaService)
	routes := trackpalServer.RegisterRoutes()
	trackpalServer.ListenAndServe(serverAddr, routes)
}
