package rpn

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var (
	TIME_ADDITION_MS        int
	TIME_SUBTRACTION_MS     int
	TIME_MULTIPLICATIONS_MS int
	TIME_DIVISIONS_MS       int
	COMPUTING_POWER         int
)

// Считывание переменной среды в виде числа
func getIntEnv(key string) int {
	str, has := os.LookupEnv(key)
	if !has {
		log.Panicf("System has not %s", key)
	}
	res, err := strconv.Atoi(str)
	if err != nil {
		log.Panicf("Env %s is not int", key)
	}
	return res
}

// Иницилизация переменных Go из файла .env
func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	TIME_ADDITION_MS = getIntEnv("TIME_ADDITION_MS")
	TIME_SUBTRACTION_MS = getIntEnv("TIME_SUBTRACTION_MS")
	TIME_MULTIPLICATIONS_MS = getIntEnv("TIME_MULTIPLICATIONS_MS")
	TIME_DIVISIONS_MS = getIntEnv("TIME_DIVISIONS_MS")
	COMPUTING_POWER = getIntEnv("COMPUTING_POWER")
}

type (
	TaskArg1Type   = float64
	TaskArg2Type   = float64
	TaskResultType = float64
)

type ExpressionResultType = float64

// Задача
type Task struct {
	Arg1          TaskArg1Type   `json:"arg1"`
	Arg2          TaskArg2Type   `json:"arg2"`
	Operation     string         `json:"operation"`
	OperationTime int            `json:"operation_time"`
	Status        string         `json:"-"`
	Result        TaskResultType `json:"-"`
}

// Мап задач
type TaskMap = map[IDTask]*Task

// Конкурентный мап задач
type ConcurrentTaskMap struct {
	m  TaskMap
	mx sync.Mutex
}

// Новый конкурентный мап задач
func NewConcurrentTaskMap() *ConcurrentTaskMap {
	return &ConcurrentTaskMap{make(map[IDTask]*Task), sync.Mutex{}}
}

// Получение ссылки на задачу.
// Если такого значения нет, то оно добавляется само.
func (cm *ConcurrentTaskMap) Get(id IDTask) *Task {
	cm.mx.Lock()
	res, ok := cm.m[id]
	if ok == false {
		t := &Task{}
		cm.m[id] = t
		cm.mx.Unlock()
		return t
	}
	cm.mx.Unlock()
	return res
}

// Получение простого мапа
func (cm *ConcurrentTaskMap) Map() *map[IDTask]*Task {
	return &cm.m
}

// Задача с ID
type TaskID struct {
	ID IDTask `json:"id"`
	Task
}

func (t *Task) Run() (res float64) {
	switch t.Operation {
	case "+":
		res = t.Arg1 + t.Arg2
	case "-":
		res = t.Arg1 - t.Arg2
	case "*":
		res = t.Arg1 * t.Arg2
	case "/":
		res = t.Arg1 / t.Arg2
	}
	return
}

func convertString(str string) float64 {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return res
}

func isSign(value rune) bool {
	return value == '+' || value == '-' || value == '*' || value == '/'
}

type IDTask = uint32

var Errorexp = errors.New("Expression is not valid")
var Errordel = errors.New("/0!")

