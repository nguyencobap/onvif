package gosoap

import (
	"fmt"
	"sort"
	"strings"

	"github.com/beevik/etree"
)

const soap12Namespace = "http://www.w3.org/2003/05/soap-envelope"

// BuildSOAP builds a minimal envelope from operation and optional header fragments.
// Namespace candidates are filtered before they are attached to the envelope.
func BuildSOAP(body string, namespaces map[string]string, headers ...string) (string, error) {
	doc := etree.NewDocument()
	root := doc.CreateElement("s:Envelope")
	addFragment := func(parent *etree.Element, fragment string) error {
		part := etree.NewDocument()
		part.ReadSettings.PreserveCData = true
		if err := part.ReadFromString(fragment); err != nil {
			return err
		}
		if len(part.ChildElements()) != 1 {
			return fmt.Errorf("expected one XML element")
		}
		for _, token := range part.Child {
			if data, ok := token.(*etree.CharData); ok && strings.TrimSpace(data.Data) != "" {
				return fmt.Errorf("unexpected text outside XML element")
			}
		}
		parent.AddChild(part.Root())
		return nil
	}
	var header *etree.Element
	for _, fragment := range headers {
		if strings.TrimSpace(fragment) == "" {
			continue
		}
		if header == nil {
			header = root.CreateElement("s:Header")
		}
		if err := addFragment(header, fragment); err != nil {
			return "", fmt.Errorf("parse SOAP header: %w", err)
		}
	}
	bodyElement := root.CreateElement("s:Body")
	if strings.TrimSpace(body) != "" {
		if err := addFragment(bodyElement, body); err != nil {
			return "", fmt.Errorf("parse SOAP body: %w", err)
		}
	}
	candidates := map[string]string{
		"soap-env": soap12Namespace,
		"soap-enc": "http://www.w3.org/2003/05/soap-encoding",
	}
	for prefix, uri := range namespaces {
		candidates[prefix] = uri
	}
	elements := soapElements(root)
	signed := false
	for _, element := range elements {
		if element.Tag != "Signature" {
			continue
		}
		ns := element.NamespaceURI()
		if ns == "" {
			ns = candidates[element.Space]
		}
		if element.Tag == "Signature" && ns == "http://www.w3.org/2000/09/xmldsig#" {
			signed = true
		}
	}
	prefix := "s"
	for suffix := 0; ; suffix++ {
		conflict := candidates[prefix] != "" && candidates[prefix] != soap12Namespace
		for _, element := range elements {
			for _, attr := range element.Attr {
				if attr.Space == "xmlns" && attr.Key == prefix && attr.Value != soap12Namespace {
					conflict = true
				}
			}
		}
		if !conflict {
			break
		}
		prefix = fmt.Sprintf("soap%d", suffix)
	}
	if signed {
		// Preserve the legacy envelope context around signed fragments.
		prefix = "soap-env"
		instruction := doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
		doc.RemoveChild(instruction)
		doc.InsertChildAt(0, instruction)
		if header == nil {
			header = etree.NewElement("soap-env:Header")
			root.InsertChildAt(0, header)
		}
	}
	root.Space = prefix
	bodyElement.Space = prefix
	if header != nil {
		header.Space = prefix
	}
	candidates[prefix] = soap12Namespace
	attrs := make([]etree.Attr, 0, len(candidates))
	for name, uri := range candidates {
		attr := etree.Attr{Space: "xmlns", Key: name, Value: uri}
		if name == "" {
			attr.Space, attr.Key = "", "xmlns"
		}
		attrs = append(attrs, attr)
	}
	if signed {
		root.Attr = attrs
	} else {
		root.Attr = requiredNamespaceAttrs(root, attrs)
	}
	return doc.WriteToString()
}

func soapElements(root *etree.Element) []*etree.Element {
	var elements []*etree.Element
	var collect func(*etree.Element)
	collect = func(element *etree.Element) {
		elements = append(elements, element)
		for _, child := range element.ChildElements() {
			collect(child)
		}
	}
	collect(root)
	return elements
}

func requiredNamespaceAttrs(root *etree.Element, candidates []etree.Attr) []etree.Attr {
	used := map[string]bool{}
	var values []string
	for _, element := range soapElements(root) {
		used[element.Space] = true
		for _, attr := range element.Attr {
			if attr.Space == "xmlns" || (attr.Space == "" && attr.Key == "xmlns") {
				continue
			}
			used[attr.Space] = true
			values = append(values, attr.Value)
		}
		for _, child := range element.Child {
			if data, ok := child.(*etree.CharData); ok {
				values = append(values, data.Data)
			}
		}
	}
	attrs := make([]etree.Attr, 0, len(candidates))
	for _, attr := range candidates {
		if attr.Space == "xmlns" && !used[attr.Key] {
			// QName and XPath prefixes can occur in values, not just XML names.
			// ponytail: O(prefixes * value bytes); tokenize if namespace counts become large.
			for _, value := range values {
				if strings.Contains(value, attr.Key+":") {
					used[attr.Key] = true
					break
				}
			}
			if !used[attr.Key] {
				continue
			}
		}
		attrs = append(attrs, attr)
	}
	sort.SliceStable(attrs, func(i, j int) bool {
		if attrs[i].Space == "xmlns" && attrs[j].Space == "xmlns" {
			if attrs[i].Key == attrs[j].Key {
				return false
			}
			if attrs[i].Key == root.Space {
				return true
			}
			if attrs[j].Key == root.Space {
				return false
			}
			return attrs[i].Key < attrs[j].Key
		}
		return attrs[i].Space == "xmlns" && attrs[j].Space != "xmlns"
	})
	return attrs
}
