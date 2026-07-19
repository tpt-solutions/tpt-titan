// backend/services/outbound_policy_test.go
// Tests for the admin outbound domain allowlist policy layer.

package services

import (
	"encoding/json"
	"testing"

	"tpt-titan/backend/config"
	"tpt-titan/backend/models"
)

func setOutboundDomains(t *testing.T, domains []string) {
	t.Helper()
	db := config.GetDatabase()
	if db == nil {
		t.Skip("no database configured for test")
	}
	encoded, err := json.Marshal(domains)
	if err != nil {
		t.Fatalf("marshal domains: %v", err)
	}
	var setting models.SystemSetting
	if err := db.Where("key = ?", "outbound_domains").First(&setting).Error; err != nil {
		setting = models.SystemSetting{Key: "outbound_domains", Value: string(encoded)}
		db.Create(&setting)
	} else {
		setting.Value = string(encoded)
		db.Save(&setting)
	}
	t.Cleanup(func() {
		db.Where("key = ?", "outbound_domains").Delete(&models.SystemSetting{})
	})
}

func TestOutboundDomainAllowed_AllowAllWhenUnset(t *testing.T) {
	// Ensure no setting exists.
	db := config.GetDatabase()
	if db == nil {
		t.Skip("no database configured for test")
	}
	db.Where("key = ?", "outbound_domains").Delete(&models.SystemSetting{})

	allowed, err := outboundDomainAllowed("https://api.example.com/x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Fatal("expected allow-all when no allowlist is configured")
	}
}

func TestOutboundDomainAllowed_ExactAndSuffixMatch(t *testing.T) {
	setOutboundDomains(t, []string{"api.example.com", "*.other.com"})

	cases := []struct {
		url    string
		wantOk bool
	}{
		{"https://api.example.com/x", true},
		{"https://sub.api.example.com/x", true}, // suffix match
		{"https://a.other.com/x", true},          // wildcard
		{"https://other.com/x", true},            // wildcard base
		{"https://evil.com/x", false},
	}
	for _, c := range cases {
		got, err := outboundDomainAllowed(c.url)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", c.url, err)
		}
		if got != c.wantOk {
			t.Fatalf("%s: expected allowed=%v, got %v", c.url, c.wantOk, got)
		}
	}
}

func TestOutboundDomainAllowed_BlocksWhenNotListed(t *testing.T) {
	setOutboundDomains(t, []string{"allowed.com"})

	allowed, err := outboundDomainAllowed("https://blocked.com/x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("expected blocked.com to be denied by the allowlist")
	}
}

// jsonMarshalDomains is a tiny local helper to avoid importing encoding/json
// at call sites above (kept here to keep the test file self-contained).
func jsonMarshalDomains(domains []string) ([]byte, error) {
	return json.Marshal(domains)
}
