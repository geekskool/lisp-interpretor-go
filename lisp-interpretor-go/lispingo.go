package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type number float64
type symbol string
type boolean bool

//symbols, numbers, expressions, procedures, lists, ... all implement this interface, which enables passing them along in the interpreter
type inter interface{}

// environments
// An environment is a mapping from variable names to their values.

var identifiers = make(map[inter]inter) // anything can be mapped to anything

type procs struct { //  handling procedures
	arguments []symbol //  saving arguments
	body      inter    //  body of function
}

var argbod procs // global variable for saving procedures

// take input from user
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

// arithematic operations
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

func gt(a []number) inter {
	num1 := a[0]
	num2 := a[1]

	return num1 > num2
}

func gteq(a []number) inter {
	num1 := a[0]
	num2 := a[1]

	return num1 >= num2
}

func lt(a []number) inter {
	num1 := a[0]
	num2 := a[1]

	return num1 < num2
}

func lteq(a []number) inter {
	num1 := a[0]
	num2 := a[1]

	return num1 <= num2
}

func equals(a []number) inter {
	num1 := a[0]
	num2 := a[1]

	return num1 == num2
}

// converts type interface to number(float64) : used to get operands from expression
func typeconverter(arr1 inter) []number {
	var result = make([]number, 0) // problem defining len of arr
	switch ai := arr1.(type) {
	case []inter:
		for i := 1; i < len(ai); i++ {
			switch x := ai[i].(type) {
			case number:
				result = append(result, x)
			case []inter:
				n := mathsop(ai[i]) // returns inter (result of operator function (+))
				result = append(result, n.(number))
				// save the resultant answer evaluated to the result array and typecast
			case symbol:
				n := identifiers[x]
				result = append(result, n.(number))
			}
		}
	}
	fmt.Println(result)
	return result
}

// performs maths(arithematics) and condition checking operations
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
			case ">":
				tc := typeconverter(arr)
				ans = gt(tc)
			case ">=":
				tc := typeconverter(arr)
				ans = gteq(tc)
			case "<":
				tc := typeconverter(arr)
				ans = lt(tc)
			case "<=":
				tc := typeconverter(arr)
				ans = lteq(tc)
			case "==":
				tc := typeconverter(arr)
				ans = equals(tc)

			}
		}
	}
	return ans
}

// converts expression into string used with "quote"
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

// println like function
// syntax : ( print "string" anything )
func doprint(arr inter) {
	switch s := arr.(type) {
	case []inter:
		switch i := s[1].(type) {
		case number:
			fmt.Println(i)
		case symbol:
			str := stringify(i)
			if strings.HasPrefix(str, "\"") {
				fmt.Println(str)
			} else {
				fmt.Println(identifiers[i])
			}
		}
	}
}

// define calls save to save the identifiers and their values
// syntax: ( define x 20)
// eg2:    ( define area lambda (x y) ( * x y) )  => if lambda func then call lambdafun()
func save(arr inter) inter {
	switch p := arr.(type) {
	case []inter:
		if _, ok := p[2].(symbol); ok {
			if p[2].(symbol) == "lambda" {
				fmt.Println("calling lambda function")
				lambdafun(p[2:])
			}
		} else {
			identifiers[p[1]] = p[2]
		}
	}
	return identifiers
}

// used to check if else statements
// syntax: ( if nil (print "hello") (print "bye"))   op=> bye
// eg2: (if (> 4 2) (print "greater") (print "smaller"))   op=> greater
func condition(arr inter) {
	switch s := arr.(type) {
	case []inter:
		switch tt := s[1].(type) {
		case symbol:
			if tt == "nil" {
				eval(s[3])
			} else {
				eval(s[2])
			}
		case []inter:
			x := eval(tt).(bool)
			if x {
				eval(s[2])
			} else {
				eval(s[3])
			}
		}
	}
}

