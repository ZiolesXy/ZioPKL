package main

import (
	"log"
	"voca-store/internal/app"
)

func main() {
	a := app.New()
	if err := a.Run(); err != nil {
		log.Fatal("SERVER ERROR:", err)
	}
}