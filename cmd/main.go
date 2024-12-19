package main

import app "github.com/romanSPB15/Calculator/internal/application"

const port = ":80"

func main() {
	app := app.NewApplication()
	app.Start_Server(port)
}
