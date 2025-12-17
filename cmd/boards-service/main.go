package main

import (
	"Taskify/internal/boardsapp"
	"log"
)

func main() {
	if err := boardsapp.Run(); err != nil {
		log.Fatal(err)
	}
}
