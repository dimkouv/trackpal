package server

import "github.com/dimkouv/trackpal/internal/services"

type TrackpalServer struct {
	trackingService services.TrackingService
}

func NewTrackpalServer(trackingService services.TrackingService) TrackpalServer {
	return TrackpalServer{trackingService: trackingService}
}
