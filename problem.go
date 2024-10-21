package problem

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
)

// Problem details interface, valid implementors are [MapProblem], [RegisteredProblem] and
// Custom structs that embeds [RegisteredProblem] (is recommended you embed by value, not a pointer,
// in that case you need to allocate appropiate memory for that field)
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

	if contentType == problemJsonContentType {
		// Use a MapProblem
		p = &MapProblem{}
	} else if contentType == problemXmlContentType {
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

	if contentType != problemJsonContentType && contentType != problemXmlContentType {
		return fmt.Errorf("%w: got '%s'", ErrInvalidContentType, contentType)
	}

	defer res.Body.Close()

	if contentType == problemJsonContentType {

		err := json.NewDecoder(res.Body).Decode(p)
		if err != nil {
			return err
		}

		p.setStatus(res.StatusCode)

	} else if contentType == problemXmlContentType {

		err := xml.NewDecoder(res.Body).Decode(p)
		if err != nil {
			return err
		}

		p.setStatus(res.StatusCode)

	}

	return nil
}
