package server

import (
	"net/http"

	"github.com/dimkouv/trackpal/internal/services"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type TrackpalServer struct {
	trackingService services.TrackingService
	userService     services.UserAccountService
	routes          []Route
}

func NewTrackpalServer(trackingService services.TrackingService) TrackpalServer {
	return TrackpalServer{trackingService: trackingService}
}
