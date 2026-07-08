package mcp

import (
	"net/http"
	"sync"
	"testing"
)

func TestSessionAdmissionReservationIsAtomic(t *testing.T) {
	guard := newSessionAdmissionHandler(http.NotFoundHandler(), 4).(*sessionAdmissionHandler)
	const workers = 32
	results := make(chan bool, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- guard.reserveSessionSlot()
		}()
	}
	wg.Wait()
	close(results)

	allowed := 0
	for result := range results {
		if result {
			allowed++
		}
	}
	if allowed != 4 {
		t.Fatalf("atomic admission allowed = %d, want 4", allowed)
	}
}

func TestSessionAdmissionReservationTracksAndForgetsSessions(t *testing.T) {
	guard := newSessionAdmissionHandler(http.NotFoundHandler(), 1).(*sessionAdmissionHandler)
	if !guard.reserveSessionSlot() {
		t.Fatal("initial reservation rejected")
	}
	guard.completeSessionReservation("session-a")
	if guard.reserveSessionSlot() {
		t.Fatal("reservation exceeded tracked session cap")
	}
	guard.forgetSession("session-a")
	if !guard.reserveSessionSlot() {
		t.Fatal("reservation remained blocked after session removal")
	}
}
