package main

import (
	"Taskify/internal/usersapp"
	"log"
)

func main() {
	if err := usersapp.Run(); err != nil {
		log.Fatal(err)
	}
}
