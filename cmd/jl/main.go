package main

import (
	"bufio"
	"os"
	"syscall"

	"github.com/bayashi/go-jl"
	"golang.org/x/term"
)

const (
	cmd string = "jl"

	exitOK  int = 0
	exitErr int = 1

	errPrefix = "JLERROR: "
)

func main() {
	o := &options{}
	o.parseArgs()

	if term.IsTerminal(int(syscall.Stdin)) {
		os.Exit(exitOK)
	}

	po := &jl.Options{
		NoPrettify: o.noPrettify,
		ShowErr:    o.showErr,
		SplitTab:   o.splitTab,
		SplitLF:    o.splitLF,
	}

	r := bufio.NewScanner(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	e := bufio.NewWriter(os.Stderr)
	for r.Scan() {
		in := r.Bytes()
		result, err := jl.Process(po, in)
		if err != nil && po.ShowErr {
			e.Write([]byte(errPrefix + err.Error()))
			e.Flush()
		}
		result = append(result, '\n')
		w.Write(result)
		w.Flush()
	}

	os.Exit(exitOK)
}
