package main

import (
	"os"

	"github.com/xchacha20-poly1305/gvgo"
)

func main() {
	var (
		based string

		major, minor, patch string
	)

	switch len(os.Args) {
	case 1, 2:
		fatal("too few arguments")
	case 5:
		patch = os.Args[4]
		fallthrough
	case 4:
		minor = os.Args[3]
		fallthrough
	case 3:
		major = os.Args[2]
		based = os.Args[1]
	default:
		fatal("too many arguments")
	}

	version, valid := gvgo.Parse(based)
	if !valid {
		fatal("Invalid version: " + based)
	}

	tryPlus := func(source *string, extra string) {
		next := gvgo.Plus(*source, extra)
		if extra != "0" && next == *source {
			fatal("Invalid argument: " + extra)
		}
		*source = next
	}
	tryPlus(&version.Major, major)
	tryPlus(&version.Minor, minor)
	tryPlus(&version.Patch, patch)

	_, _ = os.Stdout.WriteString(version.String())
}

func fatal(message string) {
	_, _ = os.Stderr.WriteString(message)
	os.Exit(1)
}
