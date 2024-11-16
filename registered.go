package problem

// Media type for JSON Problem Details
//
// https://datatracker.ietf.org/doc/html/rfc9457#name-iana-considerations
const MediaTypeProblemJSON = "application/problem+json"

// Media type for XML Problem Details
//
// https://datatracker.ietf.org/doc/html/rfc9457#name-iana-considerations
const MediaTypeProblemXML = "application/problem+xml"

// Problem details registered members, as those specified in
// https://datatracker.ietf.org/doc/html/rfc9457#name-members-of-a-problem-detail
//
// This struct can be embedded in a custom struct that adds "extension members" (members that are
// not registered in the above link) like this:
//
//	type Custom struct {
//	    problem.RegisteredProblem
//	    ExtensionMember string `json:"extension_member" xml:"extension_member"`
//	}
type RegisteredProblem struct {
	// NOTE: used in XML marshaling and unmarshaling
	XMLName struct{} `json:"-" xml:"urn:ietf:rfc:7807 problem"`

	Type     string `json:"type" xml:"type"`
	Status   int    `json:"status" xml:"status"`
	Title    string `json:"title" xml:"title"`
	Detail   string `json:"detail" xml:"detail"`
	Instance string `json:"instance" xml:"instance"`
}

func (r RegisteredProblem) GetType() string {
	return r.Type
}

func (r RegisteredProblem) GetStatus() int {
	return r.Status
}

func (r RegisteredProblem) GetTitle() string {
	return r.Title
}

func (r RegisteredProblem) GetDetail() string {
	return r.Detail
}

func (r RegisteredProblem) GetInstance() string {
	return r.Instance
}

func (r *RegisteredProblem) setStatus(status int) {
	r.Status = status
}

// Problem details map, this implementation is ONLY suitable for JSON marshaling/unmarshaling,
// it does not support XML.
//
// This implementation is useful when you want to unmarshal an "extension member", which you
// would otherwise have to create a custom struct for. You can get extension members with map
// access notation:
//
//	mapProblem["extension_member"]
type MapProblem map[string]any

func (m MapProblem) GetType() string {
	v, _ := m["type"].(string)
	return v
}

func (m MapProblem) GetStatus() int {
	v, _ := m["status"].(int)
	return v
}

func (m MapProblem) GetTitle() string {
	v, _ := m["title"].(string)
	return v
}

func (m MapProblem) GetDetail() string {
	v, _ := m["detail"].(string)
	return v
}

func (m MapProblem) GetInstance() string {
	v, _ := m["instance"].(string)
	return v
}

func (m MapProblem) setStatus(status int) {
	m["status"] = status
}
