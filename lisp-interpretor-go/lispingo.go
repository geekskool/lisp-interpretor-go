package main

import (
	"bufio"
	"fmt"
	"os"
	//"reflect"
	"strconv"
	"strings"
)

type number float64
type symbol string
type fun func(a ...float64) float64

/*
func add(a ...float64) float64 {
	var v float64
	for i := 0; i < len(a); i++ {
		v = v + a[i]
	}
	return v
}
*/

//symbols, numbers, expressions, procedures, lists, ... all implement this interface, which enables passing them along in the interpreter
type inter interface{}

// environments
// An environment is a mapping from variable names to their values.
type variable map[symbol]inter

var identifiers = make(map[inter]inter)

type env struct {
	variable
	other *env
}

//  init() is always called, regardless if there's main or not,
//  so if you import a package that has an init function, it will be executed.

var mathsenv = env{
	variable{ //aka an incomplete set of compiled-in functions
		"+": func(a ...inter) inter {
			v := a[0].(number)
			for _, i := range a[1:] {
				v += i.(number)
			}
			return v
		},
	},
	nil}

/*
var mathsenv = env{
	variable{ //aka an incomplete set of compiled-in functions
		"+": func(a ...number) number { // variadic functions
			v := a[0].(number)
			for _, i := range a[1:] {
				v += i.(number)
			}
			fmt.Println("func init() : ", v)
			return v
		},
	},
	nil}
*/

func scan_expression() string {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">>> ")
	text, _ := reader.ReadString('\n')
	return text
}

// Lexical Analysis
func create_tokens(s string) []string {

	s = strings.Replace(s, "(", " ( ", -1)
	s = strings.Replace(s, ")", " ) ", -1)
	s = strings.Replace(s, "\n", "", -1)

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

func add(a []number) inter {
	sum := a[0]
	for i := 1; i < len(a); i++ {
		sum += a[i]
	}
	return sum
}

func sub(a []number) inter {
	diff := a[0]
	for i := 1; i < len(a); i++ {
		diff = diff - a[i]
	}
	return diff
}

func mul(a []number) inter {
	prod := a[0]

	for i := 1; i < len(a); i++ {
		prod = prod * a[i]
	}
	return prod
}

func div(a []number) inter {
	num := a[0]

	for i := 1; i < len(a); i++ {
		num = num / a[i]
	}
	return num
}

func typeconverter(arr1 inter) []number {
	var result = make([]number, 0) // problem defining len of arr
	switch ai := arr1.(type) {
	case []inter:
		for i := 1; i < len(ai); i++ {
			switch x := ai[i].(type) {
			case number:
				result = append(result, x)
				//result[i] = x
			case []inter:
				n := mathsop(ai[i]) // returns inter (result of operator function (+))
				result = append(result, n.(number))
				// save the resultant answer evaluated to the result array and typecast

			}
		}
	}
	fmt.Println(result)
	return result
}

func mathsop(arr inter) inter {
	var ans inter
	switch ai := arr.(type) {
	case []inter:
		switch op := ai[0].(type) {
		case symbol:
			switch op {
			case "+":
				tc := typeconverter(ai)
				ans = add(tc)
			case "-":
				tc := typeconverter(arr)
				ans = sub(tc)
			case "*":
				tc := typeconverter(arr)
				ans = mul(tc)
			case "/":
				tc := typeconverter(arr)
				ans = div(tc)

			}
		}
	}
	return ans
}

func stringify(v inter) string {
	switch v := v.(type) {
	case []inter:
		l := make([]string, len(v))
		for i, x := range v {
			l[i] = stringify(x)
		}
		return "(" + strings.Join(l, " ") + ")"
	default:
		return fmt.Sprint(v)
	}
}

/*
func verify(cond inter) bool {
	switch s := cond.(type) {
	case []inter:

	}
	return true
}
*/
// The function eval takes two arguments: an expression, x, that we want to evaluate,
// and an environment, env, in which to evaluate it

func eval(expr inter, en env) inter {
	var ans inter
	defer getback() // in case of any runtime exception

	switch s := expr.(type) {
	case symbol:
	case number:
	case []inter:
		switch tt := s[0].(type) {
		case symbol:
			switch tt {
			case "+", "-", "*", "/", ">", ">=", "<", "<=", "==", "%":
				ans := mathsop(s)
				fmt.Println(ans)
			case "quote":
				fmt.Println(stringify(s))
			case "if":
				//				verify(s)
			}
		}
	}

	return ans
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

	//var expression inter
	expression := readFrom(split_string)

	fmt.Println(expression)

	result := eval(expression, mathsenv)
	fmt.Println(result)
}
