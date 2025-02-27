package application

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/romanSPB15/Calculator_Service/pckg/rpn"
)

// Выражение
type Expression struct {
	Data   string  `json:"-"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

// Выражение с ID
type ExpressionWithID struct {
	ID IDExpression `json:"id"`
	Expression
}

// ID выражения
type IDExpression = uint32

// Выражения
var Expressions = make(map[IDExpression]*Expression)

// Задачи
var Tasks = rpn.NewConcurrentTaskMap()

type AddHandlerResult struct {
	ID uint32 `json:"id"`
}

func AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	var req map[string]string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var code []byte
	_, err = fmt.Fprintf(w, "%s\r\n", code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := uuid.New().ID()
	str, has := req["expression"]
	if has == false {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	e := Expression{str, "Not OK", 0}
	Expressions[id] = &e
	ChanIDTasks := make(chan rpn.IDTask)
	go func() {
		res, err := rpn.Calc(str, Tasks, ChanIDTasks)
		if err != nil {
			e.Status = err.Error()
		} else {
			e.Status = "OK"
			e.Result = res
		}
	}()
	data, err := json.Marshal(AddHandlerResult{id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusCreated)
}

type GetExpressionHandlerResult struct {
	Expression ExpressionWithID `json:"expression"`
}

func GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	strid := r.FormValue("id")
	i, err := strconv.Atoi(strid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	id := IDExpression(i)
	exp, has := Expressions[id]
	if !has {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, err := json.Marshal(GetExpressionHandlerResult{ExpressionWithID{id, *exp}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

type GetExpressionsHandlerResult struct {
	Expressions []ExpressionWithID `json:"expressions"`
}

func GetExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var ExpressionsID []ExpressionWithID
	for id, e := range Expressions {
		ExpressionsID = append(ExpressionsID, ExpressionWithID{id, *e})
	}
	data, err := json.Marshal(GetExpressionsHandlerResult{ExpressionsID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

type GetTaskHandlerResult struct {
	Task rpn.TaskID `json:"task"`
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var tid rpn.TaskID
	for id, t := range *Tasks.Map() {
		tid = rpn.TaskID{Task: *t, ID: id}
		break
	}
	if tid.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(GetTaskHandlerResult{tid})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

type Application struct {
	// Агент
	Agent http.Client
}

func New() *Application {
	return &Application{http.Client{}}
}

var AgentReqestTime = time.Millisecond * 1000

func (a *Application) worker(t rpn.TaskID) {
	res := t.Run()
	fmt.Println(res)
}

// Запуск агента
func (a *Application) runAgent(startServer chan struct{}) error {
	<-startServer
	var res error
	ch := make(chan struct{})
	go func() {
		log.Println("Agent Runned")
		for {
			time.Sleep(AgentReqestTime)
			resp, err := a.Agent.Get("http://localhost/internal/task:8080")
			if err != nil {
				res = err
				ch <- struct{}{}
				return
			}
			if resp.StatusCode == http.StatusNotFound {
				continue
			}
			fmt.Println(resp.Body)
			defer resp.Body.Close()
			go func() {
				var ResultServer GetTaskHandlerResult
				log.Printf("Task %d Runned\r\n", ResultServer.Task.ID)
				b, _ := io.ReadAll(resp.Body)
				json.Unmarshal(b, &ResultServer)
				fmt.Println(ResultServer.Task.Run())
			}()
		}
	}()
	<-ch
	return res
}

// Запуск всей системы
func (a *Application) RunServer() {
	rpn.InitEnv()
	http.HandleFunc("/api/v1/calculate", AddExpressionHandler)
	http.HandleFunc("/api/v1/expressions/:id", GetExpressionHandler)
	http.HandleFunc("/api/v1/expressions", GetExpressionsHandler)
	http.HandleFunc("/api/v1/internal/task", GetTaskHandler)
	startServer := make(chan struct{})
	go func() {
		startServer <- struct{}{}
		log.Println("Orkestrator Runned")
		err := http.ListenAndServe(":8080", nil)
		panic(err)
	}()
	panic(a.runAgent(startServer))
}
