package gomagtek

import (
	"encoding/json"
	"fmt"
	"os"
)

type DeviceInfo struct {
	HostName	string	`json:"host_name"`
	VendorID	string	`json:"vendor_id"`
	ProductID	string	`json:"product_id"`
	VendorName	string	`json:"vendor_name"`
	ProductName	string	`json:"product_name"`
	ProductVer	string	`json:"product_ver"`
	SoftwareID	string	`json:"software_id"`
	DeviceSN	string	`json:"device_sn"`
	FactorySN	string	`json:"factory_sn"`
	DescriptSN	string	`json:"descript_sn"`
	BusNumber	string	`json:"bus_number"`
	BusAddress	string	`json:"bus_address"`
	USBSpec		string	`json:"usb_spec"`
	USBClass	string	`json:"usb_class"`
	USBSubclass	string	`json:"usb_subclass"`
	USBProtocol	string	`json:"usb_protocol"`
	DeviceSpeed	string	`json:"device_speed"`
	DeviceVer	string	`json:"device_ver"`
	MaxPktSize	string	`json:"max_pkt_size"`
	BufferSize	string	`json:"buffer_size"`
	Errors		[]error	`json:"errors"`
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

func NewDeviceInfo(d* Device) (i *DeviceInfo, err error) {

	var e error
	var errs []error

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

	i.HostName, e = os.Hostname(); errs = append(errs, e)
	i.VendorName, e = d.GetVendorName(); errs = append(errs, e)
	i.ProductName, e = d.GetProductName(); errs = append(errs, e)
	i.ProductVer, e = d.GetProductVer(); errs = append(errs, e)
	i.SoftwareID, e = d.GetSoftwareID(); errs = append(errs, e)
	i.DeviceSN, e = d.GetDeviceSN(); errs = append(errs, e)
	i.FactorySN, e = d.GetFactorySN(); errs = append(errs, e)
	i.DescriptSN, e = d.GetDescriptSN(); errs = append(errs, e)
	i.BufferSize, e = d.GetBufferSize(); errs = append(errs, e)

	for _, e = range errs {
		if e != nil {
			i.Errors = append(i.Errors, e)
		}
	}

	if len(i.Errors) > 0 {
		err = fmt.Errorf("%s: errs recorded in object's Errors field", getFunctionInfo())
	}

	return i, err
}

func (i *DeviceInfo) JSON() (s string, err error) {
	j, err := json.Marshal(i)
	return string(j), err
}
