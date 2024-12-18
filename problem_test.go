package problem

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"
)

func equalProblems(a Problem, b Problem) bool {
	return a.GetType() == b.GetType() &&
		a.GetTitle() == b.GetTitle() &&
		a.GetStatus() == b.GetStatus() &&
		a.GetInstance() == b.GetInstance() &&
		a.GetDetail() == b.GetDetail()
}

func compactJson(s string) string {
	b := bytes.Buffer{}
	err := json.Compact(&b, []byte(s))
	if err != nil {
		panic(err)
	}
	return b.String()
}

func responseFactory(statusCode int, contentType, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header: http.Header{
			"Content-Type": []string{contentType},
		},
		Body: io.NopCloser(strings.NewReader(body)),
	}
}

func responseFactoryErrorBody(statusCode int, contentType string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header: http.Header{
			"Content-Type": []string{contentType},
		},
		Body: io.NopCloser(iotest.ErrReader(errors.New("TEST"))),
	}
}

func TestRegisteredProblemMarshaling(t *testing.T) {
	testCases := map[string]struct {
		InputProblem *RegisteredProblem
		ExpectedJSON string

		// TODO: find a way to compact a XML document
		//
		// ExpectedMarshalXML string
	}{
		"OK": {
			InputProblem: NewRegistered(http.StatusBadRequest, "test"),
			ExpectedJSON: compactJson(`
				{
					"type": "about:blank",
					"status": 400,
					"title": "Bad Request",
					"detail": "test",
					"instance": ""
				}
			`),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(tc.InputProblem)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b, []byte(tc.ExpectedJSON)) {
				t.Errorf("expected %s, got %s", tc.ExpectedJSON, b)
			}

			// TODO: see above, find a way to compact a XML document
			//
			// b, err = xml.Marshal(tc.InputMarshalProblem)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// if !bytes.Equal(b, []byte(tc.ExpectedMarshalXML)) {
			// 	t.Errorf("marshal problem XML: expected %s, got %s", tc.ExpectedMarshalJSON, b)
			// }
		})
	}
}

