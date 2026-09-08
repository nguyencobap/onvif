package gosoap

import (
	"strings"
	"testing"

	"github.com/beevik/etree"
)

func TestBuildSOAP(t *testing.T) {
	namespaces := map[string]string{"trt": "urn:media", "tt": "urn:types", "unused": "urn:unused", "wsa": "http://www.w3.org/2005/08/addressing", "xsi": "http://www.w3.org/2001/XMLSchema-instance"}
	for _, body := range []string{`<trt:GetProfiles/>`, `<trt:Request xsi:type="tt:Type"><trt:Value> tt:Qualified </trt:Value></trt:Request>`} {
		got, err := BuildSOAP(body, namespaces)
		if err != nil {
			t.Fatal(err)
		}
		doc := etree.NewDocument()
		if err := doc.ReadFromString(got); err != nil {
			t.Fatal(err)
		}
		if doc.Root().NamespaceURI() != soap12Namespace || doc.FindElement("./Envelope/Header") != nil || !strings.Contains(got, body) {
			t.Fatal("operation or envelope changed")
		}
		if strings.Contains(got, "xmlns:unused") {
			t.Fatal("unused namespace attached")
		}
	}
	header := `<wsa:Action soap-env:mustUnderstand="1">urn:action</wsa:Action>`
	got, err := BuildSOAP(`<trt:Request><![CDATA[tt:Data & <value>]]></trt:Request>`, namespaces, header)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, header) || !strings.Contains(got, "<![CDATA[tt:Data & <value>]]>") || !strings.Contains(got, "xmlns:tt=") {
		t.Fatal("header, CDATA or QName lost")
	}
	for _, body := range []string{`<trt:Request xmlns:s="urn:custom"><s:Value/></trt:Request>`, `<s:Request/>`} {
		custom := map[string]string{"trt": "urn:media", "s": "urn:custom", "soap0": "urn:also-custom"}
		got, err := BuildSOAP(body, custom)
		if err != nil {
			t.Fatal(err)
		}
		doc := etree.NewDocument()
		doc.ReadFromString(got)
		if doc.Root().NamespaceURI() != soap12Namespace || doc.Root().Space == "s" || !strings.Contains(got, body) {
			t.Fatal("prefix collision changed payload")
		}
		if custom["s"] != "urn:custom" || len(custom) != 3 {
			t.Fatal("caller namespace map modified")
		}
	}
	for _, invalid := range []string{`<broken>`, `<a/><b/>`, `text<a/>`} {
		if _, err := BuildSOAP(invalid, namespaces); err == nil {
			t.Fatal("invalid body accepted")
		}
		if _, err := BuildSOAP("", namespaces, invalid); err == nil {
			t.Fatal("invalid header accepted")
		}
	}
	signedHeader := `<ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#"><ds:SignedInfo/></ds:Signature>`
	signed, err := BuildSOAP(`<trt:GetProfiles/>`, namespaces, signedHeader)
	if err != nil || !strings.Contains(signed, signedHeader) || !strings.Contains(signed, "<?xml") || !strings.Contains(signed, "xmlns:unused=") {
		t.Fatal("signed fragment lost its envelope context")
	}
}
