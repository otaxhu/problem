package problem

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Problem struct {
	Type             string
	Status           int
	Title            string
	Detail           string
	Instance         string
	ExtensionMembers map[string]any
}

func NewProblem(statusCode int) Problem {
	if val, ok := problems[statusCode]; ok {
		val.Status = statusCode
		val.ExtensionMembers = map[string]any{}
		return val
	}
	panic(fmt.Sprintf("problem.go: Invalid '%03d' status code at NewProblem() function", statusCode))
}

func (p *Problem) AddExtensionMember(key string, value any) {
	if _, isStandard := standardMembers[key]; isStandard {
		return
	}
	if p.ExtensionMembers == nil {
		p.ExtensionMembers = map[string]any{}
	}
	p.ExtensionMembers[key] = value
}

func (p Problem) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, 5+len(p.ExtensionMembers))
	for k, v := range p.ExtensionMembers {
		if _, isStandard := standardMembers[k]; !isStandard {
			m[k] = v
		}
	}
	if p.Type != "" {
		m["type"] = p.Type
	}
	if ok := http.StatusText(p.Status); ok != "" {
		m["status"] = p.Status
	} else {
		panic(fmt.Sprintf("problem.go: Invalid '%03d' status code at (Problem).MarshalJSON() function", p.Status))
	}
	if p.Title != "" {
		m["title"] = p.Title
	}
	if p.Detail != "" {
		m["detail"] = p.Detail
	}
	if p.Instance != "" {
		m["instance"] = p.Instance
	}
	return json.Marshal(m)
}

func (p Problem) Send(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/problem+json")
	if ok := http.StatusText(p.Status); ok == "" {
		panic(fmt.Sprintf("problem.go: Invalid '%03d' status code at (Problem).Send() function", p.Status))
	}
	w.WriteHeader(p.Status)
	json.NewEncoder(w).Encode(p)
}
