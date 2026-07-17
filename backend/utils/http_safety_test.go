// backend/utils/http_safety_test.go
// Run with: cd backend && go test ./utils/... -run TestValidateOutboundURL -v

package utils

import "testing"

func TestValidateOutboundURL_RejectsLoopback(t *testing.T) {
	if err := ValidateOutboundURL("http://127.0.0.1/admin"); err == nil {
		t.Error("expected loopback address to be rejected")
	}
}

func TestValidateOutboundURL_RejectsLocalhostHostname(t *testing.T) {
	if err := ValidateOutboundURL("http://localhost:8080/internal"); err == nil {
		t.Error("expected localhost to be rejected")
	}
}

func TestValidateOutboundURL_RejectsLinkLocalMetadataAddress(t *testing.T) {
	// 169.254.169.254 is the well-known cloud metadata endpoint; it's covered
	// by IsLinkLocalUnicast with no special-casing needed.
	if err := ValidateOutboundURL("http://169.254.169.254/latest/meta-data/"); err == nil {
		t.Error("expected the cloud metadata address to be rejected")
	}
}

func TestValidateOutboundURL_RejectsPrivateRanges(t *testing.T) {
	for _, u := range []string{
		"http://10.0.0.5/",
		"http://192.168.1.1/",
		"http://172.16.0.1/",
	} {
		if err := ValidateOutboundURL(u); err == nil {
			t.Errorf("expected private address %q to be rejected", u)
		}
	}
}

func TestValidateOutboundURL_RejectsUnsupportedScheme(t *testing.T) {
	if err := ValidateOutboundURL("ftp://example.com/file"); err == nil {
		t.Error("expected non-http(s) scheme to be rejected")
	}
}

func TestValidateOutboundURL_RejectsInvalidURL(t *testing.T) {
	if err := ValidateOutboundURL("not a url at all :://"); err == nil {
		t.Error("expected an unparseable URL to be rejected")
	}
}

func TestValidateOutboundURL_AcceptsPublicLiteralIP(t *testing.T) {
	// Using a well-known public literal IP avoids depending on real DNS
	// resolution succeeding in CI.
	if err := ValidateOutboundURL("https://8.8.8.8/"); err != nil {
		t.Errorf("expected a public IP to be accepted, got error: %v", err)
	}
}
