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
	prettify bool
	showErr  bool
}

func (o *options) parseArgs() {
	var flagHelp bool
	var flagVersion bool
	flag.BoolVarP(&flagHelp, "help", "h", false, "Display help (This message) and exit")
	flag.BoolVarP(&flagVersion, "version", "v", false, "Display version and build info and exit")
	flag.BoolVarP(&o.prettify, "prettify", "p", false, "Prettify the JSON")
	flag.BoolVarP(&o.showErr, "show-error", "e", false, "Set this option to show errors")
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
	putErr(fmt.Sprintf("Usage: %s [OPTIONS] FILE", cmd))
}

func putHelp(message string) {
	putErr(message)
	putUsage()
	putErr("Options:")
	flag.PrintDefaults()
}
