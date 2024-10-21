package problem

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
	}{
		"JSON: OK": {
			InputBody: responseFactory(http.StatusInternalServerError, problemJsonContentType, `
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
		},
		"XML: OK": {
			InputBody: responseFactory(http.StatusInternalServerError, problemXmlContentType, xml.Header+`
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
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			outProblem, err := ParseResponse(tc.InputBody)
			if err != nil {
				t.Fatal(err)
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
			`),
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

			// TODO: See above, find a way to compact a XML document.
			//
			// recorder = httptest.NewRecorder()
			//
			// ServeXML(tc.InputProblem).ServeHTTP(recorder, req)
			//
			// if recorder.Body.String() != tc.ExpectedXML {
			// 	t.Errorf("expected %s, got %s", tc.ExpectedJSON, recorder.Body.String())
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
	}{
		"JSON: OK": {
			InputBody: responseFactory(http.StatusInternalServerError, problemJsonContentType, `
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
		},
		"XML: OK": {
			InputBody: responseFactory(http.StatusInternalServerError, problemXmlContentType, xml.Header+`
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
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			outProblem := &Embed{}
			err := ParseResponseCustom(tc.InputBody, outProblem)
			if err != nil {
				t.Fatal(err)
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
