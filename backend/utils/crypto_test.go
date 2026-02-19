// backend/utils/crypto_test.go
// Run with: cd backend && go test ./utils/... -v

package utils

import (
	"bytes"
	"testing"
)

// ─── KeyManager ──────────────────────────────────────────────────────────────

func TestNewKeyManager_CreatesSalt(t *testing.T) {
	km, err := NewKeyManager("password123")
	if err != nil {
		t.Fatalf("NewKeyManager failed: %v", err)
	}
	if len(km.GetSalt()) != 32 {
		t.Errorf("expected 32-byte salt, got %d", len(km.GetSalt()))
	}
}

func TestEncryptDecryptRoundTrip(t *testing.T) {
	km, err := NewKeyManager("super-secret")
	if err != nil {
		t.Fatalf("NewKeyManager: %v", err)
	}

	plaintext := []byte("Hello, TPT Titan!")

	ciphertext, err := km.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := km.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("round-trip failed: got %q, want %q", decrypted, plaintext)
	}
}

func TestDecrypt_ShortCiphertextReturnsError(t *testing.T) {
	km, err := NewKeyManager("pass")
	if err != nil {
		t.Fatal(err)
	}

	_, err = km.Decrypt([]byte("short"))
	if err == nil {
		t.Error("expected error for short ciphertext, got nil")
	}
}

func TestEncrypt_ProducesUniqueNonces(t *testing.T) {
	km, err := NewKeyManager("pass")
	if err != nil {
		t.Fatal(err)
	}
	msg := []byte("same message")

	ct1, _ := km.Encrypt(msg)
	ct2, _ := km.Encrypt(msg)

	if bytes.Equal(ct1, ct2) {
		t.Error("two encryptions of same message should produce different ciphertext (different nonces)")
	}
}

// ─── DeriveKeyFromPassword ────────────────────────────────────────────────────

func TestDeriveKeyFromPassword_ReproducesEncryption(t *testing.T) {
	km1, err := NewKeyManager("mypassword")
	if err != nil {
		t.Fatal(err)
	}

	plaintext := []byte("sensitive data")
	ciphertext, err := km1.Encrypt(plaintext)
	if err != nil {
		t.Fatal(err)
	}

	// Recreate using the same salt
	km2, err := DeriveKeyFromPassword("mypassword", km1.GetSalt())
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := km2.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("DeriveKeyFromPassword decrypt: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("got %q, want %q", decrypted, plaintext)
	}
}