func TestParseResponse(t *testing.T) {
	testCases := map[string]struct {
		InputBody       *http.Response
		ExpectedProblem *RegisteredProblem
		ExpectedError   bool
	}{
		"JSON: OK": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"type": "about:blank",
					"status": 500,
					"title": "Internal Server Error",
					"detail": "test",
					"instance": "/test"
				}
			`),
			ExpectedProblem: &RegisteredProblem{
				Type:     "about:blank",
				Status:   http.StatusInternalServerError,
				Title:    "Internal Server Error",
				Detail:   "test",
				Instance: "/test",
			},
			ExpectedError: false,
		},
		"XML: OK": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemXML, xml.Header+`
				<problem xmlns="urn:ietf:rfc:7807">
					<type>about:blank</type>
					<status>500</status>
					<title>Internal Server Error</title>
					<detail>test</detail>
					<instance>/test</instance>
				</problem>
			`),
			ExpectedProblem: &RegisteredProblem{
				Type:     "about:blank",
				Status:   http.StatusInternalServerError,
				Title:    "Internal Server Error",
				Detail:   "test",
				Instance: "/test",
			},
			ExpectedError: false,
		},
		"Bad Content Type": {
			InputBody:       responseFactory(http.StatusInternalServerError, "text/plain", "Test plain"),
			ExpectedProblem: nil,
			ExpectedError:   true,
		},
		"JSON: Bad Syntax": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"status": 500 // Invalid JSON
				}
			`),
			ExpectedProblem: nil,
			ExpectedError:   true,
		},
		"XML: Bad Syntax": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemXML, xml.Header+`
				<problem>
					<status>500</status>
					No Closing Tag
			`),
			ExpectedProblem: nil,
			ExpectedError:   true,
		},
		"XML: Missing Type": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemXML, xml.Header+`
				<problem xmlns="urn:ietf:rfc:7807">
					<status>500</status>
					<title>Internal Server Error</title>
					<detail>test</detail>
					<instance>/test</instance>
				</problem>
			`),
			ExpectedProblem: &RegisteredProblem{
				Type:     "about:blank",
				Status:   500,
				Title:    "Internal Server Error",
				Detail:   "test",
				Instance: "/test",
			},
			ExpectedError: false,
		},
		"XML: Empty Type": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemXML, xml.Header+`
				<problem xmlns="urn:ietf:rfc:7807">
					<type></type>
					<status>500</status>
					<title>Internal Server Error</title>
					<detail>test</detail>
					<instance>/test</instance>
				</problem>
			`),
			ExpectedProblem: &RegisteredProblem{
				Type:     "",
				Status:   500,
				Title:    "Internal Server Error",
				Detail:   "test",
				Instance: "/test",
			},
			ExpectedError: false,
		},
		"JSON: Missing Type": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"status": 500,
					"title": "Internal Server Error",
					"detail": "test",
					"instance": "/test"
				}
			`),
			ExpectedProblem: &RegisteredProblem{
				Type:     "about:blank",
				Status:   500,
				Title:    "Internal Server Error",
				Detail:   "test",
				Instance: "/test",
			},
			ExpectedError: false,
		},
		"JSON: Empty Type": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"type": "",
					"status": 500,
					"title": "Internal Server Error",
					"detail": "test",
					"instance": "/test"
				}
			`),
			ExpectedProblem: &RegisteredProblem{
				Type:     "",
				Status:   500,
				Title:    "Internal Server Error",
				Detail:   "test",
				Instance: "/test",
			},
			ExpectedError: false,
		},
		"JSON: Incorrect JSON Type For Type Member": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"type": 123,
					"status": 500,
					"title": "Internal Server Error",
					"detail": "test",
					"instance": "/test"
				}
			`),
			ExpectedProblem: &RegisteredProblem{
				Type:     "about:blank",
				Status:   500,
				Title:    "Internal Server Error",
				Detail:   "test",
				Instance: "/test",
			},
			ExpectedError: false,
		},
		"Bad Body": {
			InputBody:       responseFactoryErrorBody(http.StatusInternalServerError, MediaTypeProblemJSON),
			ExpectedProblem: nil,
			ExpectedError:   true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			outProblem, err := ParseResponse(tc.InputBody)
			if tc.ExpectedError && err == nil {
				t.Fatalf("expected error to be non-nil, got <nil>")
			} else if !tc.ExpectedError && err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}
			if err != nil {
				return
			}
			if !equalProblems(outProblem, tc.ExpectedProblem) {
				t.Errorf("expected %+v, got %+v", tc.ExpectedProblem, outProblem)
			}
		})
	}
}

func TestServe(t *testing.T) {
	testCases := map[string]struct {
		InputProblem *RegisteredProblem
		ExpectedJSON string

		// TODO: find a way to compact a XML documet
		//
		// ExpectedXML string
	}{
		"OK": {
			InputProblem: NewRegistered(http.StatusBadRequest, "test"),
			ExpectedJSON: compactJson(`
				{
					"type": "about:blank",
					"status": 400,
					"title": "Bad Request",
					"detail": "test",
					"instance": ""
				}
			`) + "\n", // Append newline, because ServeXXX functions uses Encoder, which appends a newline at the end of the stream
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("", "/", nil)

			ServeJSON(tc.InputProblem).ServeHTTP(recorder, req)

			if recorder.Body.String() != tc.ExpectedJSON {
				t.Errorf("expected %s, got %s", tc.ExpectedJSON, recorder.Body.String())
			}

			contentType := recorder.Result().Header.Get("Content-Type")

			if contentType != MediaTypeProblemJSON {
				t.Errorf("expected %s, got %s", MediaTypeProblemJSON, contentType)
			}

			// TODO: See above, find a way to compact a XML document.
			//
			// recorder = httptest.NewRecorder()
			//
			// ServeXML(tc.InputProblem).ServeHTTP(recorder, req)
			//
			// if recorder.Body.String() != tc.ExpectedXML {
			// 	t.Errorf("expected %s, got %s", tc.ExpectedJSON, recorder.Body.String())
			// }
			//
			// contentType = recorder.Result().Header.Get("Content-Type")
			//
			// if contentType != problemJsonContentType {
			// 	t.Errorf("expected %s, got %s", problemJsonContentType, contentType)
			// }
		})
	}
}

type Embed struct {
	Extension1 string `json:"extension1" xml:"extension1"`
	RegisteredProblem
	Extension2 string `json:"extension2" xml:"extension2"`
}

func TestEmbeddedParseResponseCustom(t *testing.T) {
	testCases := map[string]struct {
		InputBody       *http.Response
		ExpectedProblem *Embed
		ExpectedError   bool
	}{
		"JSON: OK": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"type": "about:blank",
					"status": 500,
					"title": "Internal Server Error",
					"detail": "test",
					"instance": "/test",
					"extension1": "e1",
					"extension2": "e2"
				}
			`),
			ExpectedProblem: &Embed{
				Extension1: "e1",
				Extension2: "e2",
				RegisteredProblem: RegisteredProblem{
					Type:     "about:blank",
					Status:   http.StatusInternalServerError,
					Title:    "Internal Server Error",
					Detail:   "test",
					Instance: "/test",
				},
			},
			ExpectedError: false,
		},
		"XML: OK": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemXML, xml.Header+`
				<problem xmlns="urn:ietf:rfc:7807">
					<type>about:blank</type>
					<status>500</status>
					<title>Internal Server Error</title>
					<detail>test</detail>
					<instance>/test</instance>
					<extension1>e1</extension1>
					<extension2>e2</extension2>
				</problem>
			`),
			ExpectedProblem: &Embed{
				Extension1: "e1",
				Extension2: "e2",
				RegisteredProblem: RegisteredProblem{
					Type:     "about:blank",
					Status:   http.StatusInternalServerError,
					Title:    "Internal Server Error",
					Detail:   "test",
					Instance: "/test",
				},
			},
			ExpectedError: false,
		},
		"Bad Content Type": {
			InputBody:       responseFactory(http.StatusInternalServerError, "text/plain", "Test plain"),
			ExpectedProblem: nil,
			ExpectedError:   true,
		},
		"JSON: Bad Syntax": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"status": 500 // Invalid JSON
				}
			`),
			ExpectedProblem: nil,
			ExpectedError:   true,
		},
		"XML: Bad Syntax": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemXML, xml.Header+`
				<problem>
					<status>500</status>
					No Closing Tag
			`),
			ExpectedProblem: nil,
			ExpectedError:   true,
		},
		"XML: Missing Type": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemXML, xml.Header+`
				<problem xmlns="urn:ietf:rfc:7807">
					<status>500</status>
					<title>Internal Server Error</title>
					<detail>test</detail>
					<instance>/test</instance>
					<extension1>e1</extension1>
					<extension2>e2</extension2>
				</problem>
			`),
			ExpectedProblem: &Embed{
				Extension1: "e1",
				Extension2: "e2",
				RegisteredProblem: RegisteredProblem{
					Type:     "about:blank",
					Status:   500,
					Title:    "Internal Server Error",
					Detail:   "test",
					Instance: "/test",
				},
			},
			ExpectedError: false,
		},
		"XML: Empty Type": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemXML, xml.Header+`
				<problem xmlns="urn:ietf:rfc:7807">
					<type></type>
					<status>500</status>
					<title>Internal Server Error</title>
					<detail>test</detail>
					<instance>/test</instance>
					<extension1>e1</extension1>
					<extension2>e2</extension2>
				</problem>
			`),
			ExpectedProblem: &Embed{
				Extension1: "e1",
				Extension2: "e2",
				RegisteredProblem: RegisteredProblem{
					Type:     "",
					Status:   500,
					Title:    "Internal Server Error",
					Detail:   "test",
					Instance: "/test",
				},
			},
			ExpectedError: false,
		},
		"JSON: Missing Type": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"status": 500,
					"title": "Internal Server Error",
					"detail": "test",
					"instance": "/test",
					"extension1": "e1",
					"extension2": "e2"
				}
			`),
			ExpectedProblem: &Embed{
				Extension1: "e1",
				Extension2: "e2",
				RegisteredProblem: RegisteredProblem{
					Type:     "about:blank",
					Status:   500,
					Title:    "Internal Server Error",
					Detail:   "test",
					Instance: "/test",
				},
			},
			ExpectedError: false,
		},
		"JSON: Empty Type": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"type": "",
					"status": 500,
					"title": "Internal Server Error",
					"detail": "test",
					"instance": "/test",
					"extension1": "e1",
					"extension2": "e2"
				}
			`),
			ExpectedProblem: &Embed{
				Extension1: "e1",
				Extension2: "e2",
				RegisteredProblem: RegisteredProblem{
					Type:     "",
					Status:   500,
					Title:    "Internal Server Error",
					Detail:   "test",
					Instance: "/test",
				},
			},
		},
		"JSON: Incorrect JSON Type For Type Member": {
			InputBody: responseFactory(http.StatusInternalServerError, MediaTypeProblemJSON, `
				{
					"type": 123,
					"status": 500,
					"title": "Internal Server Error",
					"detail": "test",
					"instance": "/test",
					"extension1": "e1",
					"extension2": "e2"
				}
			`),
			ExpectedProblem: &Embed{
				Extension1: "e1",
				Extension2: "e2",
				RegisteredProblem: RegisteredProblem{
					Type:     "about:blank",
					Status:   500,
					Title:    "Internal Server Error",
					Detail:   "test",
					Instance: "/test",
				},
			},
			ExpectedError: false,
		},
		"Bad Body": {
			InputBody:       responseFactoryErrorBody(http.StatusInternalServerError, MediaTypeProblemJSON),
			ExpectedProblem: nil,
			ExpectedError:   true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			outProblem := &Embed{}
			err := ParseResponseCustom(tc.InputBody, outProblem)
			if tc.ExpectedError && err == nil {
				t.Fatalf("expected error to be non-nil, got <nil>")
			} else if !tc.ExpectedError && err != nil {
				t.Fatalf("expected error to be nil, got %v", err)
			}
			if err != nil {
				return
			}

			if !equalProblems(outProblem, tc.ExpectedProblem) {
				t.Errorf("expected %+v, got %+v", tc.ExpectedProblem, outProblem)
			}

			if outProblem.Extension1 != tc.ExpectedProblem.Extension1 {
				t.Errorf("expected %s, got %s", tc.ExpectedProblem.Extension1, outProblem.Extension1)
			}

			if outProblem.Extension2 != tc.ExpectedProblem.Extension2 {
				t.Errorf("expected %s, got %s", tc.ExpectedProblem.Extension2, outProblem.Extension2)
			}
		})
	}
}

func TestEmbeddedMarshaling(t *testing.T) {
	testCases := map[string]struct {
		InputProblem Embed
		ExpectedJSON string

		// TODO: find a way to compact a XML document.
		//
		// ExpectedXML string
	}{
		"OK": {
			InputProblem: Embed{
				RegisteredProblem: *NewRegistered(http.StatusBadRequest, "test"),
				Extension1:        "e1",
				Extension2:        "e2",
			},
			ExpectedJSON: compactJson(`
				{
					"extension1": "e1",
					"type": "about:blank",
					"status": 400,
					"title": "Bad Request",
					"detail": "test",
					"instance": "",
					"extension2": "e2"
				}
			`),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(tc.InputProblem)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b, []byte(tc.ExpectedJSON)) {
				t.Errorf("expected %s, got %s", tc.ExpectedJSON, b)
			}

			// TODO: see above, find a way to compact a XML document.
			//
			// b, err = xml.Marshal(tc.InputProblem)
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// if !bytes.Equal(b, []byte(tc.ExpectedXML)) {
			// 	t.Errorf("expected %s, got %s", tc.ExpectedXML, b)
			// }
		})
	}
}

// IMPORTANT:
//
// Since maps are unordered data structures by nature, this test will only check if
// InputProblem produces valid JSON when marshaled, it will not do a very deep comparison.
func TestMapProblemMarshaling(t *testing.T) {

	testCases := map[string]struct {
		InputProblem MapProblem
	}{
		"OK": {
			InputProblem: NewMap(http.StatusInternalServerError, "Test"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			b, err := json.Marshal(tc.InputProblem)
			if err != nil {
				t.Fatal(err)
			}
			if !json.Valid(b) {
				t.Fatalf("expected JSON to be valid")
			}
		})
	}
}
