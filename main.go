package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

var (
	nocreate = flag.Bool("c", false, "not create new empty file even if that does not exists.")
	times    = flag.String("t", "", "change the access and the modification times.")
	//timesregexp = regexp.MustCompile("^(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])([01][0-9]|2[0-3])([0-5][0-9])(\\.[0-5][0-9])?$")
	timesregexp = regexp.MustCompile("^(0[1-9]|1[0-2])(0[1-9]|[12][0-9]|3[01])([01][0-9]|2[0-3])([0-5][0-9])$")
	names       []string
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		usage()
	}

	if *times != "" && !timesregexp.MatchString(*times) {
		eprintln("goutouch: out of range or illegal time specification: MMDDhhmm")
		os.Exit(1)
	}

	for _, name := range flag.Args() {
		if !exists(name) && !createEmptyfile(name) {
			continue
		}
		names = append(names, name)
	}

	if *times != "" {
		for _, name := range names {
			y := fmt.Sprint(getThisYear())
			t, _ := time.Parse("20060102150405-0700", y+*times+"00+0900")
			os.Chtimes(name, t, t)
		}
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

func getThisYear() int {
	t := time.Now()
	return t.Year()
}

func usage() {
	eprintln("usage:")
	eprintln("gotouch [-c] [-t MMDDhhmm] file ...")
	os.Exit(1)
}

func eprintln(message string) {
	fmt.Fprintln(os.Stderr, message)
}
