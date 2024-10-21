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
			InputBody: &http.Response{
				Body: io.NopCloser(strings.NewReader(`
					{
						"type": "about:blank",
						"status": 500,
						"title": "Internal Server Error",
						"detail": "test",
						"instance": "/test"
					}
				`)),
				Header: http.Header{
					"Content-Type": []string{problemJsonContentType},
				},
				StatusCode: http.StatusInternalServerError,
			},
			ExpectedProblem: &RegisteredProblem{
				Type:     "about:blank",
				Status:   http.StatusInternalServerError,
				Title:    "Internal Server Error",
				Detail:   "test",
				Instance: "/test",
			},
		},
		"XML: OK": {
			InputBody: &http.Response{
				Body: io.NopCloser(strings.NewReader(xml.Header + `
					<problem xmlns="urn:ietf:rfc:7807">
						<type>about:blank</type>
						<status>500</status>
						<title>Internal Server Error</title>
						<detail>test</detail>
						<instance>/test</instance>
					</problem>
				`)),
				Header: http.Header{
					"Content-Type": []string{problemXmlContentType},
				},
				StatusCode: http.StatusInternalServerError,
			},
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
