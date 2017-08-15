package gomagtek

import "github.com/google/gousb"
import "path/filepath"
import "runtime"
import "errors"
import "fmt"

// ==============
// Exported Names
// ==============

var BufferSizes = []int {24, 60}

/*
 * Get the string representation of the vendor ID from the device descriptor.
 */
func GetVendorID(dev *gousb.Device) (string) {
	return dev.Desc.Vendor.String()
}

/*
 * Get the string representation of the product ID from the device descriptor.
 */
func GetProductID(dev *gousb.Device) (string) {
	return dev.Desc.Product.String()
}

/*
 * Get the software ID from device NVRAM using vendor control transfer.
 */
func GetDeviceSoftwareID(dev *gousb.Device, bufsize int) (string, error) {
	return getDeviceProperty(dev, PropertySoftwareID, bufsize)
}

/*
 * Get the serial number from device NVRAM using vendor control transfer.
 */
func GetDeviceSerialNumber(dev *gousb.Device, bufsize int) (string, error) {
	return getDeviceProperty(dev, PropertySerialNum, bufsize)
}

/*
 * Set the serial number in device NVRAM using vendor control transfer.
 */
func SetDeviceSerialNumber(dev *gousb.Device, bufsize int, prop string) (error) {
	return setDeviceProperty(dev, PropertySerialNum, bufsize, prop)
}

/*
 * Erase the serial number from device NVRAM using vendor control transfer.
 */
func EraseDeviceSerialNumber(dev *gousb.Device, bufsize int) (error) {
	return setDeviceProperty(dev, PropertySerialNum, bufsize, "")
}

/*
 * Get the manufacturer name of the device from the device descriptor.
 */
func GetManufacturerName(dev *gousb.Device) (string, error) {

	var prop string

	desc := new(DeviceDescriptor)
	err := desc.Get(dev)

	if err != nil {
		return prop, err // Get method generates error message
	}

	if desc.ManufacturerIndex > 0 {
		prop, err = dev.GetStringDescriptor(int(desc.ManufacturerIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

/*
 * Get the product name of the device from the device descriptor.
 */
func GetProductName(dev *gousb.Device) (string, error) {

	var prop string

	desc := new(DeviceDescriptor)
	err := desc.Get(dev)

	if err != nil {
		return prop, err // Get method generates error message
	}

	if desc.ProductIndex > 0 {
		prop, err = dev.GetStringDescriptor(int(desc.ProductIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

/*
 * Get the serial number of the device from the device descriptor. Changes made
 * to the serial number on the device using a control transfer are not reflected
 * in the device descriptor until the device is power-cycled (unplugged). The 
 * most current information is always stored on the device.
 */
func GetSerialNumber(dev *gousb.Device) (string, error) {

	var prop string

	desc := new(DeviceDescriptor)
	err := desc.Get(dev)

	if err != nil {
		return prop, err // Get method generates error message
	}

	if desc.SerialNumberIndex > 0 {
		prop, err = dev.GetStringDescriptor(int(desc.SerialNumberIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

/*
 * Use trial and error to find the control transfer data buffer size.
 */
func FindDeviceBufferSize(dev *gousb.Device) (int, error) {

	var rc int
	var err error

	for _, size := range BufferSizes {

		data := make([]byte, size)
		copy(data, []byte{0x00, 0x01, 0x00})
		rc, err = dev.Control(0x21, 0x09, 0x0300, 0x0000, data)

		// If command is successful, the return code will be a
		// non-zero positive integer reflecting the buffer size.

		if rc == size {
			return size, nil
		}
	}

	return 0, fmt.Errorf("%s: unsupported device: %v)", getFunctionInfo(), err)
}

// ================
// Unexported Names
// ================

/*
 * Get function filename, line number, and function name for error reporting.
 */
func getFunctionInfo() string {
	pc, file, line, success := runtime.Caller(1)
	function := runtime.FuncForPC(pc)
	if !success {
		return "Unknown goroutine"
	}
	return fmt.Sprintf("%s:%d: %s()", filepath.Base(file), line, function.Name())
}

/*
 * Set a property on the device NVRAM using vendor commands in control transfer.
 */
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
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
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
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
		return err
	}

	if data[0] > 0x00 {
		err = fmt.Errorf("%s: Vendor command error: return code %d",
			getFunctionInfo(), int(data[0]))
	}

	return err
}

/*
 * Get a property from the device NVRAM using vendor commands in control transfer.
 */
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
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
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
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
		return prop, err
	}

	if data[0] > 0x00 {
		err = fmt.Errorf("%s: Vendor command error: return code %d",
			getFunctionInfo(), int(data[0]))
	}

	if data[1] > 0x00 {
		prop = string(data[2:int(data[1])+2])
	}

	return prop, err
}

/*
 * Reset the device using vendor commands in control transfer.
 */
func resetDevice(dev *gousb.Device, bufsize int) (int, error) {

	data := make([]byte, bufsize)
	data[0] = 0x02

	rc, err := dev.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return rc, err
}