func TestDeriveKeyFromPassword_WrongPasswordFails(t *testing.T) {
	km, err := NewKeyManager("correct-password")
	if err != nil {
		t.Fatal(err)
	}

	ciphertext, err := km.Encrypt([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}

	kmWrong, err := DeriveKeyFromPassword("wrong-password", km.GetSalt())
	if err != nil {
		t.Fatal(err)
	}

	_, err = kmWrong.Decrypt(ciphertext)
	if err == nil {
		t.Error("expected decryption to fail with wrong password")
	}
}

// ─── SecureString ─────────────────────────────────────────────────────────────

func TestSecureString_RoundTrip(t *testing.T) {
	original := "This is a secure message"
	password := "mypassword"

	ss, err := NewSecureString(original, password)
	if err != nil {
		t.Fatalf("NewSecureString: %v", err)
	}

	if ss.Algorithm != "AES-256-GCM" {
		t.Errorf("expected AES-256-GCM algorithm, got %s", ss.Algorithm)
	}

	recovered, err := ss.Decrypt(password)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	if recovered != original {
		t.Errorf("got %q, want %q", recovered, original)
	}
}

func TestSecureString_WrongPasswordFails(t *testing.T) {
	ss, err := NewSecureString("secret text", "correct-pass")
	if err != nil {
		t.Fatal(err)
	}

	_, err = ss.Decrypt("wrong-pass")
	if err == nil {
		t.Error("expected error with wrong password")
	}
}

func TestSecureString_Base64RoundTrip(t *testing.T) {
	original := "round-trip test"
	password := "testpass"

	ss, err := NewSecureString(original, password)
	if err != nil {
		t.Fatal(err)
	}

	b64 := ss.ToBase64()
	if b64 == "" {
		t.Fatal("ToBase64 returned empty string")
	}

	restored, err := FromBase64(b64)
	if err != nil {
		t.Fatalf("FromBase64: %v", err)
	}

	recovered, err := restored.Decrypt(password)
	if err != nil {
		t.Fatalf("Decrypt after FromBase64: %v", err)
	}

	if recovered != original {
		t.Errorf("got %q, want %q", recovered, original)
	}
}

func TestFromBase64_InvalidInputReturnsError(t *testing.T) {
	_, err := FromBase64("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestFromBase64_TooShortReturnsError(t *testing.T) {
	import64 := "dGVzdA==" // "test" base64 — only 4 bytes decoded, < 32
	_, err := FromBase64(import64)
	if err == nil {
		t.Error("expected error for too-short decoded data")
	}
}

// ─── ShamirSecretSharing ─────────────────────────────────────────────────────

func TestCreateShares_ThresholdGreaterThanTotal(t *testing.T) {
	_, err := CreateShares([]byte("secret"), 5, 3)
	if err == nil {
		t.Error("expected error when threshold > total shares")
	}
}

func TestCreateShares_ProducesCorrectCount(t *testing.T) {
	sss, err := CreateShares([]byte("my secret key!!"), 2, 3)
	if err != nil {
		t.Fatalf("CreateShares: %v", err)
	}

	shares := sss.GetShares()
	if len(shares) != 3 {
		t.Errorf("expected 3 shares, got %d", len(shares))
	}
}

func TestShamirReconstructSecret_Success(t *testing.T) {
	secret := []byte("abcdefghijklmnop") // 16 bytes for clean split

	sss, err := CreateShares(secret, 2, 2)
	if err != nil {
		t.Fatalf("CreateShares: %v", err)
	}

	shares := sss.GetShares()
	recovered, err := sss.ReconstructSecret(shares)
	if err != nil {
		t.Fatalf("ReconstructSecret: %v", err)
	}

	if !bytes.Equal(recovered, secret) {
		t.Errorf("reconstructed %q, want %q", recovered, secret)
	}
}

func TestShamirReconstructSecret_TooFewShares(t *testing.T) {
	sss, err := CreateShares([]byte("secret1234567890"), 3, 5)
	if err != nil {
		t.Fatal(err)
	}

	shares := sss.GetShares()
	_, err = sss.ReconstructSecret(shares[:2]) // Only 2 of required 3
	if err == nil {
		t.Error("expected error with insufficient shares")
	}
}

// ─── compareFaceTemplates ──────────────────────────────────────────────────

func TestCompareFaceTemplates_Match(t *testing.T) {
	template := []byte{0x01, 0x02, 0x03, 0x04}
	if !compareFaceTemplates(template, template) {
		t.Error("identical templates should match")
	}
}

func TestCompareFaceTemplates_NoMatch(t *testing.T) {
	t1 := []byte{0x01, 0x02, 0x03}
	t2 := []byte{0x01, 0x02, 0x04}
	if compareFaceTemplates(t1, t2) {
		t.Error("different templates should not match")
	}
}

func TestCompareFaceTemplates_DifferentLength(t *testing.T) {
	t1 := []byte{0x01, 0x02}
	t2 := []byte{0x01, 0x02, 0x03}
	if compareFaceTemplates(t1, t2) {
		t.Error("different-length templates should not match")
	}
}

// ─── HardwareSecurityKey ──────────────────────────────────────────────────────

func TestCreateHardwareKey_Success(t *testing.T) {
	faceTemplate := []byte("fake-face-biometric-data")
	encKey := []byte("my-encryption-key-32bytesxxxxxxx")

	hsk, err := CreateHardwareKey("device-001", encKey, faceTemplate)
	if err != nil {
		t.Fatalf("CreateHardwareKey: %v", err)
	}

	if hsk == nil {
		t.Fatal("expected non-nil HardwareSecurityKey")
	}
	if len(hsk.publicKey) == 0 {
		t.Error("expected non-empty public key")
	}
	if len(hsk.encryptedKey) == 0 {
		t.Error("expected non-empty encrypted key")
	}
}

func TestRecoverFromHardware_WrongFaceFails(t *testing.T) {
	faceTemplate := []byte("correct-face-template")
	encKey := []byte("my-encryption-key-32bytesxxxxxxx")

	hsk, err := CreateHardwareKey("device-001", encKey, faceTemplate)
	if err != nil {
		t.Fatal(err)
	}

	_, err = hsk.RecoverFromHardware([]byte("wrong-face-template"))
	if err == nil {
		t.Error("expected error for wrong biometric scan")
	}
}
