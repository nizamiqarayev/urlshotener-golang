package urlutil

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

func NormalizeHTTPURL(rawURL string) (string, bool) {
	rawURL = strings.TrimSpace(rawURL)
	rawURL = strings.Trim(rawURL, "\"'`<>")

	if rawURL == "" || containsSpaceOrControl(rawURL) {
		return "", false
	}

	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", false
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", false
	}
	if parsedURL.User != nil || parsedURL.Host == "" || parsedURL.Hostname() == "" {
		return "", false
	}

	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)
	parsedURL.Host = strings.ToLower(parsedURL.Host)

	if !isValidHost(parsedURL.Hostname()) || !isValidPort(parsedURL.Port()) {
		return "", false
	}

	return parsedURL.String(), true
}

func containsSpaceOrControl(value string) bool {
	for _, char := range value {
		if unicode.IsSpace(char) || unicode.IsControl(char) {
			return true
		}
	}

	return false
}

func isValidHost(host string) bool {
	if host == "localhost" || net.ParseIP(host) != nil {
		return true
	}
	if len(host) > 253 || !strings.Contains(host, ".") {
		return false
	}

	labels := strings.Split(host, ".")
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return false
		}
		for _, char := range label {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				return false
			}
		}
	}

	return true
}

func isValidPort(port string) bool {
	if port == "" {
		return true
	}

	portNumber, err := strconv.Atoi(port)
	return err == nil && portNumber > 0 && portNumber <= 65535
}
