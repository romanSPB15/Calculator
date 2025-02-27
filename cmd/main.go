package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/romanSPB15/Calculator_Service/internal/application"
)

func main() {
	go func() {
		a := application.New()
		a.RunServer()
	}()
	TestClient := http.Client{}
	resp, err := TestClient.Post("http://localhost/api/v1/calculate", "application/json", strings.NewReader(`{"expression": "2+2"}`))
	fmt.Println(resp, err)
	for {
	}
}
