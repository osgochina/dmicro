package mdns

import "testing"

func TestServer_StartStop(t *testing.T) {
	s := makeService(t)
	serv, err := NewServer(&Config{Zone: s, LocalhostChecking: true})
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	defer serv.Shutdown()
}
