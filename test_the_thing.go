// package main

// import (
// 	"fmt"

// 	"github.com/electr1fy0/encraft/internal/crypto"
// 	"github.com/electr1fy0/encraft/storage"
// )

// func testCrypto() {
// 	plaintext := []byte("Hello, this is a cat speaking")
// 	password := "meowforever70"

// 	fmt.Println("Testing encryption...")

// 	encrypted, err := crypto.Encrypt(plaintext, password)
// 	if err != nil {
// 		fmt.Println("Encryption failed: ", err)
// 		return
// 	}

// 	fmt.Println("Successful encryption")
// 	decrypted, err := crypto.Decrypt(*encrypted, password)
// 	if err != nil {
// 		fmt.Println("Decryption failed: ", err)

// 	}
// 	fmt.Printf("Original: %s\n", plaintext)
// 	fmt.Printf("Decrypted: %s\n", decrypted)

// }

// func testVault() {
// 	fmt.Println("Testing vault functionality...")
// 	vault := storage.NewVault()
// 	entry := &storage.Entry{
// 		Name:     "Kylo the Doge",
// 		Password: "meowmeowmeow47",
// 		Notes:    "who put this cat here",
// 		URL:      "https://google.com",
// 	}

// 	vault.AddEntry(entry)

// 	pass := "kylo4ever"
// 	err := storage.SaveVault(vault, pass)
// 	if err != nil {
// 		fmt.Println("Error saving to vault: ", err)
// 		return
// 	}

// 	fmt.Println("Vault saved")

// 	loaded, err := storage.LoadVault(pass)
// 	if err != nil {
// 		fmt.Println("Error loading vault: ", err)
// 		return
// 	}

// 	loadedEntry, exists := loaded.GetEntry("Kylo the Doge")
// 	if !exists {
// 		fmt.Println("Error loading entry inside vault: ", exists)
// 		return
// 	}

// 	if loadedEntry.Password == entry.Password {
// 		fmt.Println("Vault test passed")
// 	} else {
// 		fmt.Println("Vault test failed")
// 	}

// }

// func main() {
// 	testCrypto()
// 	testVault()
// }
