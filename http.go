package problem

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
)

type problemHTTPWrapper struct {
	p Problem

	contentType string
}

func (p *problemHTTPWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	h := w.Header()

	switch p.contentType {
	case problemJsonContentType:
		b, err := json.Marshal(p.p)
		if err != nil {
			return
		}
		h.Set("Content-Type", problemJsonContentType)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Content-Length", strconv.Itoa(len(b)))
		w.WriteHeader(p.p.GetStatus())
		w.Write(b)
	case problemXmlContentType:
		b, err := xml.Marshal(p.p)
		if err != nil {
			return
		}
		h.Set("Content-Type", problemXmlContentType)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Content-Length", strconv.Itoa(len(xml.Header)+len(b)))
		w.WriteHeader(p.p.GetStatus())
		w.Write(append([]byte(xml.Header), b...))
	}
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
