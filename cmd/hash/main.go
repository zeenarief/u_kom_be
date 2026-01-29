package main

import (
	"fmt"
	"log"
	"u_kom_be/internal/utils"
)

func main() {
	password := "admin123" // password default admin
	hash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Hashed password:", hash)
}
