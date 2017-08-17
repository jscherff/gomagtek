package gomagtek

import "github.com/google/gousb"
import "fmt"

/*
 * The gomagtek DeviceDescriptor specifies some basic information about
 * the USB device, such as the supported USB version, maximum packet size,
 * vendor and product IDs and the number of possible configurations the
 * device can have. This differs from the descriptor provided by gousb in
 * that it includes all the fields in raw format.
 */
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
	SerialNumIndex uint8		// Index of Serial Number String Descriptor
	NumConfigurations uint8		// Number of Possible Configurations
}

/*
 * Construct a new gomagtek DeviceDescriptor from a gousb Device.
 */
func NewDeviceDescriptor(d *gousb.Device) (ndd *DeviceDescriptor, err error) {

	ndd = new(DeviceDescriptor)
	data := make([]byte, BufferSizeDeviceDescriptor)

	_, err = d.Control(
		RequestTypeStandardDeviceIn,
		RequestGetDescriptor,
		TypeDeviceDescriptor,
		InterfaceNumber,
		data)

	if err == nil {

		*ndd = DeviceDescriptor {
			data[0],
			data[1],
			uint16(data[2]) + (uint16(data[3])<<8),
			data[4],
			data[5],
			data[6],
			data[7],
			uint16(data[8]) + (uint16(data[9])<<8),
			uint16(data[10]) + (uint16(data[11])<<8),
			uint16(data[12]) + (uint16(data[13])<<8),
			data[14],
			data[15],
			data[16],
			data[17]}
	} else {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return ndd, err
}

/*
 * The gomagtek ConfigDescriptor represents the active configuration of
 * the USB device. A device can have several different configurations,
 * though most have only one.
 */
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

/*
 * Construct a new gomagtek ConfigDescriptor from a gousb Device.
 */
func NewConfigDescriptor(d *gousb.Device) (ncd *ConfigDescriptor, err error) {

	ncd = new(ConfigDescriptor)
	data := make([]byte, BufferSizeConfigDescriptor)

	_, err = d.Control(
		RequestTypeStandardDeviceIn,
		RequestGetDescriptor,
		TypeConfigDescriptor,
		InterfaceNumber,
		data)

	if err == nil {

		*ncd = ConfigDescriptor {
			data[0],
			data[1],
			uint16(data[2]) + (uint16(data[3]) << 8),
			data[4],
			data[5],
			data[6],
			data[7],
			data[8]}
	} else {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return ncd, err
}
