package main

const (
	MagtekVendorID uint16 = 0x0801
	SureswipeKbPID uint16 = 0x0001
	SureswipeHidPID uint16 = 0x0002

	DynamagSwipeKbPID uint16 = 0x0001
	DynamagInsertKbPID uint16 = 0x0001
	DynamagSwipeHidPID uint16 = 0x0011
	DynamagInsertHidPID uint16 = 0x0013
	DynamagWirelessHidPID uint16 = 0x0014

	RequestTypeVendorDeviceOut uint8 = 0x21
	RequestTypeVendorDeviceIn uint8 = 0xA1
	RequestTypeStandardDeviceOut uint8 = 0x00
	RequestTypeStandardDeviceIn uint8 = 0x80
	RequestTypeStandardInterfaceIn uint8 = 0x81

	RequestGetReport uint8 = 0x01
	RequestSetReport uint8 = 0x09
	RequestGetDescriptor uint8 = 0x06

	TypeDeviceDescriptor uint16 = 0x0100
	TypeConfigDescriptor uint16 = 0x0200
	TypeHidDescriptor uint16 = 0x2200
	TypeFeatureReport uint16 = 0x0300

	InterfaceNumber uint16 = 0x0000

	BufferSizeSureswipe int = 24
	BufferSizeDynamag int = 60
	BufferSizeDeviceDescriptor int = 18
	BufferSizeConfigDescriptor int = 9

	CommandGetProperty uint8 = 0x00
	CommandSetProperty uint8 = 0x01
	CommandResetDevice uint8 = 0x02

	ResultCodeSuccess uint8 = 0x00
	ResultCodeFailure uint8 = 0x01
	ResultCodeBadParam uint8 = 0x02

	PropertySoftwareID uint8 = 0x00
	PropertySerialNum uint8 = 0x01
)

var BufferSizes = []int {24, 60}

type DeviceDescriptor struct {
	Length uint8			// Size of the Descriptor in Bytes
	DescriptorType uint8		// Device Descriptor Type (0x01)
	UsbSpecification uint16		// BCD of Device USB Specification Number
	DeviceClass uint8		// Device Class Code or Vendor Specified
	DeviceSubClass uint8		// Subclass Code Assigned by USB Org
	DeviceProtocol uint8		// Protocol Code Assigned by USB Org
	MaxPacketSize uint8		// Maximum Packet Size for Zero Endpoint
	VendorID uint16			// Vendor ID Assigned by USB Org
	ProductID uint16		// Product ID (Assigned by Manufacturer)
	DeviceReleaseNumber uint16	// BCD of Device Release Number
	ManufacturerIndex uint8		// Index Manufacturer String Descriptor
	ProductIndex uint8		// Index of Product String Descriptor
	SerialNumberIndex uint8		// Index of Serial Number String Descriptor
	NumConfigurations uint8		// Number of Possible Configurations
}

type ConfigDescriptor struct {
	Length uint8			// Size of Descriptor in Bytes
	DescriptorType uint8		// Configuration Descriptor Type (0x02)
	TotalLength uint16		// Total Length of Data Returned
	NumInterfaces uint8		// Number of Interfaces
	ConfigurationValue uint8	// Value to Select This Configuration
	ConfigurationIndex uint8	// Index of String Descriptor for Configuration
	Attributes uint8		// Bitmap of Power Attributes
	MaxPower uint8			// Maximum Power Consumption in 2mA units
}
