package browser

import (
	"context"
	"fmt"
	"time"
)

type PlaywrightClient struct {
	browser_  any
	page_     any
	capEngine interface{}
}

func NewPlaywrightClient(home string, capEngine interface{}) (*PlaywrightClient, error) {
	if capEngine != nil {
		if ce, ok := capEngine.(interface{ Ensure(context.Context, string) error }); ok {
			ce.Ensure(context.Background(), "playwright")
		}
	}

	return &PlaywrightClient{
		capEngine: capEngine,
	}, nil
}

func (c *PlaywrightClient) Navigate(url string) error {
	return fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) GetHTML() (string, error) {
	return "", fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) GetText() (string, error) {
	return "", fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) Click(selector string) error {
	_ = selector
	return fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) Fill(selector, value string) error {
	_, _ = selector, value
	return fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) Submit(selector string) error {
	_ = selector
	return fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) Screenshot() ([]byte, error) {
	return nil, fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) WaitForSelector(selector string, timeout time.Duration) error {
	_, _ = selector, timeout
	return fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) WaitForNavigation(timeout time.Duration) error {
	_ = timeout
	return fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) EvalJS(script string) (any, error) {
	_ = script
	return nil, fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) GetCookies() ([]*Cookie, error) {
	return nil, fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) SetCookie(cookie *Cookie) error {
	_ = cookie
	return fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) ScrollDown(pixels int) error {
	_ = pixels
	return fmt.Errorf("playwright not available in this build")
}

func (c *PlaywrightClient) Close() error {
	return nil
}
