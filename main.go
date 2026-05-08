package main

import (
	"company-site/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
