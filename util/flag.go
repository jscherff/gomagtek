package main

import (
	"github.com/jscherff/gomagtek"
	"strings"
	"strconv"
	"flag"
	"fmt"
)

// stringValue is a string flag value that knows if it has been set.
type stringValue struct {
	value string
	set bool
}

func (sv *stringValue) String() string {
	return sv.value
}

func (sv *stringValue) Set(s string) (err error) {
	sv.value = s
	sv.set = true
	return err
}

func (sv *stringValue) IsSet() bool {
	return sv.set
}

// intValue is an integer flag value that knows if it has been set.
type intValue struct {
	value int
	set bool
}

func (iv *intValue) String() string {
	return strconv.Itoa(iv.value)
}

func (iv *intValue) Set(s string) (err error) {
	iv.value, err = strconv.Atoi(s)
	iv.set = true
	return err
}

func (iv *intValue) IsSet() bool {
	return iv.set
}

// intValue is a string slice flag value that knows if it has been set.
type reportFields struct {
	values []string
	set bool
}

func (rf *reportFields) String() string {
	return strings.Join(rf.values, ",")
}

func (rf *reportFields) Set(s string) (err error) {
	rf.values = strings.Split(s, ",")
	rf.set = true
	return err
}

func (rf *reportFields) IsSet() bool {
	return rf.set
}

var includeUsage = "Include `<fields>` in report (comma-separated list):"
var formatUsage = "Write report output in `<format>` format:"

var (
	fsMode = flag.NewFlagSet("mode", flag.ExitOnError)
	fModeReport = fsMode.Bool("report", false, "Report mode")
	fModeConfig = fsMode.Bool("config", false, "Config mode")
	fModeReset = fsMode.Bool("reset", false, "Reset mode")
)

var (
	fsReport = flag.NewFlagSet("report", flag.ExitOnError)
	fReportFile = fsReport.String ("file", "", "Write output to `<file>`")
	fReportRaw = fsReport.Bool("raw", false, "Write output without headings")
	fReportStdout = fsReport.Bool("stdout", false, "Write output to stdout")
	fReportFormat *string
	fReportInclude *string
)

var (
	fsConfig = flag.NewFlagSet("config", flag.ExitOnError)
	fConfigErase = fsConfig.Bool("erase", false, "Erase serial number")
	fConfigEmpty = fsConfig.Bool("empty", true, "Set serial number ONLY if it's empty")
	fConfigSet = fsConfig.String("set", "", "Set serial number to `<string>`")
	fConfigUrl = fsConfig.String("url", "", "Set serial number from URL `<url>`")
	fConfigCopy = fsConfig.Int("copy", 0, "Copy `<n>` characters of factory SN to device SN")

)

var (
	fsReset = flag.NewFlagSet("reset", flag.ExitOnError)
	fResetUsb = fsReset.Bool("usb", false, "Perform a USB reset")
	fResetDev = fsReset.Bool("dev", false, "Perform a device reset")
)

func init() {

	for _, f := range gomagtek.FieldFlags {
		includeUsage += fmt.Sprintf("\n\t%q\t%s", f, gomagtek.FieldTitleMap[gomagtek.FlagFieldMap[f]])
	}

	//fReportInclude = fsReport.String("include", strings.Join(gomagtek.FieldDefaults, ","), includeUsage)
	fReportInclude = fsReport.String("include", "", includeUsage)

	for _, t := range gomagtek.FormatTypes {
		formatUsage += fmt.Sprintf("\n\t%q\t%s", t, gomagtek.FormatTitle[t])
	}

	fReportFormat = fsReport.String("format", "", formatUsage)
}
