package wsdiscovery

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDevicesFromProbeResponsesPreservesEndpoint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RequestURI() != "/custom/device?camera=1" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprint(w, `<Envelope><Body><GetCapabilitiesResponse><Capabilities/></GetCapabilitiesResponse></Body></Envelope>`)
	}))
	defer srv.Close()
	xaddr := srv.URL + "/custom/device?camera=1"
	probe := fmt.Sprintf(`<Envelope><Body><ProbeMatches><ProbeMatch><XAddrs>%s</XAddrs></ProbeMatch></ProbeMatches></Body></Envelope>`, xaddr)
	devices, err := DevicesFromProbeResponses([]string{probe, probe})
	if err != nil {
		t.Fatal(err)
	}
	if len(devices) != 1 {
		t.Fatalf("got %d devices, want 1", len(devices))
	}
	if got := devices[0].GetEndpoint("device"); got != xaddr {
		t.Fatalf("got endpoint %q, want %q", got, xaddr)
	}
}
