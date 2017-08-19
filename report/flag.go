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

var includeUsage = fmt.Sprintf("Include `<field list>` in report (comma-separated):" +
	strings.Repeat("\n\t%q\t%s", 11),
	"hostname",	"Host Name",
	"vendor_id",	"Vendor ID",
	"vendor_name",	"Vendor Name",
	"product_id",	"Product ID",
	"product_name",	"Product Name",
	"product_ver",	"Product Version",
	"software_id",	"Software ID",
	"buffer_size",	"Buffer Size",
	"device_sn",	"Device Serial Number",
	"factory_sn",	"Device Factory Serial Number",
	"descrip_sn",	"Device Descriptor Serial Number")

var (
	fsMode = flag.NewFlagSet("mode", flag.ExitOnError)
	fReportMode = fsMode.Bool("report", false, "Report mode")
	fConfigMode = fsMode.Bool("config", false, "Config mode")
)

var (
	fsReport = flag.NewFlagSet("report_fields", flag.ExitOnError)
	fInclude = fsReport.String("include", "", includeUsage)
	fFormatCSV = fsReport.Bool("format_csv", true, "Write output in CSV format")
	fFormatNVP = fsReport.Bool("format_nvp", false, "Write output as name/value pairs")
	fWriteScreen = fsReport.Bool("write_screen", false, "Write output to screen")
	fWriteFile = fsReport.String ("write_file", "", "Write output to `<file>`")
)

var (
	fsConfig = flag.NewFlagSet("config_options", flag.ExitOnError)
	fEraseSN = fsConfig.Bool("erase_sn", false, "Erase device serial number")
	fSetSnString = fsConfig.String("set_sn_string", "", "Set device serial number to `<string>`")
	fSetSnFactory = fsConfig.Int("set_sn_factory", 7, "Set device serial number to `<n>` bytes of factory SN")
	fSetSnServer = fsConfig.String("set_sn_server", "", "Set device serial number from server `<url>`")
	fFetchIfSN = fsConfig.Bool("set_only_empty", true, "Set device serial number ONLY if empty")
	fUsbReset = fsConfig.Bool("usb_reset", false, "Perform a USB reset")
	fDevReset = fsConfig.Bool("dev_reset", false, "Perform a device reset")
)
