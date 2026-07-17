package utils

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

// maxRedirects bounds how many redirects a SafeHTTPClient will follow before
// giving up — this also bounds how many times CheckRedirect re-validates the
// destination URL below.
const maxRedirects = 5

// ValidateOutboundURL rejects any URL that would cause the server to make a
// request to itself or to internal/private network space. Workflow nodes let
// a user configure an arbitrary destination URL that this server will call —
// without this check that's a direct SSRF path to internal services,
// localhost admin endpoints, or a cloud metadata endpoint (169.254.169.254,
// which IsLinkLocalUnicast already covers with no special case needed).
func ValidateOutboundURL(rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme %q (only http/https are allowed)", parsed.Scheme)
	}

	host := parsed.Hostname()
	if host == "" {
		return fmt.Errorf("URL has no host")
	}

	var ips []net.IP
	if ip := net.ParseIP(host); ip != nil {
		ips = []net.IP{ip}
	} else {
		resolved, err := net.LookupIP(host)
		if err != nil {
			return fmt.Errorf("failed to resolve host %q: %w", host, err)
		}
		ips = resolved
	}

	for _, ip := range ips {
		if isDisallowedIP(ip) {
			return fmt.Errorf("URL resolves to a disallowed address (%s) — internal/private network destinations are not permitted", ip)
		}
	}

	return nil
}

func isDisallowedIP(ip net.IP) bool {
	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified()
}

// SafeHTTPClient returns an http.Client with the given timeout that
// re-validates the destination of every redirect via ValidateOutboundURL and
// gives up after maxRedirects — so a URL that passes validation can't be used
// to jump to a private address via a redirect chain.
func SafeHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("stopped after %d redirects", maxRedirects)
			}
			if err := ValidateOutboundURL(req.URL.String()); err != nil {
				return fmt.Errorf("redirect target rejected: %w", err)
			}
			return nil
		},
	}
}