func Calc(expression string, tasks *ConcurrentTaskMap, newID chan IDTask) (res ExpressionResultType, err0 error) {
	if len(expression) < 3 {
		return 0, Errorexp
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	b := ""
	c := rune(0)
	resflag := false
	isc := -1
	scc := 0
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	if isSign(rune(expression[0])) || isSign(rune(expression[len(expression)-1])) {
		return 0, Errorexp
	}
	if strings.Contains(expression, "(") || strings.Contains(expression, ")") {
		for i := 0; i < len(expression); i++ {
			value := expression[i]
			if value == '(' {
				if scc == 0 {
					isc = i
				}
				scc++
			}
			if value == ')' {
				scc--
				if scc == 0 {
					exp := expression[isc+1 : i]
					calc, err := Calc(exp, tasks, newID)
					if err != nil {
						return 0, err
					}
					calcstr := strconv.FormatFloat(calc, 'f', 0, 64)
					expression = strings.Replace(expression, expression[isc:i+1], calcstr, 1) // Меняем скобки на результат выражения в них

					i -= len(exp)
					isc = -1
				}
			}
		}
	}
	if isc != -1 {
		return 0, Errorexp
	}
	priority := strings.ContainsRune(expression, '*') || strings.ContainsRune(expression, '/')
	notpriority := strings.ContainsRune(expression, '+') || strings.ContainsRune(expression, '-')
	if priority && notpriority {
		for i := 1; i < len(expression); i++ {
			value := rune(expression[i])
			///////////////////////////////////////////////////////////////////////////////////////////////////////////////
			//Умножение и деление
			if value == '*' || value == '/' {
				var imin int = i - 1
				if imin != 0 {
					for imin >= 0 {
						if imin >= 0 {
							if isSign(rune(expression[imin])) {
								break
							}
						}
						imin--
					}
					imin++
				}
				imax := i + 1
				if imax == len(expression) {
					imax--
				} else {
					for !isSign(rune(expression[imax])) && imax < len(expression)-1 {
						imax++
					}
				}
				if imax == len(expression)-1 {
					imax++
				}
				exp := expression[imin:imax]
				calc, err := Calc(exp, tasks, newID)
				if err != nil {
					return 0, err
				}
				calcstr := strconv.FormatFloat(calc, 'f', 0, 64)
				expression = strings.Replace(expression, expression[imin:imax], calcstr, 1) // Меняем скобки на результат выражения в них
				i -= len(exp) - 1
			}
			if value == '+' || value == '-' || value == '*' || value == '/' {
				c = value
			}
		}
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	for _, value := range expression + "s" {
		switch {
		case value == ' ':
			continue
		case value > 47 && value < 58 || value == '.': // Если это цифра
			b += string(value)
		case isSign(value) || value == 's': // Если это знак
			if resflag {
				switch c {
				case '+':
					uuid := uuid.New()
					id := uuid.ID()
					t := Task{
						Arg1:          res,
						Arg2:          convertString(b),
						Operation:     "+",
						OperationTime: TIME_ADDITION_MS,
					}
					log.Println("New Task With ID", id)
					*tasks.Get(id) = t // Записываем задачу
					status := &t.Status
					newID <- id
					for *status != "OK" {
					}
					log.Printf("m[%d].Status = \"OK\"", id)
					res = t.Result
				case '-':
					uuid := uuid.New()
					id := uuid.ID()
					t := Task{
						Arg1:          res,
						Arg2:          convertString(b),
						Operation:     "-",
						OperationTime: TIME_SUBTRACTION_MS,
					}
					log.Println("rpn.Calc: Create New Task With ID", id)
					*tasks.Get(id) = t // Записываем задачу
					status := &t.Status
					newID <- id
					for *status != "OK" {
					}
					log.Printf("m[%d].Status = \"OK\"", id)
					res = t.Result
				case '*':
					uuid := uuid.New()
					id := uuid.ID()
					t := Task{
						Arg1:          res,
						Arg2:          convertString(b),
						Operation:     "*",
						OperationTime: TIME_MULTIPLICATIONS_MS,
					}
					log.Println("rpn.Calc: Create New Task With ID", id)
					*tasks.Get(id) = t // Записываем задачу
					status := &t.Status
					newID <- id
					for *status != "OK" {
					}
					log.Printf("m[%d].Status = \"OK\"", id)
					res = t.Result
				case '/':
					uuid := uuid.New()
					id := uuid.ID()
					t := Task{
						Arg1:          res,
						Arg2:          convertString(b),
						Operation:     "/",
						OperationTime: TIME_DIVISIONS_MS,
					}
					log.Println("rpn.Calc: Create New Task With ID", id)
					*tasks.Get(id) = t // Записываем задачу
					status := &t.Status
					newID <- id
					for *status != "OK" {
					}
					log.Printf("m[%d].Status = \"OK\"", id)
					res = t.Result
				}
			} else {
				resflag = true
				res = convertString(b)
			}
			b = ""
			c = value

			/////////////////////////////////////////////////////////////////////////////////////////////
		case value == 's':
		default:
			return 0, Errorexp
		}
	}
	return res, nil
}
