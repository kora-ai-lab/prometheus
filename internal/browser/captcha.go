package browser

import (
	"strings"
)

// captchaSelectors contains strings that identify CAPTCHA elements
var captchaSelectors = []string{
	"captcha",
	"g-recaptcha",
	"h-captcha",
	"recaptcha",
	"data-sitekey",
	"cf-challenge",
	"challenge",
}

// DetectCAPTCHA checks if HTML contains any CAPTCHA elements.
// It performs case-insensitive detection by converting HTML to lowercase
// and searching for known CAPTCHA selectors.
func DetectCAPTCHA(html string) bool {
	// Convert HTML to lowercase for case-insensitive detection
	lowerHTML := strings.ToLower(html)

	// Check if any CAPTCHA selector is present in the HTML
	for _, selector := range captchaSelectors {
		if strings.Contains(lowerHTML, selector) {
			return true
		}
	}

	return false
}
