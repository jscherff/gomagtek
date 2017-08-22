package gomagtek

import (
	"strings"
	"fmt"
	"os"
)

// ReportField represents one device property in a device report. It contains
// a descriptive name, the property value, and an error object to notify the
// caller when there was an error retrieving the property.
type ReportField struct {
	Flag string
	Field string
	Value string
	Error error
}

// Report represents a collection of ordered report fields for a single device.
type Report []ReportField

var FieldFlags = []string {
	"hn",
	"vid",
	"vn",
	"pid",
	"pn",
	"pv",
	"sid",
	"bs",
	"sn",
	"fsn",
	"dsn"}

var FlagFieldMap = map[string]string {
	"hn":	"hostname",
	"vid":	"vendor_id",
	"vn":	"vendor_name",
	"pid":	"product_id",
	"pn":	"product_name",
	"pv":	"product_ver",
	"sid":	"software_id",
	"bs":	"buffer_size",
	"sn":	"device_sn",
	"fsn":	"factory_sn",
	"dsn":	"descript_sn"}

var FieldTitleMap = map[string]string {
	"hostname":	"Host Name",
	"vendor_id":	"Vendor ID",
	"vendor_name":	"Vendor Name",
	"product_id":	"Product ID",
	"product_name":	"Product Name",
	"product_ver":	"Product Version",
	"software_id":	"Software ID",
	"buffer_size":	"Buffer Size",
	"device_sn":	"Device Serial Number",
	"factory_sn":	"Factory Serial Number",
	"descript_sn":	"Descriptor Serial Number"}

var FormatTypes = []string {
	"csv",
	"nvp"}

var FormatTitle = map[string]string {
	"csv":	"Comma-separated values",
	"nvp":	"Name-value pairs"}

// Report receives a slice of properties desired in the report and returns a
// populated report.
func (d *Device) Report(fields []string) (r Report, err error) {

	for _, f := range fields {

		f := strings.ToLower(f)
		rf := ReportField {Flag: f, Field: FlagFieldMap[f]}

		switch f {

		case "hn", FlagFieldMap["hn"]:
			rf.Value, rf.Error = os.Hostname()
		case "vid", FlagFieldMap["vid"]:
			rf.Value, rf.Error = d.GetVendorID(), nil
		case "vn", FlagFieldMap["vn"]:
			rf.Value, rf.Error = d.GetVendorName()
		case "pid", FlagFieldMap["pid"]:
			rf.Value, rf.Error = d.GetProductID(), nil
		case "pn", FlagFieldMap["pn"]:
			rf.Value, rf.Error = d.GetProductName()
		case "pv", FlagFieldMap["pv"]:
			rf.Value, rf.Error = d.GetProductVer()
		case "sid", FlagFieldMap["sid"]:
			rf.Value, rf.Error = d.GetSoftwareID()
		case "bs", FlagFieldMap["bs"]:
			rf.Value, rf.Error = d.GetBufferSize()
		case "sn", FlagFieldMap["sn"]:
			rf.Value, rf.Error = d.GetDeviceSN()
		case "fsn", FlagFieldMap["fsn"]:
			rf.Value, rf.Error = d.GetFactorySN()
		case "dsn", FlagFieldMap["dsn"]:
			rf.Value, rf.Error = d.GetDescriptSN()
		default:
			if err == nil {
				err = fmt.Errorf("%s: unsupported field(s):", getFunctionInfo())
			}
			err = fmt.Errorf("%v %s", err, f)
			continue
		}

		r = append(r, rf)
	}

	return r, err
}

func (r Report) CSV(raw bool) (out string) {

	var fields, values []string

	for _, rf := range r {
		fields = append(fields, rf.Field)
		values = append(values, fmt.Sprintf("%q", rf.Value))
	}

	if !raw {
		out += fmt.Sprintf("%s\n", strings.Join(fields, ","))
	}

	out += fmt.Sprintf("%s", strings.Join(values, ","))

	return out
}

func (r Report) NVP(raw bool) (out string) {

	for _, rf := range r {
		if !raw {
			out += fmt.Sprintf("%s:", rf.Field)
		}
		out += fmt.Sprintf("%s\n", rf.Value)
	}

	return out
}
