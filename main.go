package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	nocreate = flag.Bool("c", false, "not create new empty file even if that does not exists.")
	names    []string
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	for _, name := range flag.Args() {
		if !exists(name) && !createEmptyfile(name) {
			continue
		}
		names = append(names, name)
	}

	for _, name := range names {
		t, _ := time.Parse("2006-Jan-02", "2010-Oct-10")
		os.Chtimes(name, t, t)
	}

	os.Exit(0)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func createEmptyfile(filename string) (create bool) {
	if *nocreate {
		return false
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, filename+":", err)
		os.Exit(1)
	}
	f.Close()

	return true
}

func usage() {
	eprintln("usage:")
	eprintln("gotouch file ...")
	os.Exit(1)
}

func eprintln(message string) {
	fmt.Fprintln(os.Stderr, message)
}
