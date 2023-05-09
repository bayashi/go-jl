package jl

import (
	"testing"
)

func FuzzJlProcess(f *testing.F) {
	f.Fuzz(func(t *testing.T, prettify bool, showErr bool, splitTab bool, data []byte) {
		options := &Options{
			Prettify: prettify,
			ShowErr:  showErr,
			SplitTab: splitTab,
		}

		Process(options, data)
	})
}
