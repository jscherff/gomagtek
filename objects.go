package main

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
