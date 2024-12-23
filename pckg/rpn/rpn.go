package rpn

import (
	"errors"
	"strconv"
	"strings"
)

func stringToFloat64(str string) (res float64) {
	degree := float64(1)
	if strings.ContainsRune(str, '.') {
		drob := 0.00
		i := len(str) - 1
		for ; str[i] != '.'; i-- {
			drob += float64(9-int('9'-str[i])) * degree
			degree *= 10
		}
		degree = 1
		for ; i > 0; i-- {
			res += float64(9-int('9'-str[i-1])) * degree
			degree *= 10
		}
		for drob > 1 {
			drob /= 10
		}
		res += drob
	} else {
		for i := len(str); i > 0; i-- {
			res += float64(9-int('9'-str[i-1])) * degree
			degree *= 10
		}
	}
	return res
}

func isSign(value rune) bool {
	return value == '+' || value == '-' || value == '*' || value == '/'
}

var Errorexp = errors.New("Expression is not valid")
var Errordel = errors.New("/0!")

func Calc(expression string) (res float64, err0 error) {
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
					calc, err := Calc(exp)
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
					for !isSign(rune(expression[imin])) && imin > 0 {
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
				calc, err := Calc(exp)
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
					res += stringToFloat64(b)
				case '-':
					res -= stringToFloat64(b)
				case '*':
					res *= stringToFloat64(b)
				case '/':
					if b == "0" {
						return 0, Errordel
					}
					res /= stringToFloat64(b)
				}
			} else {
				resflag = true
				res = stringToFloat64(b)
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
