package problem

import (
	"encoding/json"
	"errors"
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
func (p *Problem) AddExtensionMember(key string, value any) {
	if _, isStandard := standardMembers[key]; isStandard {
		return
	}
	if p.ExtensionMembers == nil {
		p.ExtensionMembers = map[string]any{}
	}
	p.ExtensionMembers[key] = value
}

// Implements json.Marshaler interface
//
// NOTE:
// Should not be used in application code, instead you should use (json.Encoder).Encode() method or
// json.Marshal() function to encode the Problem struct
func (p Problem) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, 5+len(p.ExtensionMembers))
	for k, v := range p.ExtensionMembers {
		if _, isStandard := standardMembers[k]; !isStandard {
			m[k] = v
		}
	}
	if p.Type != "" || !p.OmitEmpty {
		m["type"] = p.Type
	}
	if p.Status != 0 || !p.OmitEmpty {
		m["status"] = p.Status
	}
	if p.Title != "" || !p.OmitEmpty {
		m["title"] = p.Title
	}
	if p.Detail != "" || !p.OmitEmpty {
		m["detail"] = p.Detail
	}
	if p.Instance != "" || !p.OmitEmpty {
		m["instance"] = p.Instance
	}
	return json.Marshal(m)
}

// Implements json.Unmarshaler interface
//
// NOTE:
// Just implemented for convinience in tests, unmarshals the buffer to the Problem struct
func (p *Problem) UnmarshalJSON(b []byte) error {
	var x any
	if err := json.Unmarshal(b, &x); err != nil {
		return err
	}
	m, ok := x.(map[string]any)
	if !ok {
		return errors.New("problem.go: Cannot unmarshall the JSON to the Problem struct, x is not a map[string]any")
	}
	for k, v := range m {
		switch k {
		case "type":
			val, ok := v.(string)
			if !ok {
				return errors.New("problem.go: type of 'type' JSON member is not string")
			}
			p.Type = val
		case "status":
			val, ok := v.(float64)
			if !ok {
				return errors.New("problem.go: type of 'status' JSON member is not float64")
			}
			p.Status = int(val)
		case "title":
			val, ok := v.(string)
			if !ok {
				return errors.New("problem.go: type of 'title' JSON member is not string")
			}
			p.Title = val
		case "detail":
			val, ok := v.(string)
			if !ok {
				return errors.New("problem.go: type of 'detail' JSON member is not string")
			}
			p.Detail = val
		case "instance":
			val, ok := v.(string)
			if !ok {
				return errors.New("problem.go: type of 'instance' JSON member is not string")
			}
			p.Instance = val
		default:
			p.ExtensionMembers[k] = v
		}
	}
	return nil
}

// Writes "Problem" struct to "w" with "application/problem+json" as Content-Type header and status
// code set to "p.Status" field value.
//
// Further calls to "w" methods should not be performed after call to this method
func (p Problem) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p)
}
