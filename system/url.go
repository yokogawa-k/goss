package system

import (
	"crypto/tls"
	"io"
	"net/http"
)

type URL interface {
	URL() string
	Status() (interface{}, error)
	Body() (io.Reader, error)
	Exists() (interface{}, error)
	SetAllowInsecure(bool)
}

type DefURL struct {
	url           string
	allowInsecure bool
	resp          *http.Response
	loaded        bool
	err           error
}

func NewDefURL(url string, system *System) URL {
	return &DefURL{url: url}
}

func (u *DefURL) setup() error {
	if u.loaded {
		return u.err
	}
	u.loaded = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: u.allowInsecure},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(u.url)

	if err != nil {
		u.err = err
	}

	u.resp = resp

	return u.err
}

func (u *DefURL) Exists() (interface{}, error) { return u.Status() }

func (u *DefURL) SetAllowInsecure(t bool) {
	u.allowInsecure = t
}

func (u *DefURL) ID() string {
	return u.url
}
func (u *DefURL) URL() string {
	return u.url
}

func (u *DefURL) Status() (interface{}, error) {
	err := u.setup()

	if err != nil {
		return 0, nil
	}
	return u.resp.StatusCode, nil
}

func (u *DefURL) Body() (io.Reader, error) {
	err := u.setup()

	if err != nil {
		return nil, err
	}
	return u.resp.Body, nil
}
