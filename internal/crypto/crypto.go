package crypto

import (
	"crypto/sha256"
	// "math/rand"
	"crypto/rand"

	"golang.org/x/crypto/pbkdf2" // password based key derivation function 2 (phew)
)

const (
	saltSize   = 32
	keySize    = 32
	iterations = 10_000
)

type EncryptedData struct {
	Salt       []byte `json:"salt"`
	Nonce      []byte `json:"nonce"`
	Ciphertext []byte `json:"ciphertext"`
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	return salt, err

}

func deriveKey(pass string, salt []byte) []byte {
	return pbkdf2.Key([]byte(pass), salt, iterations, keySize, sha256.New)

}

func Encrypt(plaintext []byte) (*EncryptedData, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

}
