// backend/services/email_encryption_test.go
// Run with: cd backend && go test ./services/... -run TestPGP -v

package services

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

// ─── Service construction ─────────────────────────────────────────────────────

func TestNewPGPEncryptionService(t *testing.T) {
	svc := NewPGPEncryptionService()
	if svc == nil {
		t.Fatal("expected non-nil PGPEncryptionService")
	}
	if svc.keyStore == nil {
		t.Fatal("expected non-nil key store")
	}
}

func TestNewPGPKeyStore(t *testing.T) {
	ks := NewPGPKeyStore()
	if ks == nil {
		t.Fatal("expected non-nil PGPKeyStore")
	}
}

// ─── Key store – empty queries ────────────────────────────────────────────────

func TestGetKeyPairs_EmptyStore(t *testing.T) {
	svc := NewPGPEncryptionService()
	pairs, err := svc.GetKeyPairs(uuid.New())
	if err != nil {
		t.Fatalf("GetKeyPairs: %v", err)
	}
	if len(pairs) != 0 {
		t.Errorf("expected 0 key pairs on empty store, got %d", len(pairs))
	}
}

func TestGetPublicKeys_EmptyStore(t *testing.T) {
	svc := NewPGPEncryptionService()
	keys, err := svc.GetPublicKeys(uuid.New())
	if err != nil {
		t.Fatalf("GetPublicKeys: %v", err)
	}
	if len(keys) != 0 {
		t.Errorf("expected 0 public keys, got %d", len(keys))
	}
}

func TestGetKeyStats_EmptyStore(t *testing.T) {
	svc := NewPGPEncryptionService()
	stats, err := svc.GetKeyStats(uuid.New())
	if err != nil {
		t.Fatalf("GetKeyStats: %v", err)
	}
	if stats["key_pairs"] != 0 {
		t.Errorf("expected 0 key_pairs, got %v", stats["key_pairs"])
	}
	if stats["public_keys"] != 0 {
		t.Errorf("expected 0 public_keys, got %v", stats["public_keys"])
	}
}

func TestGetPublicKeyEntity_UnknownEmail(t *testing.T) {
	ks := NewPGPKeyStore()
	_, err := ks.GetPublicKeyEntity("nobody@example.com")
	if err == nil {
		t.Error("expected error for unknown email")
	}
}

func TestGetPrivateKeyEntity_UnknownUser(t *testing.T) {
	ks := NewPGPKeyStore()
	_, err := ks.GetPrivateKeyEntity(uuid.New(), "passphrase")
	if err == nil {
		t.Error("expected error for unknown user")
	}
}

// ─── encryptWithPassphrase / decryptWithPassphrase ────────────────────────────

func TestEncryptDecryptWithPassphrase_RoundTrip(t *testing.T) {
	svc := NewPGPEncryptionService()
	plaintext := []byte("top secret message")
	passphrase := "hunter2"

	ciphertext, err := svc.encryptWithPassphrase(plaintext, passphrase)
	if err != nil {
		t.Fatalf("encryptWithPassphrase: %v", err)
	}

	if string(ciphertext) == string(plaintext) {
		t.Fatal("ciphertext must differ from plaintext")
	}

	recovered, err := svc.decryptWithPassphrase(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("decryptWithPassphrase: %v", err)
	}

	if string(recovered) != string(plaintext) {
		t.Errorf("round-trip failed: got %q, want %q", recovered, plaintext)
	}
}

func TestEncryptWithPassphrase_DifferentOutputEachTime(t *testing.T) {
	svc := NewPGPEncryptionService()
	data := []byte("same data")
	pass := "samepass"

	ct1, err := svc.encryptWithPassphrase(data, pass)
	if err != nil {
		t.Fatal(err)
	}
	ct2, err := svc.encryptWithPassphrase(data, pass)
	if err != nil {
		t.Fatal(err)
	}

	if string(ct1) == string(ct2) {
		t.Error("expected random salt to produce different ciphertexts each time")
	}
}

func TestDecryptWithPassphrase_WrongPassphraseReturnsError(t *testing.T) {
	svc := NewPGPEncryptionService()
	ct, err := svc.encryptWithPassphrase([]byte("secret"), "correct-pass")
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.decryptWithPassphrase(ct, "wrong-pass")
	if err == nil {
		t.Error("expected error when decrypting with wrong passphrase")
	}
}

func TestDecryptWithPassphrase_TooShortReturnsError(t *testing.T) {
	svc := NewPGPEncryptionService()
	_, err := svc.decryptWithPassphrase([]byte("tooshort"), "pass")
	if err == nil {
		t.Error("expected error for too-short ciphertext")
	}
}

// ─── SaveKeyPair / GetKeyPairs ────────────────────────────────────────────────

func TestSaveKeyPair_ThenGetKeyPairs(t *testing.T) {
	ks := NewPGPKeyStore()
	userID := uuid.New()

	pair := &PGPKeyPair{
		ID:       uuid.New(),
		UserID:   userID,
		Name:     "Test User",
		Email:    "test@example.com",
		IsActive: true,
	}

	if err := ks.SaveKeyPair(pair); err != nil {
		t.Fatalf("SaveKeyPair: %v", err)
	}

	pairs, err := ks.GetKeyPairs(userID)
	if err != nil {
		t.Fatalf("GetKeyPairs: %v", err)
	}

	if len(pairs) != 1 {
		t.Fatalf("expected 1 key pair, got %d", len(pairs))
	}
	if pairs[0].Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %s", pairs[0].Email)
	}
}

