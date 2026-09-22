package onvif

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nguyencobap/onvif/media"
	"github.com/stretchr/testify/assert"
)

func TestNewDeviceHTTPAndHTTPS(t *testing.T) {
	for _, tc := range []struct {
		name, scheme, path, auth string
	}{
		{name: "legacy host", scheme: "http"},
		{name: "HTTP URL", scheme: "http", path: "/"},
		{name: "HTTPS URL", scheme: "https"},
		{name: "HTTPS trailing slash", scheme: "https", path: "/"},
		{name: "HTTP custom endpoint", scheme: "http", path: "/custom/device?camera=1"},
		{name: "HTTPS digest custom endpoint", scheme: "https", path: "/custom/device?camera=1", auth: DigestAuth},
	} {
		t.Run(tc.name, func(t *testing.T) {
			devicePath := tc.path
			if devicePath == "" || devicePath == "/" {
				devicePath = "/onvif/device_service"
			}
			srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tc.auth == DigestAuth && r.Header.Get("Authorization") == "" {
					w.Header().Set("WWW-Authenticate", `Digest realm="camera", nonce="test", qop="auth"`)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				switch r.URL.RequestURI() {
				case devicePath:
					fmt.Fprint(w, `<Envelope><Body><GetCapabilitiesResponse><Capabilities>
<Media><XAddr>http://192.0.2.1/onvif/media_service?camera=1</XAddr></Media>
<PTZ><XAddr>http://192.0.2.1/onvif/ptz_service</XAddr></PTZ>
</Capabilities></GetCapabilitiesResponse></Body></Envelope>`)
				case "/onvif/ptz_service":
					fmt.Fprint(w, `<Envelope><Body><GetNodesResponse><PTZNode token="node1"/></GetNodesResponse></Body></Envelope>`)
				case "/onvif/media_service?camera=1":
					fmt.Fprint(w, `<Envelope><Body><GetProfilesResponse/></Body></Envelope>`)
				default:
					http.NotFound(w, r)
				}
			}))
			if tc.scheme == "https" {
				srv.StartTLS()
			} else {
				srv.Start()
			}
			defer srv.Close()
			client := srv.Client()
			client.Timeout = 2 * time.Second
			xaddr := srv.URL + tc.path
			if tc.name == "legacy host" {
				xaddr = strings.TrimPrefix(srv.URL, "http://")
			}
			dev, err := NewDevice(DeviceParams{Xaddr: xaddr, HttpClient: client, AuthMode: tc.auth, Username: "admin", Password: "password"})
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, xaddr, dev.GetDeviceParams().Xaddr)
			assert.Equal(t, srv.URL+devicePath, dev.GetEndpoint("device"))
			assert.Equal(t, srv.URL+"/onvif/media_service?camera=1", dev.GetEndpoint("media"))
			assert.Equal(t, dev.GetEndpoint("media"), dev.GetEndpoint("media2"))
			assert.Len(t, dev.GetPTZNodes(), 1)
			resp, err := dev.CallMethod(media.GetProfiles{})
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			_, err = io.Copy(io.Discard, resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestNewDeviceInvalidXaddr(t *testing.T) {
	for _, xaddr := range []string{"", "https://", "https://:443", "ftp://camera", "https://camera:bad", "https://camera/%zz", "https://user:password@camera", "https://camera/#fragment"} {
		t.Run(xaddr, func(t *testing.T) {
			dev, err := NewDevice(DeviceParams{Xaddr: xaddr})
			assert.Nil(t, dev)
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), "invalid Xaddr")
			}
		})
	}
}

func TestDeviceAddEndpointLegacyAndIPv6(t *testing.T) {
	for _, tc := range []struct{ xaddr, want string }{
		{"camera:8080", "https://camera:8080/media"},
		{"[::1]:8080", "https://[::1]:8080/media"},
		{"https://[::1]:8443/custom/device", "https://[::1]:8443/media"},
	} {
		dev := Device{params: DeviceParams{Xaddr: tc.xaddr}, endpoints: make(map[string]string)}
		dev.addEndpoint("Media", "https://192.0.2.1/media")
		assert.Equal(t, tc.want, dev.GetEndpoint("media"))
	}
}

func TestNewDeviceHTTPSRejectsUntrustedCertificate(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("untrusted TLS connection reached the SOAP handler")
	}))
	defer srv.Close()
	dev, err := NewDevice(DeviceParams{Xaddr: srv.URL, HttpClient: &http.Client{Timeout: 2 * time.Second}})
	assert.Nil(t, dev)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "certificate")
	}
}

func TestDevice_SetDeviceInfoFromScopes(t *testing.T) {
	const (
		name     = "DeviceName"
		hardware = "M9000"
	)
	scopes := []string{
		"onvif://www.onvif.org/Profile/Streaming",
		"onvif://www.onvif.org/SomethingElse/value",
		"onvif://www.onvif.org/name/" + name,
		"onvif://www.onvif.org/hardware/" + hardware,
	}
	device := Device{}
	device.SetDeviceInfoFromScopes(scopes)
	assert.Equal(t, device.info.Name, name)
	assert.Equal(t, device.info.Model, hardware)
}
