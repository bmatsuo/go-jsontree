// Copyright 2013, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// jsonpath.go [created: Fri, 21 Jun 2013]

/*
jsonpath is a command line utility for manipulating and filtering json data.
jsonpath reads json objects from standard input and prints selected data to
standard output.

	$ echo '{"thing":"hello"}' '{"thing":"world"}' | jsonpath $.thing
	"hello"
	"world"

by default, selected strings are printed as json strings. to print the decoded
string instead of the json representation use the -printstrings option

	$ echo '{"thing":"hello"}' '{"thing":"world"}' | jsonpath -printstrings $.thing
	hello
	world

multiple paths can be selected for each object. these objects can be printed
on the same line (tab separated) for easier scripting.

	$ echo '{"date":"2012-12-12","event":"apocalypse"}' > test.json
	$ echo '{"date":"2012-12-13","event":"false alarm"}' >> test.json
	$ cat test.json | jsonpath -oneline -printstrings $.date $.event
	2012-12-12	apocalypse
	2012-12-13	false alarm
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/bmatsuo/go-jsonpath"
	"github.com/bmatsuo/go-simplejson"
)

func main() {
	oneline := flag.Bool("oneline", false, "one line printed per input object")
	onelinesep := flag.String("sep", "\t", "result separator when -oneline is given")
	printstrings := flag.Bool("printstrings", false, "do not marshal selected strings as json")
	mustexist := flag.Bool("mustexist", true, "exits with non-zero status if a selector has no results")
	pretty := flag.Bool("p", false, "pretty-print output")
	flag.Parse()
	paths := flag.Args()
	if len(paths) < 1 {
		fmt.Fprintf(os.Stderr, "usage: %s PATH ...\n", os.Args[0])
		os.Exit(1)
	}
	selectors := make([]jsonpath.Selector, len(paths))
	for i := range paths {
		sel, err := jsonpath.Parse(paths[i])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		selectors[i] = sel
	}
	dec := json.NewDecoder(os.Stdin)
	exitcode := 0
	for cont := true; cont; {
		js := new(simplejson.Json)
		err := dec.Decode(js)
		switch err {
		case nil:
		case io.EOF:
			cont = false
			continue
		default:
			fmt.Fprintln(os.Stderr, err)
			exitcode = 1
			cont = false
			continue
		}
		first := true
		for _, sel := range selectors {
			results := jsonpath.Lookup(js, sel)
			if len(results) == 0 && *mustexist {
				exitcode = 1
			}
			for i := range results {
				if *oneline {
					if first {
						first = false
					} else {
						fmt.Print(*onelinesep)
					}
				}
				if *printstrings {
					if str, ok := results[i].Data.(string); ok {
						if *oneline {
							fmt.Print(str)
						} else {
							fmt.Println(str)
						}
						continue
					}
				}
				var p []byte
				var err error
				if *pretty {
					p, err = json.MarshalIndent(results[i], "", "\t")
				} else {
					p, err = json.Marshal(results[i])
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
				} else {
					if *oneline {
						fmt.Print(string(p))
					} else {
						fmt.Println(string(p))
					}
				}
			}
		}
		if *oneline {
			fmt.Println()
		}
	}
	os.Exit(exitcode)
}
