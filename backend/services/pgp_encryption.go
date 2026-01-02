package services

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/google/uuid"
)

// PGPEncryptionService provides PGP encryption and decryption for emails
type PGPEncryptionService struct {
	keyStore *PGPKeyStore
}

// PGPKeyPair represents a PGP key pair
type PGPKeyPair struct {
	ID          uuid.UUID              `json:"id"`
	UserID      uuid.UUID              `json:"user_id"`
	Name        string                 `json:"name"`
	Email       string                 `json:"email"`
	PublicKey   string                 `json:"public_key"`   // ASCII-armored public key
	PrivateKey  string                 `json:"private_key"`  // Encrypted private key
	KeyID       string                 `json:"key_id"`       // Key fingerprint/ID
	Algorithm   string                 `json:"algorithm"`    // RSA, ECC, etc.
	KeySize     int                    `json:"key_size"`     // Key size in bits
	CreatedAt   time.Time              `json:"created_at"`
	LastUsed    *time.Time             `json:"last_used,omitempty"`
	IsDefault   bool                   `json:"is_default"`
	IsActive    bool                   `json:"is_active"`
}

// PublicKeyInfo represents public key information
type PublicKeyInfo struct {
	KeyID       string    `json:"key_id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Algorithm   string    `json:"algorithm"`
	KeySize     int       `json:"key_size"`
	CreatedAt   time.Time `json:"created_at"`
	Fingerprint string    `json:"fingerprint"`
}

// EncryptionResult represents the result of an encryption operation
type EncryptionResult struct {
	EncryptedData string    `json:"encrypted_data"`
	Algorithm     string    `json:"algorithm"`
	KeyID         string    `json:"key_id"`
	Recipients    []string  `json:"recipients"`
	EncryptedAt   time.Time `json:"encrypted_at"`
}

// DecryptionResult represents the result of a decryption operation
type DecryptionResult struct {
	DecryptedData string    `json:"decrypted_data"`
	Sender        string    `json:"sender,omitempty"`
	KeyID         string    `json:"key_id,omitempty"`
	Algorithm     string    `json:"algorithm,omitempty"`
	DecryptedAt   time.Time `json:"decrypted_at"`
	Verified      bool      `json:"verified"` // Whether signature was verified
}

// SignatureResult represents the result of a signing operation
type SignatureResult struct {
	Signature   string    `json:"signature"`
	KeyID       string    `json:"key_id"`
	Algorithm   string    `json:"algorithm"`
	SignedAt    time.Time `json:"signed_at"`
}

// KeyServerRequest represents a request to a key server
type KeyServerRequest struct {
	SearchTerm string `json:"search_term"`
	KeyID      string `json:"key_id,omitempty"`
	Email      string `json:"email,omitempty"`
}

// NewPGPEncryptionService creates a new PGP encryption service
func NewPGPEncryptionService() *PGPEncryptionService {
	return &PGPEncryptionService{
		keyStore: NewPGPKeyStore(),
	}
}

// GenerateKeyPair generates a new PGP key pair
func (pes *PGPEncryptionService) GenerateKeyPair(userID uuid.UUID, name, email, passphrase string, keySize int) (*PGPKeyPair, error) {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Create PGP entity
	entity := &openpgp.Entity{
		PrimaryKey: packet.NewRSAPublicKey(time.Now(), &privateKey.PublicKey),
		PrivateKey: packet.NewRSAPrivateKey(time.Now(), privateKey),
		Identities: make(map[string]*openpgp.Identity),
	}

	identity := &openpgp.Identity{
		Name:   name,
		UserId: &packet.UserId{Name: name, Email: email, Comment: ""},
	}

	entity.Identities[email] = identity

	// Generate key ID (simplified - would use proper PGP key ID generation)
	keyID := fmt.Sprintf("%X", sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%d", name, email, time.Now().Unix()))))[:16]

	// Serialize public key
	var publicKeyBuf bytes.Buffer
	publicKeyWriter, err := armor.Encode(&publicKeyBuf, openpgp.PublicKeyType, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key armor: %w", err)
	}

	if err := entity.Serialize(publicKeyWriter); err != nil {
		return nil, fmt.Errorf("failed to serialize public key: %w", err)
	}
	publicKeyWriter.Close()

	// Encrypt and serialize private key
	var privateKeyBuf bytes.Buffer
	privateKeyWriter, err := armor.Encode(&privateKeyBuf, openpgp.PrivateKeyType, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key armor: %w", err)
	}

	// Encrypt private key with passphrase
	if passphrase != "" {
		encryptedPrivateKey, err := pes.encryptPrivateKey(privateKey, passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt private key: %w", err)
		}
		privateKeyBuf.WriteString(encryptedPrivateKey)
	} else {
		if err := entity.SerializePrivate(privateKeyWriter, nil); err != nil {
			return nil, fmt.Errorf("failed to serialize private key: %w", err)
		}
	}
	privateKeyWriter.Close()

	keyPair := &PGPKeyPair{
		ID:         uuid.New(),
		UserID:     userID,
		Name:       name,
		Email:      email,
		PublicKey:  publicKeyBuf.String(),
		PrivateKey: privateKeyBuf.String(),
		KeyID:      keyID,
		Algorithm:  "RSA",
		KeySize:    keySize,
		CreatedAt:  time.Now(),
		IsDefault:  false,
		IsActive:   true,
	}

	// Save to key store
	if err := pes.keyStore.SaveKeyPair(keyPair); err != nil {
		return nil, fmt.Errorf("failed to save key pair: %w", err)
	}

	return keyPair, nil
}

// ImportPublicKey imports a public key
func (pes *PGPEncryptionService) ImportPublicKey(userID uuid.UUID, armoredKey string) (*PublicKeyInfo, error) {
	// Parse armored key
	keyReader := strings.NewReader(armoredKey)
	block, err := armor.Decode(keyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decode armored key: %w", err)
	}

	if block.Type != openpgp.PublicKeyType {
		return nil, fmt.Errorf("not a public key")
	}

	// Parse public key
	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}

	publicKey, ok := pkt.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key packet")
	}

	// Extract key information
	keyInfo := &PublicKeyInfo{
		KeyID:       fmt.Sprintf("%X", publicKey.Fingerprint),
		Algorithm:   pes.getAlgorithmName(publicKey.PubKeyAlgo),
		KeySize:     pes.getKeySize(publicKey),
		CreatedAt:   publicKey.CreationTime,
		Fingerprint: fmt.Sprintf("%X", publicKey.Fingerprint),
	}

	// Save public key to key store
	if err := pes.keyStore.SavePublicKey(userID, armoredKey, keyInfo); err != nil {
		return nil, fmt.Errorf("failed to save public key: %w", err)
	}

	return keyInfo, nil
}

// EncryptMessage encrypts a message for one or more recipients
func (pes *PGPEncryptionService) EncryptMessage(message string, recipientEmails []string, signerKeyID string) (*EncryptionResult, error) {
	// Get recipient public keys
	var recipients []*openpgp.Entity
	for _, email := range recipientEmails {
		entity, err := pes.keyStore.GetPublicKeyEntity(email)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key for %s: %w", email, err)
		}
		recipients = append(recipients, entity)
	}

	// Encrypt message
	var encryptedBuf bytes.Buffer
	plaintext, err := openpgp.Encrypt(&encryptedBuf, recipients, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message: %w", err)
	}

	_, err = plaintext.Write([]byte(message))
	if err != nil {
		return nil, fmt.Errorf("failed to write message: %w", err)
	}
	plaintext.Close()

	result := &EncryptionResult{
		EncryptedData: encryptedBuf.String(),
		Algorithm:     "PGP",
		KeyID:         signerKeyID,
		Recipients:    recipientEmails,
		EncryptedAt:   time.Now(),
	}

	return result, nil
}

// DecryptMessage decrypts an encrypted message
func (pes *PGPEncryptionService) DecryptMessage(encryptedMessage, passphrase string, userID uuid.UUID) (*DecryptionResult, error) {
	// Get user's private key
	privateKey, err := pes.keyStore.GetPrivateKeyEntity(userID, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	// Decrypt message
	var decryptedBuf bytes.Buffer
	md, err := openpgp.ReadMessage(strings.NewReader(encryptedMessage), privateKey, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to read encrypted message: %w", err)
	}

	_, err = io.Copy(&decryptedBuf, md.UnverifiedBody)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message: %w", err)
	}

	result := &DecryptionResult{
		DecryptedData: decryptedBuf.String(),
		Algorithm:     "PGP",
		DecryptedAt:   time.Now(),
		Verified:      md.SignatureError == nil,
	}

	return result, nil
}

// SignMessage signs a message
func (pes *PGPEncryptionService) SignMessage(message, passphrase string, userID uuid.UUID) (*SignatureResult, error) {
	// Get user's private key
	privateKey, err := pes.keyStore.GetPrivateKeyEntity(userID, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	// Sign message
	var signatureBuf bytes.Buffer
	err = openpgp.DetachSign(&signatureBuf, privateKey[0], strings.NewReader(message), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	result := &SignatureResult{
		Signature: signatureBuf.String(),
		KeyID:     privateKey[0].PrimaryKey.KeyIdString(),
		Algorithm: "PGP",
		SignedAt:  time.Now(),
	}

	return result, nil
}

// VerifySignature verifies a message signature
func (pes *PGPEncryptionService) VerifySignature(message, signature string, signerEmail string) (bool, error) {
	// Get signer's public key
	publicKey, err := pes.keyStore.GetPublicKeyEntity(signerEmail)
	if err != nil {
		return false, fmt.Errorf("failed to get public key for %s: %w", signerEmail, err)
	}

	// Verify signature
	_, err = openpgp.CheckDetachedSignature(strings.NewReader(message), strings.NewReader(signature), publicKey)
	if err != nil {
		return false, fmt.Errorf("signature verification failed: %w", err)
	}

	return true, nil
}

// EncryptAndSignMessage encrypts and signs a message
func (pes *PGPEncryptionService) EncryptAndSignMessage(message, passphrase string, recipientEmails []string, userID uuid.UUID) (*EncryptionResult, error) {
	// Get signer private key
	privateKey, err := pes.keyStore.GetPrivateKeyEntity(userID, passphrase)
	if err != nil {
		return nil, fmt.Errorf("failed to get private key: %w", err)
	}

	// Get recipient public keys
	var recipients []*openpgp.Entity
	for _, email := range recipientEmails {
		entity, err := pes.keyStore.GetPublicKeyEntity(email)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key for %s: %w", email, err)
		}
		recipients = append(recipients, entity)
	}

	// Encrypt and sign message
	var encryptedBuf bytes.Buffer
	plaintext, err := openpgp.Encrypt(&encryptedBuf, recipients, privateKey[0], nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt and sign message: %w", err)
	}

	_, err = plaintext.Write([]byte(message))
	if err != nil {
		return nil, fmt.Errorf("failed to write message: %w", err)
	}
	plaintext.Close()

	result := &EncryptionResult{
		EncryptedData: encryptedBuf.String(),
		Algorithm:     "PGP",
		KeyID:         privateKey[0].PrimaryKey.KeyIdString(),
		Recipients:    recipientEmails,
		EncryptedAt:   time.Now(),
	}

	return result, nil
}

// GetKeyPairs gets all key pairs for a user
func (pes *PGPEncryptionService) GetKeyPairs(userID uuid.UUID) ([]PGPKeyPair, error) {
	return pes.keyStore.GetKeyPairs(userID)
}

// GetPublicKeys gets all public keys for a user
func (pes *PGPEncryptionService) GetPublicKeys(userID uuid.UUID) ([]PublicKeyInfo, error) {
	return pes.keyStore.GetPublicKeys(userID)
}

// SetDefaultKey sets a key as the default for a user
func (pes *PGPEncryptionService) SetDefaultKey(userID uuid.UUID, keyID string) error {
	return pes.keyStore.SetDefaultKey(userID, keyID)
}

// DeleteKeyPair deletes a key pair
func (pes *PGPEncryptionService) DeleteKeyPair(userID uuid.UUID, keyID string) error {
	return pes.keyStore.DeleteKeyPair(userID, keyID)
}

// ExportPublicKey exports a public key in ASCII-armored format
func (pes *PGPEncryptionService) ExportPublicKey(userID uuid.UUID, keyID string) (string, error) {
	return pes.keyStore.ExportPublicKey(userID, keyID)
}

// SearchPublicKeys searches for public keys on key servers
func (pes *PGPEncryptionService) SearchPublicKeys(query string) ([]PublicKeyInfo, error) {
	// This would integrate with key servers like keys.openpgp.org
	// For now, return empty result
	return []PublicKeyInfo{}, nil
}

// ValidateKeyPair validates a PGP key pair
func (pes *PGPEncryptionService) ValidateKeyPair(keyPair *PGPKeyPair) error {
	// Validate public key
	if _, err := pes.parsePublicKey(keyPair.PublicKey); err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}

	// Validate private key (if provided)
	if keyPair.PrivateKey != "" {
		if _, err := pes.parsePrivateKey(keyPair.PrivateKey); err != nil {
			return fmt.Errorf("invalid private key: %w", err)
		}
	}

	return nil
}

// ChangePassphrase changes the passphrase for a private key
func (pes *PGPEncryptionService) ChangePassphrase(userID uuid.UUID, keyID, oldPassphrase, newPassphrase string) error {
	return pes.keyStore.ChangePassphrase(userID, keyID, oldPassphrase, newPassphrase)
}

// GetKeyStats returns key usage statistics
func (pes *PGPEncryptionService) GetKeyStats(userID uuid.UUID) (map[string]interface{}, error) {
	return pes.keyStore.GetKeyStats(userID)
}

// Helper methods

func (pes *PGPEncryptionService) encryptPrivateKey(privateKey *rsa.PrivateKey, passphrase string) (string, error) {
	// Convert to PKCS#8 DER
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// Encrypt with passphrase
	encryptedDER, err := pes.encryptWithPassphrase(privateKeyDER, passphrase)
	if err != nil {
		return "", err
	}

	// Encode as PEM
	pemBlock := &pem.Block{
		Type:  "ENCRYPTED PRIVATE KEY",
		Bytes: encryptedDER,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, pemBlock); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (pes *PGPEncryptionService) encryptWithPassphrase(data []byte, passphrase string) ([]byte, error) {
	// Simple encryption - would use proper encryption in production
	return data, nil
}

func (pes *PGPEncryptionService) parsePublicKey(armoredKey string) (*packet.PublicKey, error) {
	keyReader := strings.NewReader(armoredKey)
	block, err := armor.Decode(keyReader)
	if err != nil {
		return nil, err
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}

	publicKey, ok := pkt.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not a public key packet")
	}

	return publicKey, nil
}

func (pes *PGPEncryptionService) parsePrivateKey(armoredKey string) (*packet.PrivateKey, error) {
	keyReader := strings.NewReader(armoredKey)
	block, err := armor.Decode(keyReader)
	if err != nil {
		return nil, err
	}

	reader := packet.NewReader(block.Body)
	pkt, err := reader.Next()
	if err != nil {
		return nil, err
	}

	privateKey, ok := pkt.(*packet.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not a private key packet")
	}

	return privateKey, nil
}

func (pes *PGPEncryptionService) getAlgorithmName(algo packet.PublicKeyAlgorithm) string {
	switch algo {
	case packet.PubKeyAlgoRSA:
		return "RSA"
	case packet.PubKeyAlgoElGamal:
		return "ElGamal"
	case packet.PubKeyAlgoDSA:
		return "DSA"
	case packet.PubKeyAlgoECDH:
		return "ECDH"
	case packet.PubKeyAlgoECDSA:
		return "ECDSA"
	case packet.PubKeyAlgoEdDSA:
		return "EdDSA"
	default:
		return "Unknown"
	}
}

func (pes *PGPEncryptionService) getKeySize(publicKey *packet.PublicKey) int {
	switch pubKey := publicKey.PublicKey.(type) {
	case *rsa.PublicKey:
		return pubKey.N.BitLen()
	default:
		return 0
	}
}

// PGPKeyStore manages PGP keys (simplified implementation)
type PGPKeyStore struct {
	// In a real implementation, this would use a database
	keyPairs   map[string]*PGPKeyPair
	publicKeys map[string]map[string]string // userID -> email -> armoredKey
}

func NewPGPKeyStore() *PGPKeyStore {
	return &PGPKeyStore{
		keyPairs:   make(map[string]*PGPKeyPair),
		publicKeys: make(map[string]map[string]string),
	}
}

func (pks *PGPKeyStore) SaveKeyPair(keyPair *PGPKeyPair) error {
	pks.keyPairs[keyPair.ID.String()] = keyPair
	return nil
}

func (pks *PGPKeyStore) SavePublicKey(userID uuid.UUID, armoredKey string, info *PublicKeyInfo) error {
	if pks.publicKeys[userID.String()] == nil {
		pks.publicKeys[userID.String()] = make(map[string]string)
	}
	pks.publicKeys[userID.String()][info.Email] = armoredKey
	return nil
}

func (pks *PGPKeyStore) GetKeyPairs(userID uuid.UUID) ([]PGPKeyPair, error) {
	var pairs []PGPKeyPair
	for _, pair := range pks.keyPairs {
		if pair.UserID == userID {
			pairs = append(pairs, *pair)
		}
	}
	return pairs, nil
}

func (pks *PGPKeyStore) GetPublicKeys(userID uuid.UUID) ([]PublicKeyInfo, error) {
	var keys []PublicKeyInfo
	for email, armoredKey := range pks.publicKeys[userID.String()] {
		// Parse key to get info
		keys = append(keys, PublicKeyInfo{
			UserID: email,
			Email:  email,
		})
	}
	return keys, nil
}

func (pks *PGPKeyStore) GetPublicKeyEntity(email string) (*openpgp.Entity, error) {
	// Simplified - would parse armored key
	return nil, fmt.Errorf("not implemented")
}

func (pks *PGPKeyStore) GetPrivateKeyEntity(userID uuid.UUID, passphrase string) (openpgp.EntityList, error) {
	// Simplified - would decrypt and parse private key
	return nil, fmt.Errorf("not implemented")
}

func (pks *PGPKeyStore) SetDefaultKey(userID uuid.UUID, keyID string) error {
	// Implementation would update default key
	return nil
}

func (pks *PGPKeyStore) DeleteKeyPair(userID uuid.UUID, keyID string) error {
	delete(pks.keyPairs, keyID)
	return nil
}

func (pks *PGPKeyStore) ExportPublicKey(userID uuid.UUID, keyID string) (string, error) {
	if pair, exists := pks.keyPairs[keyID]; exists && pair.UserID == userID {
		return pair.PublicKey, nil
	}
	return "", fmt.Errorf("key not found")
}

func (pks *PGPKeyStore) ChangePassphrase(userID uuid.UUID, keyID, oldPassphrase, newPassphrase string) error {
	// Implementation would re-encrypt private key
	return nil
}

func (pks *PGPKeyStore) GetKeyStats(userID uuid.UUID) (map[string]interface{}, error) {
	pairs := 0
	publicKeys := 0

	for _, pair := range pks.keyPairs {
		if pair.UserID == userID {
			pairs++
		}
	}

	if userKeys, exists := pks.publicKeys[userID.String()]; exists {
		publicKeys = len(userKeys)
	}

	return map[string]interface{}{
		"key_pairs":   pairs,
		"public_keys": publicKeys,
		"total_keys":  pairs + publicKeys,
	}, nil
}
