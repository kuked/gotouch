package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var names []string

	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	for _, name := range flag.Args() {
		f, err := os.Create(name)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		names = append(names, name)
	}
}

func usage() {
	eprintln("usage:")
	eprintln("gotouch file ...")
	os.Exit(1)
}

func eprintln(message string) {
	fmt.Fprintln(os.Stderr, message)
}
