package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/donke/gotouch"
	"github.com/donke/manyflags"
)

var (
	a = flag.Bool("a", false, "")
	c = flag.Bool("c", false, "")
	m = flag.Bool("m", false, "")

	r = flag.String("r", "", "")
	t = flag.String("t", "", "")

	tregexp = regexp.MustCompile(`^(\d{2}|\d{4})?(\d{8})(\.([0-5][0-9]))?$`)
	names   []string
)

var usage = `usage: gotouch [options...] file...

Options:
  -a Change the access time of the file.
  -c Not create new empty file even if that does not exists.
  -m Change the modification time of the file.
  -r Use the access and modification time from the specified file.
  -t [[CC]YY]MMDDhhmm[.SS]
     Change the access and the modification times.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	manyflags.OverwriteArgs()

	flag.Parse()
	if flag.NArg() == 0 {
		usageAndExit("")
	}

	if *r != "" {
		if !exists(*r) {
			usageAndExit("gotouch: " + *r + ": No such file or directory")
		}
	}

	if *t != "" {
		if _, err := parseTimeSpecified(*t); err != nil {
			usageAndExit("gotouch: out of range or illegal time specification")
		}
	}

	for _, name := range flag.Args() {
		if exists(name) {
			names = append(names, name)
			continue
		}
		if *c {
			continue
		}
		if err := gotouch.Create(name); err != nil {
			fmt.Fprintln(os.Stderr, name+":", err)
			os.Exit(1)
		}
		names = append(names, name)
	}

	if *r != "" && *t == "" {
		for _, name := range names {
			if err := gotouch.UpdateTimeByFile(name, *r); err != nil {
				panic(err)
			}
		}
	}

	if *t != "" {
		var fn func(s string, t time.Time) error
		switch {
		case *a && !*m:
			fn = gotouch.UpdateAtime
		case !*a && *m:
			fn = gotouch.UpdateMtime
		default:
			fn = gotouch.UpdateTime
		}
		u, _ := parseTimeSpecified(*t)
		for _, name := range names {
			if err := fn(name, u); err != nil {
				panic(err)
			}
		}
	}

	os.Exit(0)
}

func parseTimeSpecified(t string) (time.Time, error) {
	if !tregexp.MatchString(t) {
		return time.Now(), errors.New("")
	}

	parsed, err := time.Parse("20060102150405-0700", timeValue(t))
	if err != nil {
		return time.Now(), err
	}

	return parsed, nil
}

func timeValue(t string) string {
	match := tregexp.FindStringSubmatch(t)
	ccyy := match[1]
	mmddhhmm := match[2]
	ss := match[4]
	v := ""
	if ccyy != "" {
		// year 2068 problem.
		if len(ccyy) == 2 {
			y, _ := strconv.Atoi(ccyy)
			if y >= 69 && y <= 99 {
				v = "19" + ccyy
			} else {
				v = "20" + ccyy
			}
		} else {
			v = ccyy
		}
	} else {
		y := time.Now()
		v = fmt.Sprint(y.Year())
	}
	v = v + mmddhhmm
	if ss != "" {
		v = v + ss
	} else {
		v = v + "00"
	}
	v = v + "+0900"
	return v
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
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
