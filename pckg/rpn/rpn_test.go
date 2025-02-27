package rpn_test

import (
	"log"
	"testing"

	"github.com/romanSPB15/Calculator_Service/pckg/rpn"
)

// Результат Calc
type CalcResult struct {
	Number float64
	Error  error
}

func TestRPN(t *testing.T) {
	// Канал задач
	c := make(chan rpn.IDTask)
	// Мап задач
	m := make(map[rpn.IDTask]*rpn.Task)
	go func() {
		res, err := rpn.Calc("(2+2/2)+120", rpn.NewConcurrentTaskMap(), c)
		close(c)
		log.Println("Result:", res)
		if err != nil {
			log.Fatalf("TestRPN: Error Calc: %v", err)
		}
		e := 123.0
		if res != e {
			log.Fatalf("TestRPN: Result Calc(%.1F) != %.1F", res, e)
		}
	}()
	for id := range c {
		m[id].Run()
	}
}
