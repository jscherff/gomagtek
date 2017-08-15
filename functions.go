package main

import "github.com/google/gousb"
import "errors"
import "fmt"

var BufferSizes = []int {24, 60}

// TODO: Add more verbose error messages within functions to include the package
//       and function name. Try to use reflection to get the name of the package
//       and function

// -----------------------------------------------------------------------------
// Get the string representation of the vendor ID from the device descriptor.
// -----------------------------------------------------------------------------

func getVendorID(dev *gousb.Device) (string) {
	return dev.Desc.Vendor.String()
}

// -----------------------------------------------------------------------------
// Get the string representation of the product ID from the device descriptor.
// -----------------------------------------------------------------------------

func getProductID(dev *gousb.Device) (string) {
	return dev.Desc.Product.String()
}

// -----------------------------------------------------------------------------
// Get the software ID from device NVRAM using vendor control transfer.
// -----------------------------------------------------------------------------

func getDeviceSoftwareID(dev *gousb.Device, bufsize int) (string, error) {
	return getDeviceProperty(dev, PropertySoftwareID, bufsize)
}

// -----------------------------------------------------------------------------
// Get the serial number from device NVRAM using vendor control transfer.
// -----------------------------------------------------------------------------

func getDeviceSerialNumber(dev *gousb.Device, bufsize int) (string, error) {
	return getDeviceProperty(dev, PropertySerialNum, bufsize)
}

// -----------------------------------------------------------------------------
// Set the serial number in device NVRAM using vendor control transfer.
// -----------------------------------------------------------------------------

func setDeviceSerialNumber(dev *gousb.Device, bufsize int, prop string) (error) {
	return setDeviceProperty(dev, PropertySerialNum, bufsize, prop)
}

// -----------------------------------------------------------------------------
// Erase the serial number from device NVRAM using vendor control transfer.
// -----------------------------------------------------------------------------

func eraseDeviceSerialNumber(dev *gousb.Device, bufsize int) (error) {
	return setDeviceProperty(dev, PropertySerialNum, bufsize, "")
}

// -----------------------------------------------------------------------------
// Get the manufacturer name of the device from the device descriptor.
// -----------------------------------------------------------------------------

func getManufacturerName(dev *gousb.Device) (string, error) {

	var prop string

	desc, err := getDeviceDescriptor(dev)

	if err != nil {
		return prop, err
	}

	prop, err = dev.GetStringDescriptor(int(desc.ManufacturerIndex))

	if err != nil {
		return prop, err
	}

	return prop, nil
}

// -----------------------------------------------------------------------------
// Get the product name of the device from the device descriptor.
// -----------------------------------------------------------------------------

func getProductName(dev *gousb.Device) (string, error) {

	var prop string

	desc, err := getDeviceDescriptor(dev)

	if err != nil {
		return prop, err
	}

	prop, err = dev.GetStringDescriptor(int(desc.ProductIndex))

	if err != nil {
		return prop, err
	}

	return prop, nil
}

// -----------------------------------------------------------------------------
// Get the serial number of the device from the device descriptor. Changes made
// to the serial number on the device using a control transfer are not reflected
// in the device descriptor until the device is power-cycled (unplugged). The 
// most current information is always stored on the device.
// -----------------------------------------------------------------------------

func getSerialNumber(dev *gousb.Device) (string, error) {

	var prop string

	desc, err := getDeviceDescriptor(dev)

	if err != nil {
		return prop, err
	}

	prop, err = dev.GetStringDescriptor(int(desc.SerialNumberIndex))

	if err != nil {
		return prop, err
	}

	return prop, nil
}

// -----------------------------------------------------------------------------
// Get the device descriptor of the device.
// -----------------------------------------------------------------------------

