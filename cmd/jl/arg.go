package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	flag "github.com/spf13/pflag"
)

var (
	version     = ""
	installFrom = "Source"
)

type options struct {
	noPrettify bool
	showErr    bool
	splitTab   bool
	splitLF    bool
	skip       int
}

func (o *options) parseArgs() {
	var flagHelp bool
	var flagVersion bool
	flag.BoolVarP(&flagHelp, "help", "h", false, "Display help (This message) and exit")
	flag.BoolVarP(&flagVersion, "version", "v", false, "Display version and build info and exit")
	flag.BoolVarP(&o.noPrettify, "no-prettify", "P", false, "Not prettify the JSON. Prettified by default")
	flag.BoolVarP(&o.showErr, "show-error", "e", false, "Set this option to show errors, muted by default")
	flag.BoolVarP(&o.splitTab, "split-tab", "t", false, "Split tabs in each element")
	flag.BoolVarP(&o.splitLF, "split-lf", "n", false, "Split line-feed \\n in each element")
	flag.IntVarP(&o.skip, "skip", "", 0, "Skip to parse JSON if the length of the source JSON less than this")
	flag.Parse()

	if flagHelp {
		putHelp(fmt.Sprintf("[%s] Version %s", cmd, getVersion()))
		os.Exit(exitOK)
	}

	if flagVersion {
		putErr(versionDetails())
		os.Exit(exitOK)
	}
}

func versionDetails() string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	compiler := runtime.Version()

	return fmt.Sprintf(
		"Version %s - %s.%s (compiled:%s, %s)",
		getVersion(),
		goos,
		goarch,
		compiler,
		installFrom,
	)
}

func getVersion() string {
	if version != "" {
		return version
	}
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return "Unknown"
	}

	return i.Main.Version
}

func putErr(message ...interface{}) {
	fmt.Fprintln(os.Stderr, message...)
}

func putUsage() {
	putErr(fmt.Sprintf("Usage: cat some.json | %s", cmd))
}

func putHelp(message string) {
	putErr(message)
	putUsage()
	putErr("Options:")
	flag.PrintDefaults()
}
