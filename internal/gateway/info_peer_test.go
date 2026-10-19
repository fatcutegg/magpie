package gateway

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// The gateway says how it was started and which process it is to a magpie on
// this computer, whose Routing page names the one that serves the gateway
// (#110). Another machine hears neither: it routes nothing here.
func TestInfoNamesTheProcessToThisComputerOnly(t *testing.T) {
	h := New().Handler()

	for _, c := range []struct {
		name, addr string
		named      bool
	}{
		{"this computer", "127.0.0.1:50123", true},
		{"ipv6 loopback", "[::1]:50123", true},
		{"another machine", "192.168.0.9:50123", false},
	} {
		t.Run(c.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = c.addr
			h.ServeHTTP(w, r)
			if w.Code != 200 {
				t.Fatalf("%d %s", w.Code, w.Body)
			}
			var got struct {
				Pid    int    `json:"pid"`
				Serves string `json:"serves"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
				t.Fatal(err)
			}
			if named := got.Pid != 0 && got.Serves != ""; named != c.named {
				t.Fatalf("pid %d serves %q, wanted named=%v: %s", got.Pid, got.Serves, c.named, w.Body)
			}
		})
	}
}
