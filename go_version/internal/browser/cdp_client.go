package browser

import "time"

type StubClient struct{}

func NewStubClient() *StubClient { return &StubClient{} }

func (c *StubClient) Navigate(url string) error   { _ = url; return ErrBrowserUnavailable }
func (c *StubClient) GetHTML() (string, error)    { return "", ErrBrowserUnavailable }
func (c *StubClient) GetText() (string, error)    { return "", ErrBrowserUnavailable }
func (c *StubClient) Click(selector string) error { _ = selector; return ErrBrowserUnavailable }
func (c *StubClient) Fill(selector, value string) error {
	_, _ = selector, value
	return ErrBrowserUnavailable
}
func (c *StubClient) Submit(selector string) error { _ = selector; return ErrBrowserUnavailable }
func (c *StubClient) Screenshot() ([]byte, error)  { return nil, ErrBrowserUnavailable }
func (c *StubClient) WaitForSelector(selector string, timeout time.Duration) error {
	_, _ = selector, timeout
	return ErrBrowserUnavailable
}
func (c *StubClient) WaitForNavigation(timeout time.Duration) error {
	_ = timeout
	return ErrBrowserUnavailable
}
func (c *StubClient) EvalJS(script string) (any, error) {
	_ = script
	return nil, ErrBrowserUnavailable
}
func (c *StubClient) GetCookies() ([]*Cookie, error) { return nil, ErrBrowserUnavailable }
func (c *StubClient) SetCookie(cookie *Cookie) error { _ = cookie; return ErrBrowserUnavailable }
func (c *StubClient) ScrollDown(pixels int) error    { _ = pixels; return ErrBrowserUnavailable }
func (c *StubClient) Close() error                   { return nil }
