// +build unit_test

package server

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/internal/repository"
	"github.com/dimkouv/trackpal/internal/services"
)

func (ts TrackpalServer) handleDummyRequest(req int) {
	ts.trackingService.SetUser(models.UserAccount{ID: int64(req)})
	sleepMS := 50 + rand.Intn(200) // simulate some work that takes from 50 to 200 ms
	time.Sleep(time.Duration(sleepMS) * time.Millisecond)

	if int64(req) != ts.trackingService.GetUser().ID {
		panic(fmt.Sprintf("req=%d != userID=%d", req, ts.trackingService.GetUser().ID))
	}
}

func TestServiceConcurrency(t *testing.T) {
	t.Run("service properties should support concurrent writing", func(t *testing.T) {
		trackingService := services.NewTrackingService(repository.NewTrackingRepositoryMock())
		const spawnedRequests = 2000

		srv := NewTrackpalServer(trackingService)
		wg := new(sync.WaitGroup)
		wg.Add(spawnedRequests)
		for i := 0; i < spawnedRequests; i++ {
			reqID := i
			go func() {
				srv.handleDummyRequest(reqID)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}
