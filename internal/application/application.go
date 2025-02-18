package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/romanSPB15/Calculator_Service/pckg/rpn"
)

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		w.WriteHeader(500)
		return
	}
	var req map[string]string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	result, err := rpn.Calc(req["expression"])
	var code []byte
	if err == nil {
		var structres struct {
			Result float64 `json:"result"`
		}
		structres.Result = result
		code, _ = json.Marshal(structres)
		w.WriteHeader(200)
	} else if err == rpn.Errorexp {
		var structres struct {
			Error string `json:"error"`
		}
		structres.Error = err.Error()
		code, _ = json.Marshal(structres)
		w.WriteHeader(422)
	} else {
		w.WriteHeader(500)
		return
	}
	_, err = fmt.Fprintf(w, "%s\r\n", code)
	if err != nil {
		w.WriteHeader(500)
	}
}

type Application struct{}

func NewApplication() *Application {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return &Application{}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/calc", CalcHandler)
	return http.ListenAndServe(":80", nil)
}
