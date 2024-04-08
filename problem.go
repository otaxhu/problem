package problem

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Problem details struct, you can get an instance with default values through "NewDefault"
// constructor, or you can build your own Problem by instatiate it yourself and set the wanted
// fields.
//
// Then you can use (Problem).Send() method to send the Problem details through a
// http.ResponseWriter
type Problem struct {
	Type             string
	Status           int
	Title            string
	Detail           string
	Instance         string
	ExtensionMembers map[string]any

	OmitEmpty bool
}

// Builds a new Problem with default values, if "statusCode" is an invalid HTTP status code, then
// function panics.
//
// "omitEmpty" param if present and true, omits fields with zero values during JSON marshaling.
// Can be setted again to another value through "OmitEmpty" struct field
func NewDefault(statusCode int, omitEmpty ...bool) Problem {
	if val, ok := problems[statusCode]; ok {
		if len(omitEmpty) > 0 {
			val.OmitEmpty = omitEmpty[0]
		}
		return val
	}
	panic(fmt.Sprintf("problem.go: Invalid '%03d' status code at NewDefault() function", statusCode))
}

// Adds "key" and "value" to the "ExtensionMembers" map, if "key" already exists in map, then "key"
// is overwrited.
//
// If "p.ExtensionMembers" is nil, a new map will be created and the step above will be performed
//
// If "key" is one of "standardMembers" keys defined in "common.go", then this function will be no-op
func (p *Problem) AddExtensionMember(key string, value any) {
	if _, isStandard := standardMembers[key]; isStandard {
		return
	}
	if p.ExtensionMembers == nil {
		p.ExtensionMembers = map[string]any{}
	}
	p.ExtensionMembers[key] = value
}

// Writes "Problem" struct to "w" with "application/problem+json" as Content-Type header and status
// code set to "p.Status" field value.
//
// Further calls to "w" methods should not be performed after call to this method
func (p Problem) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", problemJsonContentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p)
}
