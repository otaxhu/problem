package problem

const problemJsonContentType = "application/problem+json"

const problemXmlContentType = "application/problem+xml"

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
