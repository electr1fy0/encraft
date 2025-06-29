package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/electr1fy0/encraft/internal/crypto"
)

type Entry struct {
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	Notes     string    `json:"notes,omitempty"`
	URL       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Vault struct {
	Version int               `json:"version"`
	Entries map[string]*Entry `json:"entries"`
}

func NewVault() *Vault {
	return &Vault{
		Version: 1,
		Entries: make(map[string]*Entry),
	}
}

func (v *Vault) AddEntry(entry *Entry) {
	now := time.Now()
	entry.CreatedAt = now
	entry.UpdatedAt = now
	v.Entries[entry.Name] = entry

}
func (v *Vault) GetEntry(name string) (*Entry, bool) {
	entry, exists := v.Entries[name]
	return entry, exists
}

func (v *Vault) DeleteEntry(name string) bool {
	if _, exists := v.Entries[name]; exists {
		delete(v.Entries, name)
		return true
	}
	return false

}

func (v *Vault) ListEntries() []string {
	names := make([]string, 0, len(v.Entries))

	for name, _ := range v.Entries {
		names = append(names, name)
	}
	return names
}

func (v *Vault) ToJSON() ([]byte, error) {
	return json.Marshal(v)
}

func FromJSON(data []byte) (*Vault, error) {
	var v Vault

	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return &v, nil

}

func GetVaultPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, ".secrets-vault"), nil
}

func VaultExists() (bool, error) {
	path, err := GetVaultPath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func SaveVault(vault *Vault, password string) error {
	jsonData, err := vault.ToJSON()
	if err != nil {
		return err
	}
	encryptedData, err := crypto.Encrypt(jsonData, password)
	if err != nil {
		return err
	}

	encryptedJSON, err := json.Marshal(encryptedData)
	if err != nil {
		return err
	}

	path, err := GetVaultPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, encryptedJSON, 0600)

}

func LoadVault(pass string) (*Vault, error) {
	path, err := GetVaultPath()
	if err != nil {
		return nil, err
	}

	encryptedJSON, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var encryptedData crypto.EncryptedData
	if err := json.Unmarshal(encryptedJSON, &encryptedData); err != nil {
		return nil, err
	}

	jsonData, err := crypto.Decrypt(encryptedData, pass)
	if err != nil {
		return nil, err
	}

	return FromJSON(jsonData)
}
