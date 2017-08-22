package gomagtek

import (
	"encoding/json"
	"encoding/xml"
	"os"
)

type DeviceInfo struct {
	HostName	string
	VendorID	string
	ProductID	string
	VendorName	string
	ProductName	string
	ProductVer	string
	SoftwareID	string
	DeviceSN	string
	FactorySN	string
	DescriptSN	string
	BusNumber	string
	BusAddress	string
	USBSpec		string
	USBClass	string
	USBSubclass	string
	USBProtocol	string
	DeviceSpeed	string
	DeviceVer	string
	MaxPktSize	string
	BufferSize	string
}

var ImportMap = map[string]string {
	"host_name":	"HostName",
	"vendor_id":	"VendorID",
	"product_id":	"ProductID",
	"vendor_name":	"VendorName",
	"product_name":	"ProductName",
	"product_ver":	"ProductVer",
	"software_id":	"SoftwareID",
	"device_sn":	"DeviceSN",
	"factory_sn":	"FactorySN",
	"descript_sn":	"DescriptSN",
	"bus_number":	"BusNumber",
	"bus_address":	"BusAddress",
	"usb_spec":	"USBSpec",
	"usb_class":	"USBClass",
	"usb_subclass":	"USBSubclass",
	"usb_protocol":	"USBProtocol",
	"device_speed":	"DeviceSpeed",
	"device_ver":	"DeviceVer",
	"max_pkt_size":	"MaxPktSize",
	"buffer_size":	"BufferSize"}

var ExportMap = map[string]string {
	"HostName":	"host_name",
	"VendorID":	"vendor_id",
	"ProductID":	"product_id",
	"VendorName":	"vendor_name",
	"ProductName":	"product_name",
	"ProductVer":	"product_ver",
	"SoftwareID":	"software_id",
	"DeviceSN":	"device_sn",
	"FactorySN":	"factory_sn",
	"DescriptSN":	"descript_sn",
	"BusNumber":	"bus_number",
	"BusAddress":	"bus_address",
	"USBSpec":	"usb_spec",
	"USBClass":	"usb_class",
	"USBSubclass":	"usb_subclass",
	"USBProtocol":	"usb_protocol",
	"DeviceSpeed":	"device_speed",
	"DeviceVer":	"device_ver",
	"MaxPktSize":	"max_pkt_size",
	"BufferSize":	"buffer_size"}

var DefaultFields = []string {
	"HostName",
	"VendorID",
	"ProductID",
	"ProductVer",
	"SoftwareID",
	"DeviceSN"}

func NewDeviceInfo(d* Device) (i *DeviceInfo, errs []error) {

	var e error

	i = &DeviceInfo {
		VendorID:	d.GetVendorID(),
		ProductID:	d.GetProductID(),
		BusNumber:	d.GetBusNumber(),
		BusAddress:	d.GetBusAddress(),
		USBSpec:	d.GetUSBSpec(),
		USBClass:	d.GetUSBClass(),
		USBSubclass:	d.GetUSBSubclass(),
		USBProtocol:	d.GetUSBProtocol(),
		DeviceSpeed:	d.GetDeviceSpeed(),
		DeviceVer:	d.GetDeviceVer(),
		MaxPktSize:	d.GetMaxPktSize()}

	if i.HostName, e = os.Hostname(); e != nil {errs = append(errs, e)}
	if i.VendorName, e = d.GetVendorName(); e != nil {errs = append(errs, e)}
	if i.ProductName, e = d.GetProductName(); e != nil {errs = append(errs, e)}
	if i.ProductVer, e = d.GetProductVer(); e != nil {errs = append(errs, e)}
	if i.SoftwareID, e = d.GetSoftwareID(); e != nil {errs = append(errs, e)}
	if i.DeviceSN, e = d.GetDeviceSN(); e != nil {errs = append(errs, e)}
	if i.FactorySN, e = d.GetFactorySN(); e != nil {errs = append(errs, e)}
	if i.DescriptSN, e = d.GetDescriptSN(); e != nil {errs = append(errs, e)}
	if i.BufferSize, e = d.GetBufferSize(); e != nil {errs = append(errs, e)}

	return i, errs
}

func (i *DeviceInfo) JSON() (s string, err error) {
	j, err := json.Marshal(i)
	return string(j), err
}

func (i *DeviceInfo) XML() (s string, err error) {
	x, err := xml.Marshal(i)
	return string(x), err
}

func (i *DeviceInfo) FXML() (s string, err error) {
	x, err := xml.MarshalIndent(i, "", "\t")
	return string(x), err
}
