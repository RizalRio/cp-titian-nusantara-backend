// cmd/generate-password/main.go
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
    password := "Admin123!" // Password yang mau di-hash
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    fmt.Println("Hash untuk SQL:", string(hash))
}