package gomagtek

// https://beta.golang.org/doc/go1.8#language
// https://play.golang.org/p/QNArOeqy94


import (
	"encoding/json"
	"encoding/xml"
	"os"
)

type DeviceInfo struct {
	HostName	string
	DeviceSN	string
	VendorID	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	ProductID	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	SoftwareID	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	VendorName	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	ProductName	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	ProductVer	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	FactorySN	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	DescriptSN	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	BusNumber	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	BusAddress	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	USBSpec		string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	USBClass	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	USBSubclass	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	USBProtocol	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	DeviceSpeed	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	DeviceVer	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	MaxPktSize	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	BufferSize	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
}

type DeviceInfoMin struct {
	HostName	string
	DeviceSN	string
	VendorID	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	ProductID	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	SoftwareID	string	`json:",omitempty" xml:",omitempty" csv:",omitempty"`
	VendorName	string	`json:"-" xml:"-" csv:"-"`
	ProductName	string	`json:"-" xml:"-" csv:"-"`
	ProductVer	string	`json:"-" xml:"-" csv:"-"`
	FactorySN	string	`json:"-" xml:"-" csv:"-"`
	DescriptSN	string	`json:"-" xml:"-" csv:"-"`
	BusNumber	string	`json:"-" xml:"-" csv:"-"`
	BusAddress	string	`json:"-" xml:"-" csv:"-"`
	USBSpec		string	`json:"-" xml:"-" csv:"-"`
	USBClass	string	`json:"-" xml:"-" csv:"-"`
	USBSubclass	string	`json:"-" xml:"-" csv:"-"`
	USBProtocol	string	`json:"-" xml:"-" csv:"-"`
	DeviceSpeed	string	`json:"-" xml:"-" csv:"-"`
	DeviceVer	string	`json:"-" xml:"-" csv:"-"`
	MaxPktSize	string	`json:"-" xml:"-" csv:"-"`
	BufferSize	string	`json:"-" xml:"-" csv:"-"`
}

var ImportMap = map[string]string {
	"host_name":	"HostName",
	"device_sn":	"DeviceSN",
	"vendor_id":	"VendorID",
	"product_id":	"ProductID",
	"software_id":	"SoftwareID",
	"vendor_name":	"VendorName",
	"product_name":	"ProductName",
	"product_ver":	"ProductVer",
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
	"DeviceSN":	"device_sn",
	"VendorID":	"vendor_id",
	"ProductID":	"product_id",
	"SoftwareID":	"software_id",
	"VendorName":	"vendor_name",
	"ProductName":	"product_name",
	"ProductVer":	"product_ver",
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

func NewDeviceInfo(d *Device) (i *DeviceInfo, errs []error) {

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
	if i.DeviceSN, e = d.GetDeviceSN(); e != nil {errs = append(errs, e)}
	if i.SoftwareID, e = d.GetSoftwareID(); e != nil {errs = append(errs, e)}
	if i.VendorName, e = d.GetVendorName(); e != nil {errs = append(errs, e)}
	if i.ProductName, e = d.GetProductName(); e != nil {errs = append(errs, e)}
	if i.ProductVer, e = d.GetProductVer(); e != nil {errs = append(errs, e)}
	if i.FactorySN, e = d.GetFactorySN(); e != nil {errs = append(errs, e)}
	if i.DescriptSN, e = d.GetDescriptSN(); e != nil {errs = append(errs, e)}
	if i.BufferSize, e = d.GetBufferSize(); e != nil {errs = append(errs, e)}

	return i, errs
}

func NewDeviceInfoFromXML(x []byte) (*DeviceInfo, error) {
	i := new(DeviceInfo)
	err := xml.Unmarshal(x, i)
	return i, err
}

func NewDeviceInfoFromJSON(j []byte) (*DeviceInfo, error) {
	i := new(DeviceInfo)
	err := json.Unmarshal(j, i)
	return i, err
}

func (i *DeviceInfo) Prune(fields []string) (error) {
	// Pass a list of fields desired. Iterate through struct fields.
	// If struct field is not in the list, set its value to an empty
	// string so it will not be included in JSON/XML/CSV export. See
	// stackoverflow.com/questions/18926303/iterate-through-a-struct-in-go
	// for hints.
	return nil
}

func (i *DeviceInfo) JSON(min bool) ([]byte, error) {
	if min {return json.Marshal(DeviceInfoMin(*i))}
	return json.Marshal(i)
}

func (i *DeviceInfo) XML(min bool) ([]byte, error) {
	if min {return xml.Marshal(DeviceInfoMin(*i))}
	return xml.Marshal(i)
}

func (i *DeviceInfo) FXML(min bool) ([]byte, error) {
	if min {return xml.MarshalIndent(DeviceInfoMin(*i), "", "\t")}
	return xml.MarshalIndent(i, "", "\t")
}
