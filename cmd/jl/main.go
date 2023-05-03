package main

import (
	"bufio"
	"os"

	"github.com/bayashi/go-jl"
)

const (
	cmd string = "jl"

	exitOK  int = 0
	exitErr int = 1
)

func main() {
	o := &options{}
	o.parseArgs()

	po := &jl.Options{
		Prettify: o.prettify,
		ShowErr:  o.showErr,
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		in := s.Bytes()
		result := jl.Process(po, in)
		os.Stdout.Write(result)
		os.Stdout.WriteString("\n")
	}

	os.Exit(exitOK)
}
