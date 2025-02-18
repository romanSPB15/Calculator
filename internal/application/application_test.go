package application_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	app "github.com/romanSPB15/Calculator_Service/internal/application"
)

func TestCalcHandler(t *testing.T) {
	r := strings.NewReader("{\"expression\": \"(2+3-1)/2*5\"}")
	req, err := http.NewRequest("POST", "/api/calc", r)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	app.CalcHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("server works incorrect")
	}
	if string(body) != "{\"result\":10}\r\n" {
		t.Fatal("incorrect body", string(body), "expect {\"result\":10}")
	}
	if code := res.StatusCode; code != 200 {
		t.Fatal("incorrect code", code, "expect 200")
	}
	r = strings.NewReader("{\"expression\": \"(2+3-;.1/w*a\"}")
	req, err = http.NewRequest("POST", "/api/calc", r)
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	app.CalcHandler(w, req)
	res = w.Result()
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("server works incorrect")
	}
	if string(body) != "{\"error\":\"Expression is not valid\"}\r\n" {
		t.Fatal("incorrect answer", string(body), "expect {\"result\": \"10\"}")
	}
	if code := res.StatusCode; code != 422 {
		t.Fatal("incorrect code", code, "expect 422")
	}
	r = strings.NewReader("{\"expression\": \"(2+3-;.1/w*a\"}")
	req, err = http.NewRequest("POST", "/api/calc", r)
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	app.CalcHandler(w, req)
	res = w.Result()
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("server works incorrect")
	}
	if string(body) != "{\"error\":\"Expression is not valid\"}\r\n" {
		t.Fatal("incorrect answer", string(body), "expect {\"error\":\"Expression is not valid\"}\r\n")
	}
	if code := res.StatusCode; code != 422 {
		t.Fatal("incorrect code", code, "expect 422")
	}
	r = strings.NewReader("{\"expression\": \"10+10/0\"}")
	req, err = http.NewRequest("POST", "/api/calc", r)
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	app.CalcHandler(w, req)
	res = w.Result()
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	if err != nil {
		t.Fatal("server works incorrect")
	}
	if code := res.StatusCode; code != 500 {
		t.Fatal("incorrect code", code, "expect 500")
	}
}
