package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 4 {
		panic("usage: gc -g <password>\n       gc -c <password> <hash>")
	}
	op := os.Args[1]
	switch op {
	case "-g":
		password := os.Args[2]
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(hash))
	case "-c":
		password := os.Args[2]
		hash := os.Args[3]
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			fmt.Println("false")
		} else {
			fmt.Println("true")
		}
	default:
		panic("usage: gc -g <password>\n       gc -c <password> <hash>")
	}

}
