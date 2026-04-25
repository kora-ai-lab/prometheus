package browser

import (
	"errors"
	"time"
)

var ErrBrowserUnavailable = errors.New("browser automation is scaffolded but not implemented")

type Cookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Client interface {
	Navigate(url string) error
	GetHTML() (string, error)
	GetText() (string, error)
	Click(selector string) error
	Fill(selector, value string) error
	Submit(selector string) error
	Screenshot() ([]byte, error)
	WaitForSelector(selector string, timeout time.Duration) error
	WaitForNavigation(timeout time.Duration) error
	EvalJS(script string) (any, error)
	GetCookies() ([]*Cookie, error)
	SetCookie(c *Cookie) error
	ScrollDown(pixels int) error
	Close() error
}
