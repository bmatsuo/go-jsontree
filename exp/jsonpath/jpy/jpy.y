%{

// This tag will end up in the generated y.go, so that forgetting
// 'make clean' does not fail the next build.

// +build ignore

// jpy.y
// example of a Go yacc program
// usage is
//  go tool yacc -p "jpy_" jpy.y (produces y.go)
//  go build -o jpy y.go
//  ./units $GOROOT/src/cmd/yacc/units.txt
//  you have: c
//  you want: furlongs/fortnight
//      * 1.8026178e+12
//      / 5.5474878e-13
//  you have:

package jpy

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
    "unicode/utf8"
)

const (
    Ndim = 15  // number of dimensions
    Maxe = 695 // log of largest number
)

type Node struct {
    vval float64
    dim  [Ndim]int8
}

type Var struct {
    name string
    node Node
}

var fi *bufio.Reader // input
var fund [Ndim]*Var  // names of fundamental units
var line string      // current input line
var lineno int       // current input line number
var linep int        // index to next rune in unput
var nerrors int      // error count
var one Node         // constant one
var peekrune rune    // backup runt from input
var retnode1 Node
var retnode2 Node
var retnode Node
var sym string
%}

%union {
    node Node
    vvar *Var
    numb int
    vval float64
}

%type   <node>  prog

%token  <vval>  VÄL // dieresis to test UTF-8
%token  <vvar>  VAR
%token  <numb>  _SUP // tests leading underscore in token name
%%
prog:
    ':' VAR
    {
        fmt.Println($2.name)
    }
%%
func Poop() string {
    return "plop"
}