func TestSaveKeyPair_MultipleUsers_IsolatesResults(t *testing.T) {
	ks := NewPGPKeyStore()
	user1 := uuid.New()
	user2 := uuid.New()

	_ = ks.SaveKeyPair(&PGPKeyPair{ID: uuid.New(), UserID: user1, Email: "u1@example.com", IsActive: true})
	_ = ks.SaveKeyPair(&PGPKeyPair{ID: uuid.New(), UserID: user1, Email: "u1b@example.com", IsActive: true})
	_ = ks.SaveKeyPair(&PGPKeyPair{ID: uuid.New(), UserID: user2, Email: "u2@example.com", IsActive: true})

	pairs1, _ := ks.GetKeyPairs(user1)
	pairs2, _ := ks.GetKeyPairs(user2)

	if len(pairs1) != 2 {
		t.Errorf("user1 expected 2 pairs, got %d", len(pairs1))
	}
	if len(pairs2) != 1 {
		t.Errorf("user2 expected 1 pair, got %d", len(pairs2))
	}
}

// ─── SavePublicKey / GetPublicKeys ────────────────────────────────────────────

func TestSavePublicKey_ThenGetPublicKeys(t *testing.T) {
	ks := NewPGPKeyStore()
	userID := uuid.New()

	info := &PublicKeyInfo{
		Email: "alice@example.com",
		KeyID: "AABBCCDD",
	}

	if err := ks.SavePublicKey(userID, "armored-key-data", info); err != nil {
		t.Fatalf("SavePublicKey: %v", err)
	}

	keys, err := ks.GetPublicKeys(userID)
	if err != nil {
		t.Fatalf("GetPublicKeys: %v", err)
	}

	if len(keys) != 1 {
		t.Fatalf("expected 1 public key, got %d", len(keys))
	}
	if keys[0].Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", keys[0].Email)
	}
	if keys[0].UserID != userID.String() {
		t.Errorf("expected userID %s, got %s", userID, keys[0].UserID)
	}
}

// ─── ExportPublicKey ──────────────────────────────────────────────────────────

func TestExportPublicKey_FoundAndNotFound(t *testing.T) {
	ks := NewPGPKeyStore()
	userID := uuid.New()
	keyID := uuid.New()

	pair := &PGPKeyPair{
		ID:        keyID,
		UserID:    userID,
		PublicKey: "-----BEGIN PGP PUBLIC KEY-----",
		IsActive:  true,
	}
	_ = ks.SaveKeyPair(pair)

	// Found case
	exported, err := ks.ExportPublicKey(userID, keyID.String())
	if err != nil {
		t.Fatalf("ExportPublicKey: %v", err)
	}
	if !strings.HasPrefix(exported, "-----BEGIN") {
		t.Errorf("unexpected key content: %s", exported)
	}

	// Not found case
	_, err = ks.ExportPublicKey(userID, uuid.New().String())
	if err == nil {
		t.Error("expected error for non-existent key")
	}
}

// ─── DeleteKeyPair ─────────────────────────────────────────────────────────────

func TestDeleteKeyPair_RemovesFromStore(t *testing.T) {
	ks := NewPGPKeyStore()
	userID := uuid.New()
	keyID := uuid.New()

	_ = ks.SaveKeyPair(&PGPKeyPair{ID: keyID, UserID: userID, IsActive: true})

	if err := ks.DeleteKeyPair(userID, keyID.String()); err != nil {
		t.Fatalf("DeleteKeyPair: %v", err)
	}

	pairs, _ := ks.GetKeyPairs(userID)
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs after delete, got %d", len(pairs))
	}
}

// ─── GetKeyStats ──────────────────────────────────────────────────────────────

func TestGetKeyStats_AfterAdding(t *testing.T) {
	svc := NewPGPEncryptionService()
	userID := uuid.New()

	_ = svc.keyStore.SaveKeyPair(&PGPKeyPair{ID: uuid.New(), UserID: userID, IsActive: true})
	_ = svc.keyStore.SaveKeyPair(&PGPKeyPair{ID: uuid.New(), UserID: userID, IsActive: true})
	_ = svc.keyStore.SavePublicKey(userID, "key-data", &PublicKeyInfo{Email: "a@b.com"})

	stats, err := svc.GetKeyStats(userID)
	if err != nil {
		t.Fatal(err)
	}

	if stats["key_pairs"] != 2 {
		t.Errorf("expected 2 key_pairs, got %v", stats["key_pairs"])
	}
	if stats["public_keys"] != 1 {
		t.Errorf("expected 1 public_key, got %v", stats["public_keys"])
	}
	if stats["total_keys"] != 3 {
		t.Errorf("expected total_keys=3, got %v", stats["total_keys"])
	}
}

// ─── encryptPrivateKey ─────────────────────────────────────────────────────────

func TestEncryptWithPassphrase_OutputLongerThanSalt(t *testing.T) {
	svc := NewPGPEncryptionService()

	// Output layout: 32-byte salt | nonce (12 bytes) | ciphertext+tag (≥ data + 16)
	data := []byte("private key DER bytes here")
	ct, err := svc.encryptWithPassphrase(data, "testpass")
	if err != nil {
		t.Fatal(err)
	}
	// Minimum: 32 (salt) + 12 (nonce) + 16 (GCM tag) + len(data)
	minExpected := 32 + 12 + 16 + len(data)
	if len(ct) < minExpected {
		t.Errorf("encrypted output too short: got %d bytes, want >= %d", len(ct), minExpected)
	}
}
		t.Errorf("encrypted output should be longer than just the salt: got %d bytes", len(ct))
	}
}
