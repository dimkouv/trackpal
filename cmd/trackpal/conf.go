package main

import "github.com/dimkouv/trackpal/internal/envlib"

// nolint: gochecknoglobals
var (
	postgresDSN = envlib.GetEnvOrDefault(
		"TRACKPAL_POSTGRES_DSN",
		"user=master password=masterkey dbname=trackpal sslmode=disable",
	)

	serverAddr = envlib.GetEnvOrDefault(
		"TRACKPAL_SERVER_ADDR",
		"127.0.0.1:8080",
	)
)
