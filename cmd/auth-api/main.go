package main

import (
	"log"

	"github.com/dev-gopi/authhub/internal/bootstrap"
)

func main() {
	if err := bootstrap.RunAuthAPI(); err != nil {
		log.Fatal(err)
	}
}
