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
)

func main() {
	o := &options{}
	o.parseArgs()

	if term.IsTerminal(int(syscall.Stdin)) {
		os.Exit(exitOK)
	}

	po := &jl.Options{
		Prettify: o.prettify,
		ShowErr:  o.showErr,
		SplitTab: o.splitTab,
		SplitLF:  o.splitLF,
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		in := s.Bytes()
		result, err := jl.Process(po, in)
		if err != nil && po.ShowErr {
			os.Stderr.Write([]byte(err.Error()))
			os.Stderr.WriteString("\n")
		}
		os.Stdout.Write(result)
		os.Stdout.WriteString("\n")
	}

	os.Exit(exitOK)
}
