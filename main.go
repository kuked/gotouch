package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

var (
	nocreate = flag.Bool("c", false, "")
	times    = flag.String("t", "", "")
	tregexp  = regexp.MustCompile(`^((\d{2})?\d{2})?(\d{8})(\.[0-5][0-9])?$`)
	names    []string
)

var usage = `usage: gotouch [options...] file...

Options:
  -c Not create new empty file even if that does not exists.
  -t [[CC]YY]MMDDhhmm[.SS]
     Change the access and the modification times.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	flag.Parse()
	if flag.NArg() == 0 {
		usageAndExit("")
	}

	// XXX 正規表現にマッチしたからといって日付として正しいとは限らない
	if *times != "" && !tregexp.MatchString(*times) {
		usageAndExit("gotouch: out of range or illegal time specification")
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
			// XXX マッチした部分文字列を確認して年を足すかどうか判定しないとね
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

func usageAndExit(message string) {
	if message != "" {
		fmt.Fprintln(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
