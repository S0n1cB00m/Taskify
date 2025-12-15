package main

import (
	"log"

	"Taskify/internal/app"
)

// @title           Taskify API
// @version         1.0
// @description     This is a sample server for Taskify application.
// @host            localhost:3000
// @BasePath        /api
// @schemes         http
func main() {
	// Вся логика, включая инициализацию и запуск, скрыта внутри app.Run()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
