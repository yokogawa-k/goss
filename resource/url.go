package resource

import "github.com/aelsabbahy/goss/system"

type URL struct {
	URL           string   `json:"-"`
	Status        int      `json:"status"`
	AllowInsecure bool     `json:"allow-insecure"`
	Body          []string `json:"body"`
}

func (u *URL) ID() string      { return u.URL }
func (u *URL) SetID(id string) { u.URL = id }

func (u *URL) Validate(sys *system.System) []TestResult {
	sysURL := sys.NewURL(u.URL, sys)
	sysURL.SetAllowInsecure(u.AllowInsecure)

	var results []TestResult

	results = append(results, ValidateValue(u, "status", u.Status, sysURL.Status))
	if u.Status == 0 {
		return results
	}
	results = append(results, ValidateContains(u, "Body", u.Body, sysURL.Body))

	return results
}

func NewURL(sysURL system.URL) *URL {
	url := sysURL.URL()
	status, _ := sysURL.Status()
	return &URL{
		URL:           url,
		Status:        status.(int),
		Body:          []string{},
		AllowInsecure: false,
	}
}
