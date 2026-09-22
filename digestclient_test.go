package onvif

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestDigestNonceCount(t *testing.T) {
	for _, tc := range []struct {
		count  uint32
		want   string
		opaque string
	}{{0, "00000001", "e2a56563"}, {9, "0000000a", ""}} {
		client := NewDigestClient(nil, "Mufasa", "Circle Of Life")
		challenge := `Digest algorithm=MD5, realm="testrealm@host.com", qop="auth", nonce="dcd98b7102dd2f0e8b11d0f600bfb0c093"`
		if tc.opaque != "" {
			challenge += fmt.Sprintf(`, opaque="%s"`, tc.opaque)
		}
		response := &http.Response{Header: make(http.Header)}
		response.Header.Set("WWW-Authenticate", challenge)
		client.getDigestParts(response)
		client.nonceCount = tc.count
		header, err := client.getDigestAuth("GET", "/dir/index.html")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(header, "nc="+tc.want+",") || !strings.Contains(header, "qop=auth,") {
			t.Fatalf("invalid digest nonce count or qop: %s", header)
		}
		if tc.opaque != "" && !strings.Contains(header, fmt.Sprintf(`opaque="%s"`, tc.opaque)) {
			t.Fatal("server opaque value was not echoed")
		}
		if tc.opaque == "" && strings.Contains(header, "opaque=") {
			t.Fatal("opaque was sent without a server challenge")
		}
		cnonce := regexp.MustCompile(`cnonce="([0-9a-f]+)"`).FindStringSubmatch(header)
		if len(cnonce) != 2 {
			t.Fatal("missing cnonce")
		}
		// RFC 2617 example HA1/HA2, with this request's generated cnonce.
		digest := md5.Sum([]byte("939e7578ed9e3c518a452acee763bce9:" + client.snonce + ":" + tc.want + ":" + cnonce[1] + ":auth:39aff3a2bab6126f332b942af96d3366"))
		if !strings.Contains(header, fmt.Sprintf(`response="%x"`, digest)) {
			t.Fatal("response hash did not use the eight-digit hexadecimal nonce count")
		}
	}
}
