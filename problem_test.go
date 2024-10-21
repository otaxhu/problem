package problem

// TODO: Remade all of the tests using the new API before bumping to V1 stable

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"testing"
// )

// func compareProblems(a, b Problem) bool {
// 	if a.Type != b.Type {
// 		return false
// 	}
// 	if a.Title != b.Title {
// 		return false
// 	}
// 	if a.Detail != b.Detail {
// 		return false
// 	}
// 	if a.Status != b.Status {
// 		return false
// 	}
// 	if a.Instance != b.Instance {
// 		return false
// 	}
// 	// TODO: compare the ExtensionMembers
// 	//
// 	// if !maps.Equal(a.ExtensionMembers, b.ExtensionMembers) {
// 	// 	return false
// 	// }
// 	return true
// }

// func TestNewDefault(t *testing.T) {
// 	testCases := []struct {
// 		Name       string
// 		StatusCode int

// 		ExpectedProblem Problem
// 		ExpectsPanic    bool
// 	}{
// 		{
// 			Name:       "Success_Status200",
// 			StatusCode: 200,

// 			ExpectedProblem: Problem{
// 				Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-200-ok",
// 				Title:  "200 OK",
// 				Status: 200,
// 			},
// 			ExpectsPanic: false,
// 		},
// 		{
// 			Name:       "Fail_Status99",
// 			StatusCode: 99,

// 			ExpectedProblem: Problem{},
// 			ExpectsPanic:    true,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			defer func() {
// 				r := recover()
// 				if r == nil && tc.ExpectsPanic {
// 					t.Errorf("expected NewProblem function to panic, doesn't panicked")
// 				} else if r != nil && !tc.ExpectsPanic {
// 					t.Errorf("expected NewProblem function not to panic, panicked")
// 				}
// 			}()
// 			prob := NewDefault(tc.StatusCode)
// 			if !compareProblems(prob, tc.ExpectedProblem) {
// 				t.Errorf("unexpected Problem value, expected '%v' got '%v'", tc.ExpectedProblem, prob)
// 			}
// 		})
// 	}
// }

// type responseWriter struct {
// 	bytes.Buffer
// 	header     http.Header
// 	statusCode int
// }

// func (r *responseWriter) Header() http.Header {
// 	if r.header == nil {
// 		r.header = http.Header{}
// 	}
// 	return r.header
// }

// func (r *responseWriter) WriteHeader(statusCode int) {
// 	r.statusCode = statusCode
// }

// func TestProblemSend(t *testing.T) {
// 	testCases := []struct {
// 		Name         string
// 		InputProblem Problem

// 		ExpectedProblem     Problem
// 		ExpectedContentType string
// 		ExpectedStatusCode  int
// 	}{
// 		{
// 			Name:         "Success_NewDefault(200)",
// 			InputProblem: NewDefault(200),

// 			ExpectedProblem: Problem{
// 				Type:   "https://www.rfc-editor.org/rfc/rfc9110.html#name-200-ok",
// 				Status: 200,
// 				Title:  "200 OK",
// 			},
// 			ExpectedContentType: "application/problem+json",
// 			ExpectedStatusCode:  200,
// 		},
// 		{
// 			Name: "Fail_ProblemWithBadStatusCode_99",
// 			InputProblem: Problem{
// 				Status: 99,
// 			},

// 			ExpectedProblem: Problem{
// 				Status: 99,
// 			},
// 			ExpectedContentType: "application/problem+json",
// 			ExpectedStatusCode:  99,
// 		},
// 	}
// 	for _, tc := range testCases {
// 		t.Run(tc.Name, func(t *testing.T) {
// 			wr := &responseWriter{}
// 			tc.InputProblem.Send(wr)
// 			if wr.header.Get("Content-Type") != tc.ExpectedContentType {
// 				t.Errorf("unexpecte Content-Type header, expected %q got %q", tc.ExpectedContentType, wr.header.Get("Content-Type"))
// 			}
// 			if wr.statusCode != tc.ExpectedStatusCode {
// 				t.Errorf("unexpected status code, expected '%d' got '%d'", tc.ExpectedStatusCode, wr.statusCode)
// 			}
// 			p := Problem{}
// 			if err := json.NewDecoder(wr).Decode(&p); err != nil {
// 				t.Errorf("unexpected error during JSON unmarshalling, got following error: '%v'", err)
// 			}
// 			if !compareProblems(p, tc.ExpectedProblem) {
// 				t.Errorf("unexpected Problem value, expected '%v' got '%v'", tc.ExpectedProblem, p)
// 			}
// 		})
// 	}
// }
