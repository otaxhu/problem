package problem

import (
	"encoding/json"
	"errors"
)

const problemJsonContentType = "application/problem+json"

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
