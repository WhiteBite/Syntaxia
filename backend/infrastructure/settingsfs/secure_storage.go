package settingsfs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SecureStorage provides encrypted storage for sensitive credentials
type SecureStorage struct {
	credentialsPath string
	mu              sync.RWMutex
	machineKey      []byte
}

// Credential represents a stored credential
type Credential struct {
	Key       string `json:"key"`
	Value     string `json:"value"` // encrypted
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

// CredentialsFile represents the credentials file structure
type CredentialsFile struct {
	Version     int          `json:"version"`
	Credentials []Credential `json:"credentials"`
}

// NewSecureStorage creates a new secure storage instance
func NewSecureStorage() (*SecureStorage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".syntaxia")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	credentialsPath := filepath.Join(configDir, "credentials.json")

	// Generate machine-specific key for encryption
	machineKey := generateMachineKey()

	return &SecureStorage{
		credentialsPath: credentialsPath,
		machineKey:      machineKey,
	}, nil
}

// SaveCredential saves an encrypted credential
func (s *SecureStorage) SaveCredential(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Encrypt the value
	encryptedValue, err := s.encrypt(value)
	if err != nil {
		return fmt.Errorf("failed to encrypt credential: %w", err)
	}

	// Load existing credentials
	creds, err := s.loadCredentials()
	if err != nil {
		creds = &CredentialsFile{Version: 1, Credentials: []Credential{}}
	}

	// Update or add credential
	now := currentTimestamp()
	found := false
	for i, c := range creds.Credentials {
		if c.Key == key {
			creds.Credentials[i].Value = encryptedValue
			creds.Credentials[i].UpdatedAt = now
			found = true
			break
		}
	}

	if !found {
		creds.Credentials = append(creds.Credentials, Credential{
			Key:       key,
			Value:     encryptedValue,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	return s.saveCredentials(creds)
}

// LoadCredential loads and decrypts a credential
func (s *SecureStorage) LoadCredential(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	creds, err := s.loadCredentials()
	if err != nil {
		return "", fmt.Errorf("failed to load credentials: %w", err)
	}

	for _, c := range creds.Credentials {
		if c.Key == key {
			decrypted, err := s.decrypt(c.Value)
			if err != nil {
				return "", fmt.Errorf("failed to decrypt credential: %w", err)
			}
			return decrypted, nil
		}
	}

	return "", fmt.Errorf("credential not found: %s", key)
}

// DeleteCredential removes a credential
func (s *SecureStorage) DeleteCredential(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	creds, err := s.loadCredentials()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	newCreds := make([]Credential, 0, len(creds.Credentials))
	found := false
	for _, c := range creds.Credentials {
		if c.Key != key {
			newCreds = append(newCreds, c)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("credential not found: %s", key)
	}

	creds.Credentials = newCreds
	return s.saveCredentials(creds)
}

// HasCredential checks if a credential exists
func (s *SecureStorage) HasCredential(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	creds, err := s.loadCredentials()
	if err != nil {
		return false
	}

	for _, c := range creds.Credentials {
		if c.Key == key {
			return true
		}
	}
	return false
}

// ListCredentialKeys returns all credential keys (not values)
func (s *SecureStorage) ListCredentialKeys() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	creds, err := s.loadCredentials()
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(creds.Credentials))
	for i, c := range creds.Credentials {
		keys[i] = c.Key
	}
	return keys, nil
}

// loadCredentials loads credentials from file
func (s *SecureStorage) loadCredentials() (*CredentialsFile, error) {
	data, err := os.ReadFile(s.credentialsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &CredentialsFile{Version: 1, Credentials: []Credential{}}, nil
		}
		return nil, err
	}

	var creds CredentialsFile
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

// saveCredentials saves credentials to file
func (s *SecureStorage) saveCredentials(creds *CredentialsFile) error {
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	// Write with restricted permissions (owner only)
	return os.WriteFile(s.credentialsPath, data, 0600)
}

// encrypt encrypts a string using AES-GCM
func (s *SecureStorage) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(s.machineKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt decrypts a string using AES-GCM
func (s *SecureStorage) decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.machineKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// generateMachineKey generates a machine-specific encryption key
func generateMachineKey() []byte {
	// Use machine-specific data to generate a consistent key
	// This ensures credentials can only be decrypted on the same machine
	hostname, _ := os.Hostname()
	homeDir, _ := os.UserHomeDir()

	// Combine machine identifiers
	machineID := hostname + homeDir

	// Hash to get a 32-byte key for AES-256
	hash := sha256.Sum256([]byte(machineID))
	return hash[:]
}

// currentTimestamp returns current Unix timestamp
func currentTimestamp() int64 {
	return time.Now().Unix()
}
