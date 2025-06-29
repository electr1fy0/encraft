package main

import (
	"fmt"

	"github.com/electr1fy0/encraft/internal/crypto"
)

func testCrypto() {
	plaintext := []byte("Hello, this is a cat speaking")
	password := "meowforever70"

	fmt.Println("Testing encryption...")

	encrypted, err := crypto.Encrypt(plaintext, password)
	if err != nil {
		fmt.Println("Encryption failed: ", err)
		return
	}

	fmt.Println("Successful encryption")
	decrypted, err := crypto.Decrypt(*encrypted, password)
	if err != nil {
		fmt.Println("Decryption failed: ", err)

	}
	fmt.Printf("Original: %s\n", plaintext)
	fmt.Printf("Decrypted: %s\n", decrypted)

}

func main() {
	testCrypto()
}
