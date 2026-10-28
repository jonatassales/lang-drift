// Package rtc owns Pion data channels and their bounded local peer lifecycle.
package rtc

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/jonatassales/tao/apps/extub/src/application"
	"github.com/jonatassales/tao/apps/extub/src/domain"
	"github.com/pion/webrtc/v4"
)

type peer struct {
	connection *webrtc.PeerConnection
	mu         sync.Mutex
	channel    *webrtc.DataChannel
	expires    time.Time
}
type Hub struct {
	mu      sync.Mutex
	peers   map[*peer]struct{}
	pending int
	closed  bool
	api     *webrtc.API
}

func New() *Hub {
	settings := webrtc.SettingEngine{}
	settings.SetIncludeLoopbackCandidate(true)
	settings.SetIPFilter(func(address net.IP) bool { return address.IsLoopback() })
	settings.SetICETimeouts(5*time.Second, 10*time.Second, time.Second)
	settings.SetSCTPMaxReceiveBufferSize(65536)
	settings.SetSCTPMaxMessageSize(262144)
	return &Hub{peers: make(map[*peer]struct{}), api: webrtc.NewAPI(webrtc.WithSettingEngine(settings))}
}

// Offer admits at most eight peers including negotiations. Failure or cancellation
// closes its connection; successful peers are owned by Broadcast/Close.
func (hub *Hub) Offer(ctx context.Context, offer application.Description) (application.Description, error) {
	if offer.Type != "offer" || len(offer.SDP) == 0 || len(offer.SDP) > 60000 {
		return application.Description{}, domain.ErrInvalid
	}
	hub.mu.Lock()
	if hub.closed || hub.pending+len(hub.peers) >= 8 {
		hub.mu.Unlock()
		return application.Description{}, domain.ErrLimit
	}
	hub.pending++
	hub.mu.Unlock()
	defer func() { hub.mu.Lock(); hub.pending--; hub.mu.Unlock() }()
	connection, err := hub.api.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return application.Description{}, err
	}
	accepted := false
	defer func() {
		if !accepted {
			_ = connection.Close()
		}
	}()
	item := &peer{connection: connection, expires: time.Now().Add(10 * time.Minute)}
	connection.OnDataChannel(func(channel *webrtc.DataChannel) {
		item.mu.Lock()
		defer item.mu.Unlock()
		if channel.Label() != "extub" || item.channel != nil {
			_ = channel.Close()
			return
		}
		item.channel = channel
	})
	if err = connection.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeOffer, SDP: offer.SDP}); err != nil {
		return application.Description{}, domain.ErrInvalid
	}
	answer, err := connection.CreateAnswer(nil)
	if err != nil {
		return application.Description{}, err
	}
	gathered := webrtc.GatheringCompletePromise(connection)
	if err = connection.SetLocalDescription(answer); err != nil {
		return application.Description{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		return application.Description{}, ctx.Err()
	case <-gathered:
	}
	final := connection.LocalDescription()
	if final == nil {
		return application.Description{}, errors.New("missing_answer")
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	if hub.closed {
		return application.Description{}, context.Canceled
	}
	hub.peers[item] = struct{}{}
	accepted = true
	return application.Description{Type: "answer", SDP: final.SDP}, nil
}

// Broadcast uses Pion's nonblocking data-channel writes. Slow consumers are closed
// before queued bytes plus the next frame exceed 256 KiB. No per-frame goroutine
// or retry survives cancellation. Peer I/O occurs outside the registry lock.
func (hub *Hub) Broadcast(ctx context.Context, frame []byte) error {
	if len(frame) > 262144 {
		return domain.ErrLimit
	}
	hub.mu.Lock()
	peers := make([]*peer, 0, len(hub.peers))
	for item := range hub.peers {
		peers = append(peers, item)
	}
	hub.mu.Unlock()
	for _, item := range peers {
		if ctx.Err() != nil {
			return nil
		}
		state := item.connection.ConnectionState()
		if time.Now().After(item.expires) || state == webrtc.PeerConnectionStateClosed || state == webrtc.PeerConnectionStateFailed || state == webrtc.PeerConnectionStateDisconnected {
			hub.remove(item)
			continue
		}
		item.mu.Lock()
		channel := item.channel
		item.mu.Unlock()
		if channel == nil || channel.ReadyState() != webrtc.DataChannelStateOpen {
			continue
		}
		if channel.BufferedAmount()+uint64(len(frame)) > 262144 {
			hub.remove(item)
			continue
		}
		if err := channel.SendText(string(frame)); err != nil {
			hub.remove(item)
		}
	}
	return nil
}
func (hub *Hub) remove(item *peer) {
	hub.mu.Lock()
	delete(hub.peers, item)
	hub.mu.Unlock()
	_ = item.connection.Close()
}

// Close is idempotent, rejects new admissions and releases every retained peer.
func (hub *Hub) Close() error {
	hub.mu.Lock()
	hub.closed = true
	peers := hub.peers
	hub.peers = make(map[*peer]struct{})
	hub.mu.Unlock()
	var result error
	for item := range peers {
		result = errors.Join(result, item.connection.Close())
	}
	return result
}
