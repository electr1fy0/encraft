package storage

import (
	"time"
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
	return &Vault {
		Version : 1,
		Entries : make(map[string]*Entry),
	}
}

func (v *Vault) AddEntry(entry * Entry){

}

func (v *Vault) GetEntry(name string) (* Entry) {

}

func