// used with "set!" (manipulates value of a variable)
// syntax: ( set! x (* a b 8))  >> assume a = 10 b = 20 , then x = 1600
func assign(arr inter) {
	switch s := arr.(type) {
	case []inter:
		switch x := s[2].(type) {
		case []inter:
			identifiers[s[1]] = eval(x)
			fmt.Println(identifiers[s[1]])
		}
	}
}

// declares lambda function and maps environments function name, arguments and function body
// syntax:    ( define area lambda (x y) ( * x y) )  => if lambda func then call lambdafun()
func lambdafun(arr inter) {
	var v number
	fmt.Println("lambda : ", arr)

	switch s := arr.(type) {
	case []inter:
		switch a := s[1].(type) {
		case []inter:
			for i := 0; i < len(a); i++ {
				argbod.arguments = append(argbod.arguments, a[i].(symbol))
				identifiers[a[i]] = v //arguments
			}
		}
		argbod.body = s[2] // expression or body of function
	}
	fmt.Println("arguments saved in argbod ", argbod.arguments)
}

// calls function after defining
// suppose area function is defined in previous case, then,
// syntax: ( area 10 20 )   => op - 200
func funcall(arr inter) {
	//var temp = make([]inter, 0)
	switch a := arr.(type) {
	case []inter:
		name := a[0].(symbol)
		if _, ok := identifiers[name]; ok { // if ok -> true ; key is present // function name
			fmt.Println("function has been declared \n")
		}

		for i := 0; i < len(a)-1; i++ {
			identifiers[argbod.arguments[i]] = a[i+1].(number)
		}

		identifiers[name] = eval(argbod.body)
	}
}

// creates a list of interfaces
// syntax : ( list name_of_list e1 e2 e3 ... )
func create_list(arr inter) {
	var temp []inter
	var name symbol
	switch s := arr.(type) {
	case []inter:
		name = s[1].(symbol)
		for j := 2; j < len(s); j++ {
			temp = append(temp, s[j])
		}
		identifiers[name] = temp
	}
}

// appends (any element) interfaces to the list
// syntax : ( append name_of_list_already_declared  e1 e2 e3... )
func listconcat(arr inter) {
	var name symbol
	var temp []inter
	fmt.Println(arr)
	switch s := arr.(type) {
	case []inter:
		name = s[0].(symbol)
		temp = identifiers[name].([]inter)
		for i := 1; i < len(s); i++ {
			temp = append(temp, s[i])
		}
	}
	identifiers[name] = temp
}

// The function eval takes one argument: an expression, x, that we want to evaluate,
func eval(expr inter) inter {
	var ans inter
	defer getback() // in case of any runtime exception

	switch s := expr.(type) {
	case symbol:
	case number:
	case []inter:
		switch tt := s[0].(type) {
		case symbol:
			switch tt {
			case "+", "-", "*", "/", ">", ">=", "<", "<=", "==":
				ans = mathsop(s)
				fmt.Println(ans)
			case "quote":
				fmt.Println(stringify(s))
			case "define":
				fmt.Println(save(s))
			case "print":
				doprint(s)
			case "if":
				condition(s)
			case "set!": // Evaluate exp and assign that value to var
				assign(s)
			case "lambda":
				lambdafun(s)
			case "list":
				create_list(s)
			case "append":
				listconcat(s[1:])
			case "car":
				array := identifiers[s[1].(symbol)].([]inter)
				fmt.Println(array[0])
			case "cdr":
				array := identifiers[s[1].(symbol)].([]inter)
				fmt.Println(array[1:])
			default: // function call
				funcall(s)
			}
		}
	}

	return ans
}

// to recover from panic ( run time exception)
func getback() {
	r := recover()
	if r != nil {
		fmt.Println("Undefined format : succesfully recovered ")
	}
}

func main() {

	for {
		input_string := scan_expression()
		//fmt.Println(input_string)
		if strings.Compare(input_string, "(quit)\n") == 0 {
			fmt.Println("exiting....")
			break
		}

		split_string := create_tokens(input_string)

		expression := readFrom(split_string)

		eval(expression)

	}
}
