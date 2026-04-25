package security

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type DASTConfig struct {
	TargetURL     string
	Timeout       time.Duration
	FollowRedirects bool
}

type DASTResult struct {
	URL               string
	StatusCode       int
	Headers           map[string]string
	Server            string
	SecurityHeaders   map[string]string
	TLSVersion         string
	CipherSuite       string
	OpenEndpoints      []string
	Findings          []string
}

var commonEndpoints = []string{
	"/admin", "/administrator", "/api", "/api/v1", "/api/v2",
	"/backup", "/backups", "/config", "/configuration", "/debug",
	"/login", "/signin", "/register", "/admin.php", "/login.php",
	"/console", "/status", "/health", "/info", "/metrics",
}

func NewDAST(cfg DASTConfig) *DAST {
	return &DAST{cfg: cfg}
}

type DAST struct {
	cfg DASTConfig
	client *http.Client
}

func (d *DAST) Scan() (*DASTResult, error) {
	timeout := 10 * time.Second
	if d.cfg.Timeout > 0 {
		timeout = d.cfg.Timeout
	}

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(d.cfg.TargetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &DASTResult{
		URL:             d.cfg.TargetURL,
		StatusCode:      resp.StatusCode,
		Headers:          make(map[string]string),
		SecurityHeaders: make(map[string]string),
	}

	for k, v := range resp.Header {
		if len(v) > 0 {
			result.Headers[k] = v[0]
		}
	}

	result.Server = resp.Header.Get("Server")

	result.SecurityHeaders["Strict-Transport-Security"] = resp.Header.Get("Strict-Transport-Security")
	result.SecurityHeaders["X-Frame-Options"] = resp.Header.Get("X-Frame-Options")
	result.SecurityHeaders["X-Content-Type-Options"] = resp.Header.Get("X-Content-Type-Options")
	result.SecurityHeaders["Content-Security-Policy"] = resp.Header.Get("Content-Security-Policy")

	return result, nil
}

func (d *DAST) CheckTLS() (*DASTResult, error) {
	host := d.cfg.TargetURL

	conn, err := tls.Dial("tcp", host, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	state := conn.ConnectionState()
	cipherSuite := tls.CipherSuiteName(state.CipherSuite)
	result := &DASTResult{
		URL:           d.cfg.TargetURL,
		TLSVersion:    tlsVersion(state.Version),
		CipherSuite:   cipherSuite,
	}

	return result, nil
}

func (d *DAST) FuzzEndpoints() []string {
	var found []string
	client := &http.Client{Timeout: 5 * time.Second}

	for _, endpoint := range commonEndpoints {
		url := d.cfg.TargetURL + endpoint
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				found = append(found, endpoint)
			}
		}
	}

	return found
}

func tlsVersion(v uint16) string {
	switch v {
	case tls.VersionTLS10:
		return "1.0"
	case tls.VersionTLS11:
		return "1.1"
	case tls.VersionTLS12:
		return "1.2"
	case tls.VersionTLS13:
		return "1.3"
	default:
		return "unknown"
	}
}

func init() {
	_ = net.Dial
}