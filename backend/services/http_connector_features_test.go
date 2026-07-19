// backend/services/http_connector_features_test.go
// Tests for the http.request connector reliability/signing/extraction features.

package services

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The SSRF guard (utils.ValidateOutboundURL) correctly rejects the loopback
// address an httptest server binds to, so these tests exercise
// executeHTTPRequest directly (the request machinery after validation).

func TestHTTPRequestConnector_ExtractField_Nested(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"data":{"order":{"id":"ORD-42"}}}`)
	}))
	defer srv.Close()

	c := &HTTPRequestConnector{}
	res, err := c.executeHTTPRequest(map[string]interface{}{
		"extract_field": "data.order.id",
	}, srv.URL, "", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["extracted_field"] != "ORD-42" {
		t.Fatalf("expected extracted_field ORD-42, got %v", res["extracted_field"])
	}
	if res["extracted_field_name"] != "data.order.id" {
		t.Fatalf("expected extracted_field_name, got %v", res["extracted_field_name"])
	}
}

func TestHTTPRequestConnector_ExtractField_ArrayIndex(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"items":[{"id":1},{"id":2}]}`)
	}))
	defer srv.Close()

	c := &HTTPRequestConnector{}
	res, err := c.executeHTTPRequest(map[string]interface{}{
		"extract_field": "items.1.id",
	}, srv.URL, "", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["extracted_field"] != float64(2) {
		t.Fatalf("expected extracted_field 2, got %v", res["extracted_field"])
	}
}

func TestHTTPRequestConnector_ExtractField_Missing(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"ok":true}`)
	}))
	defer srv.Close()

	c := &HTTPRequestConnector{}
	res, err := c.executeHTTPRequest(map[string]interface{}{
		"extract_field": "nope.deep",
	}, srv.URL, "", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res["extracted_field_error"]; !ok {
		t.Fatalf("expected extracted_field_error when field is absent, got %v", res)
	}
}

func TestHTTPRequestConnector_SigningHeader(t *testing.T) {
	var gotSig, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get("X-Titan-Signature")
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
	}))
	defer srv.Close()

	secret := "shared-secret"
	c := &HTTPRequestConnector{}
	_, err := c.executeHTTPRequest(map[string]interface{}{
		"method":         "POST",
		"body":           `{"hello":"world"}`,
		"signing_secret": secret,
	}, srv.URL, "POST", []byte(`{"hello":"world"}`), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(gotSig, "sha256=") {
		t.Fatalf("expected X-Titan-Signature header, got %q", gotSig)
	}

	expected := signRequestPayload("POST", srv.URL, []byte(gotBody), secret)
	if gotSig != "sha256="+expected {
		t.Fatalf("signature does not match recomputed HMAC")
	}
}

func TestHTTPRequestConnector_RetrySucceedsOnSecondAttempt(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := &HTTPRequestConnector{}
	res, err := c.executeHTTPRequest(map[string]interface{}{
		"retry_attempts": float64(3),
	}, srv.URL, "", nil, nil)
	if err != nil {
		t.Fatalf("expected eventual success after retry, got error: %v", err)
	}
	if attempts != 2 {
		t.Fatalf("expected exactly 2 attempts, got %d", attempts)
	}
	if res["success"] != true {
		t.Fatalf("expected success from retried response, got %v", res)
	}
}

func TestHTTPRequestConnector_RetryExhausts(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	c := &HTTPRequestConnector{}
	_, err := c.executeHTTPRequest(map[string]interface{}{
		"retry_attempts": float64(2),
	}, srv.URL, "", nil, nil)
	if err == nil {
		t.Fatal("expected an error after retries exhausted")
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestExtractJSONField_Unit(t *testing.T) {
	root := map[string]interface{}{}
	_ = json.Unmarshal([]byte(`{"a":{"b":[10,20]}}`), &root)

	if v, ok := extractJSONField(root, "a.b.1"); !ok || v != float64(20) {
		t.Fatalf("expected a.b.1 == 20, got %v ok=%v", v, ok)
	}
	if _, ok := extractJSONField(root, "a.missing"); ok {
		t.Fatal("expected missing path to be not-found")
	}
}
