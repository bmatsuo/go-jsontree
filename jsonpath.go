// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// jsonpath.go [created: Mon, 10 Jun 2013]

// Package jsonpath does ....
package jsonpath

import (
	"github.com/bmatsuo/go-simplejson"
)

func Lookup(js *simplejson.Json, path ...Selector) []*simplejson.Json {
	var selected []*simplejson.Json
	jschan := make(chan *simplejson.Json, 2)
	go Chain(path...)(jschan, js)
	for js := range jschan {
		if js == nil {
			break
		}
		selected = append(selected, js)
	}
	return selected
}

// Selectors MUST send nil on the channel when there are no more elements.
type Selector func(chan<- *simplejson.Json, *simplejson.Json)

func Chain(path ...Selector) Selector {
	return func(out chan<- *simplejson.Json, js *simplejson.Json) {
		cin := make(chan *simplejson.Json, 2)
		cin <- js
		cin <- nil
		cout := make(chan *simplejson.Json, 2)
		chain := func(i int, cout chan<- *simplejson.Json, cin <-chan *simplejson.Json) {
			j := 0
			for js := range cin {
				if js == nil {
					j++
					break
				}
				_cout := make(chan *simplejson.Json)
				go path[i](_cout, js)
				for js := range _cout {
					if js == nil {
						break
					}
					cout <- js
				}
			}
			cout <- nil
		}
		for i := range path {
			if i == len(path)-1 {
				go chain(i, out, cin)
			} else {
				go chain(i, cout, cin)
				cin = cout
				cout = make(chan *simplejson.Json, 2)
			}
		}
	}
}

func RecursiveDescent(out chan<- *simplejson.Json, js *simplejson.Json) {
	recDescent(out, js)
	out <- nil
}
func recDescent(out chan<- *simplejson.Json, js *simplejson.Json) {
	out <- js
	if a, err := js.Array(); err == nil {
		for i := range a {
			elem := js.GetIndex(i)
			recDescent(out, elem)
		}
	} else if m, err := js.Map(); err == nil {
		for k := range m {
			val := js.Get(k)
			recDescent(out, val)
		}
	}
}

func All(out chan<- *simplejson.Json, js *simplejson.Json) {
	if a, err := js.Array(); err == nil {
		for i := range a {
			out <- js.GetIndex(i)
		}
	} else if m, err := js.Map(); err == nil {
		for k := range m {
			out <- js.Get(k)
		}
	}
	out <- nil
}

func Key(key string) Selector {
	return func(out chan<- *simplejson.Json, js *simplejson.Json) {
		jschild, ok := js.CheckGet(key)
		if ok {
			out <- jschild
		}
		out <- nil
	}
}

func Index(i int) Selector {
	return func(out chan<- *simplejson.Json, js *simplejson.Json) {
		if len(js.MustArray()) > i {
			out <- js.GetIndex(i)
		}
		out <- nil
	}
}

func Has(sel ...Selector) Selector {
	return func(out chan<- *simplejson.Json, js *simplejson.Json) {
		if len(Lookup(js, sel...)) > 0 {
			out <- js
		}
		out <- nil
	}
}

func EqualString(x string, sel ...Selector) Selector {
	return func(out chan<- *simplejson.Json, js *simplejson.Json) {
		for _, jschild := range Lookup(js, sel...) {
			str, err := jschild.String()
			if err == nil && str == x {
				out <- js
				break
			}
		}
		out <- nil
	}
}

func EqualFloat64(x float64, sel ...Selector) Selector {
	return func(out chan<- *simplejson.Json, js *simplejson.Json) {
		for _, jschild := range Lookup(js, sel...) {
			f, err := jschild.Float64()
			if err == nil && f == x {
				out <- js
				break
			}
		}
		out <- nil
	}
}

func EqualBool(x bool, sel ...Selector) Selector {
	return func(out chan<- *simplejson.Json, js *simplejson.Json) {
		for _, jschild := range Lookup(js, sel...) {
			b, err := jschild.Bool()
			if err == nil && b == x {
				out <- js
				break
			}
		}
		out <- nil
	}
}
