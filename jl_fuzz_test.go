package jl

import (
	"testing"
)

// go test -fuzz=FuzzJlProcess -fuzztime=5s
func FuzzJlProcess(f *testing.F) {
	f.Fuzz(func(t *testing.T, noPrettify bool, showErr bool, splitTab bool, splitLF bool, data []byte) {
		options := &Options{
			NoPrettify: noPrettify,
			ShowErr:    showErr,
			SplitTab:   splitTab,
			SplitLF:    splitLF,
		}

		Process(options, data)
	})
}
