package main

import (
	"log"

	"Taskify/internal/app"
)

func main() {
	// Вся логика, включая инициализацию и запуск, скрыта внутри app.Run()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
