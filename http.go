package problem

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
	"sync"
)

type problemHTTPWrapper struct {
	p Problem

	contentType string
}

var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (p *problemHTTPWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	buf := getBuffer()
	defer bufferPool.Put(buf)

	h := w.Header()

	switch p.contentType {
	default:
		// Any other content type should not be handled
		return
	case problemJsonContentType:
		err := json.NewEncoder(buf).Encode(p.p)
		if err != nil {
			return
		}
		h.Set("Content-Type", problemJsonContentType)
	case problemXmlContentType:
		buf.WriteString(xml.Header)
		err := xml.NewEncoder(buf).Encode(p.p)
		if err != nil {
			return
		}
		h.Set("Content-Type", problemXmlContentType)
	}

	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(p.p.GetStatus())
	buf.WriteTo(w)
}

// ServeXML returns a Handler that serves the p argument in XML format.
//
// p MUST not be a [MapProblem], since it cannot be marshaled to XML.
//
// Headers Content-Type is set to 'application/problem+xml' and X-Content-Type-Options is set to 'nosniff';
// and finally writes the status code from p.GetStatus().
func ServeXML(p Problem) http.Handler {
	return &problemHTTPWrapper{
		p:           p,
		contentType: problemXmlContentType,
	}
}

// ServeJSON returns a Handler that serves the p argument in JSON format
//
// Headers Content-Type is set to 'application/problem+json' and X-Content-Type-Options is set to
// 'nosniff'; and finally writes the status code from p.GetStatus().
func ServeJSON(p Problem) http.Handler {
	return &problemHTTPWrapper{
		p:           p,
		contentType: problemJsonContentType,
	}
}
