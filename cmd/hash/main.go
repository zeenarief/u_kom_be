package main

import (
	"belajar-golang/internal/utils"
	"fmt"
	"log"
)

func main() {
	password := "admin123" // password default admin
	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hashed password:", hash)
}
