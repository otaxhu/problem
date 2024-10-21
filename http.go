package problem

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
)

// ServeXML returns a Handler that serves the p argument in XML format.
//
// p MUST not be a [MapProblem], since it cannot be marshaled to XML.
//
// Headers Content-Type is set to 'application/problem+xml' and X-Content-Type-Options is set to 'nosniff';
// and finally writes the status code from p.GetStatus().
func ServeXML(p Problem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		b, _ := xml.Marshal(p)
		h.Set("Content-Type", problemXmlContentType)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Content-Length", strconv.Itoa(len(xml.Header)+len(b)))
		w.WriteHeader(p.GetStatus())
		w.Write(append([]byte(xml.Header), b...))
	})
}

// ServeJSON returns a Handler that serves the p argument in JSON format
//
// Headers Content-Type is set to 'application/problem+json' and X-Content-Type-Options is set to
// 'nosniff'; and finally writes the status code from p.GetStatus().
func ServeJSON(p Problem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		b, _ := json.Marshal(p)
		h.Set("Content-Type", problemJsonContentType)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Content-Length", strconv.Itoa(len(b)))
		w.WriteHeader(p.GetStatus())
		w.Write(b)
	})
}
