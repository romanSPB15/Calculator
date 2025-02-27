package application_test

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"net/http/httptest"

	app "github.com/romanSPB15/Calculator_Service/internal/application"
	"github.com/romanSPB15/Calculator_Service/pckg/rpn"
)

func TestAddHandlerAndGetExpressionsHandler(t *testing.T) {
	rpn.InitEnv()
	expressions := []string{"1+1", "2+2", "3+3", "4+4", "5+5"}
	for _, exp := range expressions {
		r := strings.NewReader(fmt.Sprintf("{\"expression\": \"%s\"}", exp))
		req := httptest.NewRequest("POST", "http://localhost/api/v1/calculate", r)
		w := httptest.NewRecorder()
		app.AddExpressionHandler(w, req)
		resp := w.Result()
		resp.Body.Close()
	}
	url := `http://localhost/api/v1/expressions`
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	app.GetExpressionsHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	log.Println(w.Code, w.Body)
}

func TestAddHandlerAndGetExpressionHandler(t *testing.T) {
	rpn.InitEnv()
	r := strings.NewReader("{\"expression\": \"(2+3-1)/2*5\"}")
	req := httptest.NewRequest("POST", "http://localhost/api/v1/calculate", r)
	w := httptest.NewRecorder()
	app.AddExpressionHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	var addres app.AddHandlerResult
	err := d.Decode(&addres)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf(`http://localhost/api/v1/expressions?id=%d`, addres.ID)
	req = httptest.NewRequest("GET", url, nil)
	w = httptest.NewRecorder()
	app.GetExpressionHandler(w, req)
	resp = w.Result()
	defer resp.Body.Close()
	fmt.Println(w.Code, w.Body)
}

func TestAddHandlerAndGetTaskHandler(t *testing.T) {
	rpn.InitEnv()
	r := strings.NewReader("{\"expression\": \"2+8\"}")
	req := httptest.NewRequest("POST", "http://localhost/api/v1/calculate", r)
	w := httptest.NewRecorder()
	app.AddExpressionHandler(w, req)
	url := `http://localhost/api/v1/internal/task`
	time.Sleep(100 * time.Millisecond)
	req = httptest.NewRequest("GET", url, nil)
	w = httptest.NewRecorder()
	app.GetTaskHandler(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	fmt.Println(w.Code, w.Body)
}
