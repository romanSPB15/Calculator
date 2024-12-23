package application

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/romanSPB15/Calculator/pckg/rpn"
)

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	data := make([]byte, 100)
	defer r.Body.Close()
	n, err := r.Body.Read(data)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	data = data[:n]
	var expstruct struct {
		Expession string
	}
	err = json.Unmarshal(data, &expstruct)

	var code []byte
	result, err := rpn.Calc(expstruct.Expession)
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

func (a *Application) Start_Server(port string) {
	http.HandleFunc("/api/calc", CalcHandler)
	http.ListenAndServe(port, nil)
}
