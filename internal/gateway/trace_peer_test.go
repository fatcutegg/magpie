package gateway

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// The routing is this process's own, so a magpie on this computer that does
// not serve the gateway reads it here to show the same Routing page (#110).
// It answers this computer alone: the trace names accounts and what they
// served, and the magpie that reads it is on the same machine.
func TestGatewayTraceIsLoopbackOnly(t *testing.T) {
	h := New().Handler()

	t.Run("this computer", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/magpie/trace", nil)
		r.RemoteAddr = "127.0.0.1:50123"
		h.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("%d %s", w.Code, w.Body)
		}
		var got struct {
			Object string            `json:"object"`
			Seq    int64             `json:"seq"`
			Routes []json.RawMessage `json:"routes"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if got.Object != "trace" || got.Routes == nil {
			t.Fatalf("object %q routes %v", got.Object, got.Routes)
		}
	})

	t.Run("another machine", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/magpie/trace", nil)
		r.RemoteAddr = "192.168.0.9:50123"
		h.ServeHTTP(w, r)
		if w.Code != 403 {
			t.Fatalf("another machine reached the trace: %d %s", w.Code, w.Body)
		}
	})
}
