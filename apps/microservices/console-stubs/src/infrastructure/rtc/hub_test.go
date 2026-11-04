package rtc

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/jonatassales/tao/apps/microservices/console-stubs/src/application"
	"github.com/pion/transport/v5"
	"github.com/pion/transport/v5/stdnet"
	"github.com/pion/webrtc/v4"
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

// listening passes every socket through to the standard network and keeps
// the local address of each UDP socket ICE asks it for.
type listening struct {
	transport.Net
	mu        sync.Mutex
	addresses []net.IP
}

func (network *listening) record(address net.IP) {
	network.mu.Lock()
	defer network.mu.Unlock()
	network.addresses = append(network.addresses, address)
}

func (network *listening) ListenUDP(kind string, address *net.UDPAddr) (transport.UDPConn, error) {
	var ip net.IP
	if address != nil {
		ip = address.IP
	}
	network.record(ip)
	return network.Net.ListenUDP(kind, address)
}

func (network *listening) ListenPacket(kind, address string) (net.PacketConn, error) {
	host, _, _ := net.SplitHostPort(address)
	network.record(net.ParseIP(host))
	return network.Net.ListenPacket(kind, address)
}

func TestPeersListenOnLoopbackOnly(t *testing.T) {
	standard, err := stdnet.NewNet()
	if err != nil {
		t.Fatal(err)
	}
	sockets := &listening{Net: standard}
	settings := loopbackSettings()
	settings.SetNet(sockets)
	hub := &Hub{peers: make(map[*peer]struct{}), api: webrtc.NewAPI(webrtc.WithSettingEngine(settings))}
	defer hub.Close()

	client, err := webrtc.NewAPI(webrtc.WithSettingEngine(loopbackSettings())).NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	channel, err := client.CreateDataChannel("console-stubs", nil)
	if err != nil {
		t.Fatal(err)
	}
	received := make(chan string, 1)
	channel.OnMessage(func(message webrtc.DataChannelMessage) {
		select {
		case received <- string(message.Data):
		default:
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	offer, err := client.CreateOffer(nil)
	if err != nil {
		t.Fatal(err)
	}
	gathered := webrtc.GatheringCompletePromise(client)
	if err = client.SetLocalDescription(offer); err != nil {
		t.Fatal(err)
	}
	select {
	case <-ctx.Done():
		t.Fatal("the client never finished gathering")
	case <-gathered:
	}

	answer, err := hub.Offer(ctx, application.Description{Type: "offer", SDP: client.LocalDescription().SDP})
	if err != nil {
		t.Fatal(err)
	}
	if err = client.SetRemoteDescription(webrtc.SessionDescription{Type: webrtc.SDPTypeAnswer, SDP: answer.SDP}); err != nil {
		t.Fatal(err)
	}

	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for frame := ""; frame == ""; {
		select {
		case <-ctx.Done():
			t.Fatal("no frame reached the peer over loopback")
		case frame = <-received:
		case <-ticker.C:
			if err := hub.Broadcast(ctx, []byte(`{"source":"stub"}`)); err != nil {
				t.Fatal(err)
			}
		}
	}

	sockets.mu.Lock()
	defer sockets.mu.Unlock()
	if len(sockets.addresses) == 0 {
		t.Fatal("the hub opened no socket through its network")
	}
	for _, address := range sockets.addresses {
		if !address.IsLoopback() {
			t.Fatalf("a peer listened on %v", address)
		}
	}
}
