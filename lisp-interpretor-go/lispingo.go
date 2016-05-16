package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type number float64
type symbol string

//symbols, numbers, expressions, procedures, lists, ... all implement this interface, which enables passing them along in the interpreter
type inter interface{}

// environments
// An environment is a mapping from variable names to their values.
type variable map[symbol]inter

type env struct {
	variable
	outer *env
}

//  init() is always called, regardless if there's main or not,
//  so if you import a package that has an init function, it will be executed.

var mathsenv = env{
	variable{ //aka an incomplete set of compiled-in functions
		"+": func(a ...inter) inter { // variadic functions
			v := a[0].(number)
			for _, i := range a[1:] {
				v += i.(number)
			}
			fmt.Println("func init() : ", v)
			return v
		},
		"-": func(a ...inter) inter {
			v := a[0].(number)
			for _, i := range a[1:] {
				v -= i.(number)
			}
			return v
		},
		"*": func(a ...inter) inter {
			v := a[0].(number)
			for _, i := range a[1:] {
				v *= i.(number)
			}
			return v
		},
		"/": func(a ...inter) inter {
			v := a[0].(number)
			for _, i := range a[1:] {
				v /= i.(number)
			}
			return v
		},
		">": func(a ...inter) inter {
			return a[0].(number) > a[1].(number)
		},
		">=": func(a ...inter) inter {
			return a[0].(number) >= a[1].(number)
		},
		"<": func(a ...inter) inter {
			return a[0].(number) < a[1].(number)
		},
		"<=": func(a ...inter) inter {
			return a[0].(number) <= a[1].(number)
		},
		"equal?": func(a ...inter) inter {
			return reflect.DeepEqual(a[0], a[1])
		},
		"length": func(a ...inter) inter {
			return number(len(a[0].([]inter)))
		},
		"append": func(a ...inter) inter {
			result := make([]inter, 0)
			for _, i := range a {
				result = append(result, i.([]inter)[:]...)
			}
			return result
		},
		"null?": func(a ...inter) inter {
			return len(a[0].([]inter)) == 0
		},
		"cons": func(a ...inter) inter {
			switch car := a[0]; cdr := a[1].(type) {
			case []inter:
				return append([]inter{car}, cdr...)
			default:
				return []inter{car, cdr}
			}
		},
		"car": func(a ...inter) inter { //
			return a[0].([]inter)[0]
		},
		"cdr": func(a ...inter) inter {
			return a[0].([]inter)[1:]
		},
		//"list": eval(readFrom(
		//	"(lambda z z)"),
		//	&globalenv),
	},
	nil}

func scan_expression() string {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">>> ")
	text, _ := reader.ReadString('\n')
	return text
}

// Lexical Analysis
func create_tokens(s string) []string {
	if strings.Contains(s, "(") {
		s = strings.Replace(s, "(", " ( ", -1)
	}
	if strings.Contains(s, ")") {
		s = strings.Replace(s, ")", " ) ", -1)
	}
	if strings.Contains(s, "\n") {
		s = strings.Replace(s, "\n", "", -1)
	}

	split_s := strings.Fields(s)

	return split_s
}

//Syntactic Analysis
func readFrom(tokens []string) inter {

	// take first element from tokens
	token := tokens[0]
	tokens = append(tokens[:0], tokens[1:]...)
	switch token {
	case "(": // list begins
		L := make([]inter, 0)
		for tokens[0] != ")" {
			i := readFrom(tokens)
			L = append(L, i)
		}
		tokens = append(tokens[:0], tokens[1:]...)
		return L
	default: // atom
		if f, err := strconv.ParseFloat(token, 64); err == nil {
			return number(f)
		} else {
			return symbol(token)
		}
	}
}

// The function eval takes two arguments: an expression, x, that we want to evaluate,
// and an environment, env, in which to evaluate it

func eval(expr inter, en *env) (res inter) {

	defer getback() // in case of any runtime exception

	switch t := expr.(type) {
	case number:
		res = t
		fmt.Println("number : ", res)
	case symbol: // there is a list of symbols (operators/keywords etc) defined in init()
		res = en.variable[t]
		fmt.Println(res)
	case []inter:
		fmt.Println("array")
		for i := 0; i < len(t); i++ {
			switch p := t[i].(type) {
			case []inter:
				//return eval(p)

			case symbol:

				// fmt.Println("val returned from init : ", string(res))
				res = en.variable[p]
				fmt.Println(res)
				//y := eval(p, en)
				//fmt.Println(y.(symbol))
				//if t[1].(type) == number && t[2].(type) == number {
				v1 := t[1].(number)
				v2 := t[2].(number)
				if p == "+" {
					fmt.Println(v1 + v2)
					res = v1 + v2
				} else if p == "-" {
					fmt.Println(v1 - v2)
					res = v1 - v2
				} else if x == "*" {
					fmt.Println(v1 * v2)
					res = v1 * v2
				} else if x == "/" {
					fmt.Println(v1 / v2)
					res = v1 / v2
				} else if x == ">" {
					fmt.Println(v1 > v2)
					res = v1 > v2
				} else if x == "<" {
					fmt.Println(v1 < v2)
					res = v1 < v2
				} else if x == ">=" {
					fmt.Println(v1 >= v2)
					res = v1 >= v2
				} else if x == "<=" {
					fmt.Println(v1 <= v2)
					res = v1 <= v2
				} else if x == "==" {
					fmt.Println(v1 == v2)
					res = v1 == v2
				}
				//}
			case number:
				//v[i] = t[i].(number)
				//fmt.Println(t[i].(number))
			default:

			}
		}
	}
	return res
	//}
}

func getback() {
	r := recover()
	if r != nil {
		fmt.Println("Undefined format : succesfully recovered ")
	}
}

func main() {

	input_string := scan_expression()
	//fmt.Println(input_string)

	split_string := create_tokens(input_string)
	fmt.Printf("%s\n", split_string)

	var expression inter
	expression = readFrom(split_string)

	fmt.Println(expression)

	result := eval(expression, &mathsenv)
	fmt.Println(result)
}
