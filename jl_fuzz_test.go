package jl

import (
	"testing"
)

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
