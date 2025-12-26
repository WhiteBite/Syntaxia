package settingsfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveCredential_Success(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()
	storage := &SecureStorage{
		credentialsPath: filepath.Join(tempDir, "credentials.json"),
		machineKey:      generateMachineKey(),
	}

	err := storage.SaveCredential("api_key", "sk-test-12345")
	if err != nil {
		t.Fatalf("SaveCredential failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(storage.credentialsPath); os.IsNotExist(err) {
		t.Fatal("credentials file was not created")
	}
}

func TestLoadCredential_Success(t *testing.T) {
	tempDir := t.TempDir()
	storage := &SecureStorage{
		credentialsPath: filepath.Join(tempDir, "credentials.json"),
		machineKey:      generateMachineKey(),
	}

	// Save a credential first
	originalValue := "sk-secret-key-67890"
	err := storage.SaveCredential("openai_key", originalValue)
	if err != nil {
		t.Fatalf("SaveCredential failed: %v", err)
	}

	// Load it back
	loadedValue, err := storage.LoadCredential("openai_key")
	if err != nil {
		t.Fatalf("LoadCredential failed: %v", err)
	}

	if loadedValue != originalValue {
		t.Errorf("LoadCredential returned %q, want %q", loadedValue, originalValue)
	}
}

func TestLoadCredential_NotFound(t *testing.T) {
	tempDir := t.TempDir()
	storage := &SecureStorage{
		credentialsPath: filepath.Join(tempDir, "credentials.json"),
		machineKey:      generateMachineKey(),
	}

	_, err := storage.LoadCredential("nonexistent_key")
	if err == nil {
		t.Fatal("LoadCredential should return error for nonexistent key")
	}
}

func TestDeleteCredential_Success(t *testing.T) {
	tempDir := t.TempDir()
	storage := &SecureStorage{
		credentialsPath: filepath.Join(tempDir, "credentials.json"),
		machineKey:      generateMachineKey(),
	}

	// Save a credential
	err := storage.SaveCredential("temp_key", "temp_value")
	if err != nil {
		t.Fatalf("SaveCredential failed: %v", err)
	}

	// Verify it exists
	if !storage.HasCredential("temp_key") {
		t.Fatal("credential should exist after save")
	}

	// Delete it
	err = storage.DeleteCredential("temp_key")
	if err != nil {
		t.Fatalf("DeleteCredential failed: %v", err)
	}

	// Verify it's gone
	if storage.HasCredential("temp_key") {
		t.Fatal("credential should not exist after delete")
	}
}

func TestHasCredential(t *testing.T) {
	tempDir := t.TempDir()
	storage := &SecureStorage{
		credentialsPath: filepath.Join(tempDir, "credentials.json"),
		machineKey:      generateMachineKey(),
	}

	// Should not exist initially
	if storage.HasCredential("test_key") {
		t.Fatal("credential should not exist initially")
	}

	// Save it
	_ = storage.SaveCredential("test_key", "test_value")

	// Should exist now
	if !storage.HasCredential("test_key") {
		t.Fatal("credential should exist after save")
	}
}

func TestListCredentialKeys(t *testing.T) {
	tempDir := t.TempDir()
	storage := &SecureStorage{
		credentialsPath: filepath.Join(tempDir, "credentials.json"),
		machineKey:      generateMachineKey(),
	}

	// Save multiple credentials
	_ = storage.SaveCredential("key1", "value1")
	_ = storage.SaveCredential("key2", "value2")
	_ = storage.SaveCredential("key3", "value3")

	keys, err := storage.ListCredentialKeys()
	if err != nil {
		t.Fatalf("ListCredentialKeys failed: %v", err)
	}

	if len(keys) != 3 {
		t.Errorf("ListCredentialKeys returned %d keys, want 3", len(keys))
	}
}
