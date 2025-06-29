package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"

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

// gen rand salt for spicing every keygen so that each key is unique even for same plaintext
// wanna make sure two peeps with same pass don't get the same key
// used before encryption
func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	return salt, err
}

// use pkb fn to gen key
func deriveKey(pass string, salt []byte) []byte {
	return pbkdf2.Key([]byte(pass), salt, iterations, keySize, sha256.New)

}

// essentially parse the plaintext into the struct
func Encrypt(plaintext []byte, pass string) (*EncryptedData, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}

	key := deriveKey(pass, salt)
	clearBytes(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)

	// makes sure encrypting same message twice with the same key gives diff ciphertext
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	return &EncryptedData{
		Salt:       salt,
		Nonce:      nonce,
		Ciphertext: ciphertext,
	}, nil

}

func Decrypt(encr EncryptedData, pass string) ([]byte, error) {
	key := deriveKey(pass, encr.Salt)
	defer clearBytes(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, encr.Nonce, encr.Ciphertext, nil)

	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// mk all bytes 0
func clearBytes(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
