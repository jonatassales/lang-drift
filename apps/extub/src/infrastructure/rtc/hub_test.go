package rtc

import (
	"context"
	"github.com/jonatassales/tao/apps/extub/src/application"
	"testing"
)

func TestRejectedOffersReleaseAdmissionAndCloseIsIdempotent(t *testing.T) {
	hub := New()
	defer hub.Close()
	for i := 0; i < 10; i++ {
		if _, err := hub.Offer(context.Background(), application.Description{Type: "offer", SDP: "invalid"}); err == nil {
			t.Fatal("invalid SDP accepted")
		}
	}
	hub.mu.Lock()
	pending, retained := hub.pending, len(hub.peers)
	hub.mu.Unlock()
	if pending != 0 || retained != 0 {
		t.Fatal("failed offer retained resources")
	}
	if err := hub.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := hub.Offer(context.Background(), application.Description{Type: "offer", SDP: "invalid"}); err == nil {
		t.Fatal("closed hub admitted peer")
	}
}
