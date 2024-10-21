package problem

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
)

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
