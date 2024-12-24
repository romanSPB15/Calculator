package main

import (
	"log"

	app "github.com/romanSPB15/Calculator/internal/application"
)

const port = ":80"

func main() {
	err := app.NewApplication().RunServer()
	if err != nil {
		log.Fatal(err)
	}
}