func getDeviceDescriptor(dev *gousb.Device) (DeviceDescriptor, error) {

	var desc DeviceDescriptor

	data := make([]byte, BufferSizeDeviceDescriptor)

	_, err := dev.Control(
		RequestTypeStandardDeviceIn,
		RequestGetDescriptor,
		TypeDeviceDescriptor,
		InterfaceNumber,
		data)

	if err != nil {
		return desc, err
	}

	desc = DeviceDescriptor {
		Length: data[0],
		DescriptorType: data[1],
		UsbSpecification: uint16(data[2]) + (uint16(data[3]) << 8),
		DeviceClass: data[4],
		DeviceSubClass: data[5],
		DeviceProtocol: data[6],
		MaxPacketSize: data[7],
		VendorID: uint16(data[8]) + (uint16(data[9]) << 8),
		ProductID: uint16(data[10]) + (uint16(data[11]) << 8),
		DeviceReleaseNumber: uint16(data[12]) + (uint16(data[13]) << 8),
		ManufacturerIndex: data[14],
		ProductIndex: data[15],
		SerialNumberIndex: data[16],
		NumConfigurations: data[17]}

	return desc, nil
}

// -----------------------------------------------------------------------------
// Get the configuration descriptor of the device.
// -----------------------------------------------------------------------------

func getConfigDescriptor(dev *gousb.Device) (ConfigDescriptor, error) {

	var desc ConfigDescriptor

	data := make([]byte, BufferSizeConfigDescriptor)

	_, err := dev.Control(
		RequestTypeStandardDeviceIn,
		RequestGetDescriptor,
		TypeConfigDescriptor,
		InterfaceNumber,
		data)

	if err != nil {
		return desc, err
	}

	desc = ConfigDescriptor {
		Length: data[0],
		DescriptorType: data[1],
		TotalLength: uint16(data[2]) + (uint16(data[3]) << 8),
		NumInterfaces: data[4],
		ConfigurationValue: data[5],
		ConfigurationIndex: data[6],
		Attributes: data[7],
		MaxPower: data[8]}

	return desc, nil
}

// -----------------------------------------------------------------------------
// Use trial and error to find the control transfer data buffer size.
// -----------------------------------------------------------------------------

func findDeviceBufferSize(dev *gousb.Device) (int, error) {

	for _, size := range BufferSizes {

		data := make([]byte, size)
		copy(data, []byte{0x00, 0x01, 0x00})
		rc, _ := dev.Control(0x21, 0x09, 0x0300, 0x0000, data)

		if rc == size {
			return size, nil
		}
	}

	return 0, fmt.Errorf("Unsupported device")
}

func setDeviceProperty(dev *gousb.Device, id uint8, bufsize int, prop string) (error) {

	if len(prop) > bufsize - 3 {
		return errors.New("setProperty() error: Property length > data buffer")
	}

	data := make([]byte, bufsize)
	copy(data[0:], []byte{CommandSetProperty, uint8(len(prop)+1), id})
	copy(data[3:], prop)

	_, err := dev.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return err
	}

	data = make([]byte, bufsize)

	_, err = dev.Control(
		RequestTypeVendorDeviceIn,
		RequestGetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return err
	}

	if data[0] > 0x00 {
		return fmt.Errorf("Vendor command error: return code %d", int(data[0]))
	}

	// err = dev.Reset()

	return err
}

func getDeviceProperty(dev *gousb.Device, id uint8, bufsize int) (string, error) {

	var prop string

	data := make([]byte, bufsize)
	copy(data[0:], []byte{CommandGetProperty, 0x01, id})

	_, err := dev.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return prop, err
	}

	data = make([]byte, bufsize)

	_, err = dev.Control(
		RequestTypeVendorDeviceIn,
		RequestGetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return prop, err
	}

	if data[0] > 0x00 {
		return prop, fmt.Errorf("Vendor command error: return code %d", int(data[0]))
	}

	if data[1] > 0x00 {
		prop = string(data[2:int(data[1])+2])
	}

	return prop, nil
}

func resetDevice(dev *gousb.Device, bufsize int) (int, error) {

	data := make([]byte, bufsize)
	data[0] = 0x02

	return dev.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)
}
