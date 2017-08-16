package gomagtek

import "path/filepath"
import "runtime"
import "fmt"

// ==============
// Exported Names
// ==============

/*
 * Get the string representation of the vendor ID from the device descriptor.
 */
func (d *Device) GetVendorID() (string) {
	return d.Desc.Vendor.String()
}

/*
 * Get the string representation of the product ID from the device descriptor.
 */
func (d *Device) GetProductID() (string) {
	return d.Desc.Product.String()
}

/*
 * Get the software ID from device NVRAM using vendor control transfer.
 */
func (d *Device) GetSoftwareID() (string, error) {
	return d.getProperty(PropertySoftwareID)
}

/*
 * Get the serial number from device NVRAM using vendor control transfer.
 */
func (d *Device) GetSerialNumber() (string, error) {
	return d.getProperty(PropertySerialNum)
}

/*
 * Set the serial number in device NVRAM using vendor control transfer.
 */
func (d *Device) SetSerialNumber(prop string) (error) {
	return d.setProperty(PropertySerialNum, prop)
}

/*
 * Erase the serial number from device NVRAM using vendor control transfer.
 */
func (d *Device) EraseSerialNumber() (error) {
	return d.setProperty(PropertySerialNum, "")
}

/*
 * Get the manufacturer name of the device from the device descriptor.
 */
func (d *Device) GetManufacturerName() (string, error) {

	var prop string

	desc := new(DeviceDescriptor)
	err := desc.Get(d)

	if err != nil {
		return prop, err // Get method generates error message
	}

	if desc.ManufacturerIndex > 0 {
		prop, err = d.GetStringDescriptor(int(desc.ManufacturerIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

/*
 * Get the product name of the device from the device descriptor.
 */
func (d *Device) GetProductName() (string, error) {

	var prop string

	desc := new(DeviceDescriptor)
	err := desc.Get(d)

	if err != nil {
		return prop, err // Get method generates error message
	}

	if desc.ProductIndex > 0 {
		prop, err = d.GetStringDescriptor(int(desc.ProductIndex))
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
func (d *Device) GetSerialNumberDesc() (string, error) {

	var prop string

	desc := new(DeviceDescriptor)
	err := desc.Get(d)

	if err != nil {
		return prop, err // Get method generates error message
	}

	if desc.SerialNumberIndex > 0 {
		prop, err = d.GetStringDescriptor(int(desc.SerialNumberIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

/*
 * Reset the device using vendor commands in control transfer.
 */
func (d *Device) VendorReset() (int, error) {

	data := make([]byte, d.BufferSize)
	data[0] = 0x02

	rc, err := d.Control(
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
 * Use trial and error to find the control transfer data buffer size.
 * Failure to use the correct size for control transfers carrying
 * vendor commands will result in a LIBUSB_ERROR_PIPE error.
 */
func (d *Device) findBufferSize() (error) {

	var err error
	var rc, size int

	for _, size = range vendorBufferSizes {

		data := make([]byte, size)
		copy(data, []byte{0x00, 0x01, 0x00})
		rc, err = d.Control(0x21, 0x09, 0x0300, 0x0000, data)

		if rc == size {
			d.BufferSize = size
		}
	}

	if err != nil {
		err = fmt.Errorf("%s: unsupported device", getFunctionInfo())
	}

	return err
}

/*
 * Get a property from the device NVRAM using vendor commands in control transfer.
 */
func (d *Device) getProperty(id uint8) (string, error) {

	var prop string

	data := make([]byte, d.BufferSize)
	copy(data[0:], []byte{CommandGetProperty, 0x01, id})

	_, err := d.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
		return prop, err
	}

	data = make([]byte, d.BufferSize)

	_, err = d.Control(
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
 * Set a property on the device NVRAM using vendor commands in control transfer.
 */
func (d *Device) setProperty(id uint8, prop string) (error) {

	if len(prop) > d.BufferSize - 3 {
		return fmt.Errorf("%s: property length > data buffer", getFunctionInfo())
	}

	data := make([]byte, d.BufferSize)
	copy(data[0:], []byte{CommandSetProperty, uint8(len(prop)+1), id})
	copy(data[3:], prop)

	_, err := d.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
		return err
	}

	data = make([]byte, d.BufferSize)

	_, err = d.Control(
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
