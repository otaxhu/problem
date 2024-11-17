package problem

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	tme_json "github.com/otaxhu/type-mismatch-encoding/encoding/json"
	tme_xml "github.com/otaxhu/type-mismatch-encoding/encoding/xml"
)

// Problem details interface, valid implementors are [MapProblem], [RegisteredProblem] and
// Custom structs that embeds [RegisteredProblem] (is recommended you embed by value, not a pointer,
// in that case you need to allocate appropriate memory for that field)
type Problem interface {
	GetType() string
	GetStatus() int
	GetTitle() string
	GetDetail() string
	GetInstance() string

	// For taking status code from Response object instead of status in problem details body.
	//
	// Also works as a "Must embed" method when creating customs Problems, customs Problems
	// must embed [RegisteredProblem]
	setStatus(status int)

	// For setting "about:blank" to type when the problem detail's type member is not present or
	// has a JSON type other than string.
	setTypeAboutBlank()
}

// NewMap returns a [MapProblem], this implementation is ONLY suitable for JSON
// marshaling/unmarshaling.
//
// If you want XML support, Instead you should use [NewRegistered] or create your own custom
// struct that embeds [RegisteredProblem] and set the fields yourself.
func NewMap(statusCode int, details string) MapProblem {
	return MapProblem{
		"status": statusCode,
		"detail": details,
		"title":  http.StatusText(statusCode),
		"type":   "about:blank",
	}
}

// NewRegistered returns a *[RegisteredProblem], this implementation is suitable for both
// XML and JSON marshaling/unmarshaling.
//
// If you want to add extension members (extra members others than the specified in
// https://datatracker.ietf.org/doc/html/rfc9457#name-members-of-a-problem-detail)
// then you will need to use [NewMap] or create your own custom struct that embeds
// [RegisteredProblem] and set the fields yourself.
func NewRegistered(statusCode int, details string) *RegisteredProblem {
	return &RegisteredProblem{
		Type:   "about:blank",
		Status: statusCode,
		Title:  http.StatusText(statusCode),
		Detail: details,
	}
}

// ParseResponse parses the [http.Response] object into a Problem details.
// The Content-Type header of the response determines the implementation used
//
//  1. If Content-Type is 'application/problem+json' (Problem JSON) then the type of the returned
//     Problem is *[MapProblem] (a pointer)
//
//  2. Else if Content-Type is 'application/problem+xml' (Problem XML) then the type of the
//     returned Problem is *[RegisteredProblem] (a pointer) and found extension members are
//     ignored.
//
// If Content-Type is not one of the first two above, then an error [ErrInvalidContentType] is
// returned, you can check for it using errors.Is(err, ErrInvalidContentType)
func ParseResponse(res *http.Response) (Problem, error) {

	contentType := res.Header.Get("Content-Type")

	var p Problem

	if contentType == MediaTypeProblemJSON {
		// Use a MapProblem
		p = &MapProblem{}
	} else if contentType == MediaTypeProblemXML {
		// Use RegisteredProblem, gonna lost extension members but it's better than failing
		p = &RegisteredProblem{}
	} else {
		return nil, fmt.Errorf("%w: got '%s'", ErrInvalidContentType, contentType)
	}

	err := ParseResponseCustom(res, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// ParseResponseCustom parses the [http.Response] object, unmarshaling it into p argument.
//
// There are some contraints you need to follow in order for this function to work.
//
//   - p must be a pointer.
//   - p must not be nil nor point to nil.
//   - If Content-Type is 'application/problem+xml' (Problem XML), then p must not be a [MapProblem]
//
// If you followed this constraints, then you should get p populated with the Problem details
// values and no errors.
func ParseResponseCustom(res *http.Response, p Problem) error {
	contentType := res.Header.Get("Content-Type")

	if contentType != MediaTypeProblemJSON && contentType != MediaTypeProblemXML {
		return fmt.Errorf("%w: got '%s'", ErrInvalidContentType, contentType)
	}

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return err
	}

	// See issue https://github.com/otaxhu/problem/issues/14
	//
	// This structs checks that "type" is present or has an incorrect type (on JSON)
	// and act accordingly by setting "type" to "about:blank"
	//
	// RFC 9457 Section 3.1 https://www.rfc-editor.org/rfc/rfc9457.html#section-3.1-1
	//
	// "...If a member's value type does not match the specified type, the member MUST
	// be ignored -- i.e., processing will continue as if the member had not been present."
	//
	// RFC 9457 Section 3.1.1 https://www.rfc-editor.org/rfc/rfc9457.html#section-3.1.1-2
	//
	// "When this member ("type") is not present, its value is assumed to be "about:blank"."

	checkTypeJSON := struct {
		Type any `json:"type"`
	}{}

	checkTypeXML := struct {
		XMLName struct{} `xml:"urn:ietf:rfc:7807 problem"`
		Type    *string  `xml:"type"`
	}{}

	br := bytes.NewReader(b)

	switch contentType {
	case MediaTypeProblemJSON:

		dec := tme_json.NewDecoder(br)
		dec.AllowTypeMismatch()

		err := dec.Decode(p)
		if err != nil {
			return err
		}

		if p.GetType() == "" {
			br.Reset(b)

			dec = tme_json.NewDecoder(br)
			dec.AllowTypeMismatch()

			err = dec.Decode(&checkTypeJSON)
			if err != nil {
				return err
			}

			if _, ok := checkTypeJSON.Type.(string); !ok {
				p.setTypeAboutBlank()
			}
		}

	case MediaTypeProblemXML:

		dec := tme_xml.NewDecoder(br)
		dec.AllowTypeMismatch = true

		err := dec.Decode(p)
		if err != nil {
			return err
		}

		if p.GetType() == "" {

			br.Reset(b)

			dec = tme_xml.NewDecoder(br)
			dec.AllowTypeMismatch = true

			err = dec.Decode(&checkTypeXML)
			if err != nil {
				return err
			}

			if checkTypeXML.Type == nil {
				p.setTypeAboutBlank()
			}
		}
	}

	p.setStatus(res.StatusCode)

	return nil
}
