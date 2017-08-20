package main

import (
	"strings"
	"flag"
	"fmt"
)

var reportFields = map[string]string{
	"hostname":	"Host Name",
	"vendor_id":	"Vendor ID",
	"vendor_name":	"Vendor Name",
	"product_id":	"Product ID",
	"product_name":	"Product Name",
	"product_ver":	"Product Version",
	"software_id":	"Software ID",
	"buffer_size":	"Buffer Size",
	"device_sn":	"Device Serial Number",
	"factory_sn":	"Device Factory Serial Number",
	"descrip_sn":	"Device Descriptor Serial Number"}

var includeUsage = fmt.Sprintf("Include `<fields>` in report (comma-separated list):" + strings.Repeat("\n\t%q\t%s", 11),
	"hn",	"Host Name",
	"vid",	"Vendor ID",
	"vn",	"Vendor Name",
	"pid",	"Product ID",
	"pn",	"Product Name",
	"pv",	"Product Version",
	"sid",	"Software ID",
	"bs",	"Buffer Size",
	"sn",	"Device Serial Number",
	"fsn",	"Factory Serial Number",
	"dsn",	"Descriptor Serial Number")

var formatUsage = fmt.Sprintf("Write report output in `<format>` format:" + strings.Repeat("\n\t%q\t%s", 2),
	"csv",	"Comma-separated values",
	"nvp",	"Name-value pairs")

var (
	fsMode = flag.NewFlagSet("mode", flag.ExitOnError)
	fModeReport = fsMode.Bool("report", false, "Report mode")
	fModeConfig = fsMode.Bool("config", false, "Config mode")
	fModeReset = fsMode.Bool("reset", false, "Reset mode")
)

var (
	fsReport = flag.NewFlagSet("report", flag.ExitOnError)
	fReportInc = fsReport.String("include", "", includeUsage)
	fReportFmt = fsReport.String("format", "", formatUsage)
	fReportFile = fsReport.String ("file", "", "Write output to `<file>`")
	fReportRaw = fsReport.Bool("raw", false, "Write output without headings")
	fReportStdout = fsReport.Bool("stdout", false, "Write output to stdout")
)

var (
	fsConfig = flag.NewFlagSet("config", flag.ExitOnError)
	fConfigErase = fsConfig.Bool("erase", false, "Erase serial number")
	fConfigSet = fsConfig.String("set", "", "Set serial number to `<string>`")
	fConfigUrl = fsConfig.String("url", "", "Set serial number from URL `<url>`")
	fConfigCopy = fsConfig.Int("copy", 7, "Copy `<n>` characters of factory SN to device SN")
	fConfigEmpty = fsConfig.Bool("empty", true, "Set serial number ONLY if it's empty")
)

var (
	fsReset = flag.NewFlagSet("reset", flag.ExitOnError)
	fResetUsb = fsReset.Bool("usb", false, "Perform a USB reset")
	fResetDev = fsReset.Bool("dev", false, "Perform a device reset")
)
