package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/scrypt"
)

// KeyManager handles user encryption keys with zero-knowledge architecture
type KeyManager struct {
	masterKey []byte // Derived from user password, never stored
	salt      []byte // Random salt for key derivation
}

// NewKeyManager creates a new key manager with Argon2 key derivation
func NewKeyManager(password string) (*KeyManager, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	// Argon2id key derivation (recommended for password hashing)
	masterKey := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return &KeyManager{
		masterKey: masterKey,
		salt:      salt,
	}, nil
}

// Encrypt encrypts data using AES-256-GCM with the master key
func (km *KeyManager) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(km.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts data using AES-256-GCM with the master key
func (km *KeyManager) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(km.masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// GetSalt returns the salt for key derivation (safe to store)
func (km *KeyManager) GetSalt() []byte {
	return km.salt
}

// DeriveKeyFromPassword recreates the key manager from password and stored salt
func DeriveKeyFromPassword(password string, salt []byte) (*KeyManager, error) {
	masterKey := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	return &KeyManager{
		masterKey: masterKey,
		salt:      salt,
	}, nil
}

// ShamirSecretSharing implements threshold cryptography for key recovery
type ShamirSecretSharing struct {
	secret []byte
	shares [][]byte
	threshold int
	totalShares int
}

// CreateShares splits a secret into shares using Shamir's algorithm
func CreateShares(secret []byte, threshold, totalShares int) (*ShamirSecretSharing, error) {
	if threshold > totalShares {
		return nil, errors.New("threshold cannot be greater than total shares")
	}

	// For simplicity, implement a basic threshold scheme
	// In production, use a proper implementation like github.com/hashicorp/vault/shamir

	shares := make([][]byte, totalShares)

	// Split the secret into equal parts (simplified version)
	shareSize := len(secret) / totalShares
	remainder := len(secret) % totalShares

	for i := 0; i < totalShares; i++ {
		start := i * shareSize
		end := start + shareSize
		if i < remainder {
			end++
		}
		shares[i] = make([]byte, end-start)
		copy(shares[i], secret[start:end])
	}

	return &ShamirSecretSharing{
		secret:      secret,
		shares:      shares,
		threshold:   threshold,
		totalShares: totalShares,
	}, nil
}

// ReconstructSecret combines shares to recover the secret
func (sss *ShamirSecretSharing) ReconstructSecret(shares [][]byte) ([]byte, error) {
	if len(shares) < sss.threshold {
		return nil, fmt.Errorf("need at least %d shares, got %d", sss.threshold, len(shares))
	}

	// Combine the shares back (simplified version)
	var secret []byte
	for _, share := range shares {
		secret = append(secret, share...)
	}

	// Verify the reconstruction (compare hash)
	originalHash := sha256.Sum256(sss.secret)
	reconstructedHash := sha256.Sum256(secret)

	if originalHash != reconstructedHash {
		return nil, errors.New("share reconstruction failed - invalid shares")
	}

	return secret, nil
}

// GetShares returns the generated shares
func (sss *ShamirSecretSharing) GetShares() [][]byte {
	return sss.shares
}

// HardwareSecurityKey represents a USB drive + face ID recovery system
type HardwareSecurityKey struct {
	deviceID   string
	publicKey  []byte
	encryptedKey []byte
	faceTemplate []byte // Facial recognition template
}

// CreateHardwareKey initializes a hardware security key
func CreateHardwareKey(deviceID string, encryptionKey []byte, faceTemplate []byte) (*HardwareSecurityKey, error) {
	// Generate a keypair for the device
	publicKey := make([]byte, 32)
	if _, err := rand.Read(publicKey); err != nil {
		return nil, fmt.Errorf("failed to generate device key: %w", err)
	}

	// Encrypt the user's encryption key with the device public key
	km, err := NewKeyManager(string(publicKey))
	if err != nil {
		return nil, err
	}

	encryptedKey, err := km.Encrypt(encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt key: %w", err)
	}

	return &HardwareSecurityKey{
		deviceID:     deviceID,
		publicKey:    publicKey,
		encryptedKey: encryptedKey,
		faceTemplate: faceTemplate,
	}, nil
}

// RecoverFromHardware recovers the encryption key using hardware + biometrics
func (hsk *HardwareSecurityKey) RecoverFromHardware(faceScan []byte) ([]byte, error) {
	// Verify face template (simplified - in reality use proper biometric matching)
	if !compareFaceTemplates(hsk.faceTemplate, faceScan) {
		return nil, errors.New("biometric verification failed")
	}

	// Decrypt the key using the device private key (simplified)
	km, err := DeriveKeyFromPassword(string(hsk.publicKey), []byte("device_private_key"))
	if err != nil {
		return nil, err
	}

	decryptedKey, err := km.Decrypt(hsk.encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}

	return decryptedKey, nil
}

// compareFaceTemplates performs biometric comparison (placeholder implementation)
func compareFaceTemplates(template1, template2 []byte) bool {
	// In reality, use proper biometric comparison algorithm
	// For demo, simple byte comparison
	if len(template1) != len(template2) {
		return false
	}
	for i := range template1 {
		if template1[i] != template2[i] {
			return false
		}
	}
	return true
}

// SecureString represents encrypted text that can be safely stored/transmitted
type SecureString struct {
	EncryptedData []byte
	Salt          []byte
	Algorithm     string
}

// NewSecureString creates an encrypted string from plaintext
func NewSecureString(plaintext, password string) (*SecureString, error) {
	km, err := NewKeyManager(password)
	if err != nil {
		return nil, err
	}

	encrypted, err := km.Encrypt([]byte(plaintext))
	if err != nil {
		return nil, err
	}

	return &SecureString{
		EncryptedData: encrypted,
		Salt:          km.GetSalt(),
		Algorithm:     "AES-256-GCM",
	}, nil
}

// Decrypt recovers the plaintext from a secure string
func (ss *SecureString) Decrypt(password string) (string, error) {
	km, err := DeriveKeyFromPassword(password, ss.Salt)
	if err != nil {
		return "", err
	}

	decrypted, err := km.Decrypt(ss.EncryptedData)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// ToBase64 converts secure string to base64 for storage/transmission
func (ss *SecureString) ToBase64() string {
	data := append(ss.Salt, ss.EncryptedData...)
	return base64.StdEncoding.EncodeToString(data)
}

// FromBase64 reconstructs secure string from base64
func FromBase64(data string) (*SecureString, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	if len(decoded) < 32 {
		return nil, errors.New("invalid secure string format")
	}

	return &SecureString{
		EncryptedData: decoded[32:],
		Salt:          decoded[:32],
		Algorithm:     "AES-256-GCM",
	}, nil
}
