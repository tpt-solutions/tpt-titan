package services

import (
	"encoding/json"
	"net/url"
	"strings"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
)

// deliveryLogBodyLimit caps request/response bodies stored in the delivery log
// so a large payload can't blow up the table.
const deliveryLogBodyLimit = 4096

// hostOf returns the lower-cased host (without port) of a URL, or "" if the
// URL can't be parsed. Used for outbound domain allowlist matching and for
// recording the delivery log.
func hostOf(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.ToLower(parsed.Hostname())
}

// outboundDomainAllowed implements the admin-configured outbound domain
// allowlist. It reads the "outbound_domains" SystemSetting (a JSON array of
// allowed host suffixes, e.g. ["api.example.com","*.example.com"]). An empty
// or unset allowlist permits all destinations — this is purely a policy layer
// on top of the SSRF guard (utils.ValidateOutboundURL), which still blocks
// internal/private addresses regardless.
func outboundDomainAllowed(rawURL string) (bool, error) {
	host := hostOf(rawURL)
	if host == "" {
		return false, nil
	}

	db := config.GetDatabase()
	if db == nil {
		return true, nil
	}
	var setting models.SystemSetting
	if err := db.Where("key = ?", "outbound_domains").First(&setting).Error; err != nil {
		// Not configured → allow all public destinations.
		return true, nil
	}

	var allowed []string
	if err := json.Unmarshal([]byte(setting.Value), &allowed); err != nil || len(allowed) == 0 {
		return true, nil
	}

	for _, entry := range allowed {
		pattern := strings.ToLower(strings.TrimSpace(entry))
		if pattern == "" {
			continue
		}
		if pattern == host {
			return true, nil
		}
		// Suffix match: "example.com" allows "api.example.com".
		if strings.HasSuffix(host, "."+pattern) {
			return true, nil
		}
		// Wildcard: "*.example.com" allows any subdomain.
		if strings.HasPrefix(pattern, "*.") {
			base := pattern[2:]
			if host == base || strings.HasSuffix(host, "."+base) {
				return true, nil
			}
		}
	}
	return false, nil
}
