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

func Calc(expression string) (float64, error) {
	if len(expression) < 3 {
		return 0, Errorexp
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	var res float64
	var b string
	var c rune = 0
	var resflag bool = false
	var isc int
	var countc int = 0
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	for _, value := range expression {
		if isSign(value) {
			countc++
		}
	}
	//////////////////////////////////////////////////////////////////////////////////////////////////////
	if isSign(rune(expression[0])) || isSign(rune(expression[len(expression)-1])) {
		return 0, Errorexp
	}
	for i, value := range expression {
		if value == '(' {
			isc = i
		}
		if value == ')' {
			calc, err := Calc(expression[isc+1 : i])
			if err != nil {
				return 0, err
			}
			calcstr := strconv.FormatFloat(calc, 'f', 0, 64)
			i2 := i
			i -= len(expression[isc:i+1]) - len(calcstr)
			expression = strings.Replace(expression, expression[isc:i2+1], calcstr, 1) // Меняем скобки на результат выражения в них
		}
	}
	if countc > 1 {
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
				var imax int = i + 1
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
				calc, err := Calc(expression[imin:imax])
				if err != nil {
					return 0, err
				}
				calcstr := strconv.FormatFloat(calc, 'f', 0, 64)
				i -= len(expression[isc:i+1]) - len(calcstr) - 1
				expression = strings.Replace(expression, expression[imin:imax], calcstr, 1) // Меняем скобки на результат выражения в них
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
			b = strings.ReplaceAll(b, b, "")
			c = value

			/////////////////////////////////////////////////////////////////////////////////////////////
		case value == 's':
		default:
			return 0, Errorexp
		}
	}
	return res, nil
}
