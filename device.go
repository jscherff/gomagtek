package gomagtek

import (
	"github.com/google/gousb"
	"math"
	"fmt"
)

// Device represents a USB device. The Device struct Desc field contains
// all information about the device. It includes the raw device descriptor,
// the config descriptor of the active config, and the size of the data
// buffer required by the device for vendor commands sent via control
// transfer.
type Device struct {
	*gousb.Device
	BufferSize int
	DeviceDescriptor *DeviceDescriptor
	ConfigDescriptor *ConfigDescriptor
}

// NewDevice constructs a new Device.
func NewDevice(d *gousb.Device) (nd *Device, err error) {

	nd = &Device{d, 0, new(DeviceDescriptor), new(ConfigDescriptor)}

	err = nd.findBufferSize()

	if err != nil {
		return nd, err
	}

	_ = nd.getDeviceDescriptor()
	_ = nd.getConfigDescriptor()

	return nd, err
}

// DeviceDesc is a partial representation of a USB device descriptor.
type DeviceDesc struct {
	*gousb.DeviceDesc
}

// GetVendorId retrieves the string representation of the vendor ID from
// the device descriptor.
func (d *Device) GetVendorID() (string) {
	return d.Desc.Vendor.String()
}

// GetProductID retrieves the string representation of the product ID from
// the device descriptor.
func (d *Device) GetProductID() (string) {
	return d.Desc.Product.String()
}

// GetBufferSize retrieves the size of the device data buffer.
func (d *Device) GetBufferSize() (int) {
	return d.BufferSize
}

// GetSoftwareID retrieves the software ID from the device NVRAM.
func (d *Device) GetSoftwareID() (string, error) {
	return d.getProperty(PropSoftwareID)
}

// GetProductVer retrieves the MagneSafe version from device NVRAM.
func (d *Device) GetProductVer() (string, error) {
	return d.getProperty(PropProductVer)
}

// GetSerialNum retrieves the configurable serial number from device NVRAM.
func (d *Device) GetSerialNum() (string, error) {
	return d.getProperty(PropSerialNum)
}

// SetSerialNum sets the configurable serial number in device NVRAM.
func (d *Device) SetSerialNum(prop string) (error) {
	return d.setProperty(PropSerialNum, prop)
}

// EraseSerialNum removes the configurable serial number from device NVRAM.
func (d *Device) EraseSerialNum() (error) {
	return d.setProperty(PropSerialNum, "")
}

// GetFactorySerialNum retrieves the factory serial number from device NVRAM.
func (d *Device) GetFactorySerialNum() (string, error) {
	return d.getProperty(PropFactorySerialNum)
}

// SetFactorySerialNum sets the factory serial number in device NVRAM. This
// command will fail with result code 07 if the serial number has already been
// configured.
func (d *Device) SetFactorySerialNum(prop string) (error) {
	return d.setProperty(PropFactorySerialNum, prop)
}

// CopyFactorySerialNum copies 'length' characters from the factory serial
// number to the configurable serial number in device NVRAM.
func (d *Device) CopyFactorySerialNum(length int) (error) {

	ds, err := d.GetFactorySerialNum()

	if err == nil {
		limit := int(math.Min(float64(length), float64(len(ds))))
		err = d.SetSerialNum(ds[:limit])
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return err
}

// GetVendorName retrieves the manufacturer name from device descriptor.
func (d *Device) GetVendorName() (prop string, err error) {

	if d.DeviceDescriptor.ManufacturerIndex > 0 {
		prop, err = d.GetStringDescriptor(int(d.DeviceDescriptor.ManufacturerIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

// GetProductName retrieves the product name from device descriptor.
func (d *Device) GetProductName() (prop string, err error) {

	if d.DeviceDescriptor.ProductIndex > 0 {
		prop, err = d.GetStringDescriptor(int(d.DeviceDescriptor.ProductIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

// GetDescriptorSerialNum retrieves the serial number of the device from the
// device descriptor. Changes made to the serial number on the device using a
// control transfer are not reflected in the device descriptor until the device
// is power-cycled (unplugged). The most current information is always stored
// on the device.
func (d *Device) GetDescriptorSerialNum() (prop string, err error) {

	if d.DeviceDescriptor.SerialNumIndex > 0 {
		prop, err = d.GetStringDescriptor(int(d.DeviceDescriptor.SerialNumIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

// UsbReset performs a USB port reset to reinitialize the device.
func (d *Device) UsbReset() (err error) {
	return d.Reset()
}


// DeviceReset resets the device using low-level vendor commands.
func (d *Device) DeviceReset() (err error) {

	data := make([]byte, d.BufferSize)
	data[0] = CommandResetDevice

	rc, err := d.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		err = fmt.Errorf("%s: %v (%d)", getFunctionInfo(), err, rc)
	}

	return err
}

// getDeviceDescriptor retrieves the raw device descriptor.
func (d *Device) getDeviceDescriptor() (err error) {

	data := make([]byte, BufferSizeDeviceDescriptor)

	rc, err := d.Control(
		RequestTypeStandardDeviceIn,
		RequestGetDescriptor,
		TypeDeviceDescriptor,
		InterfaceNumber,
		data)

	if err == nil {

		*d.DeviceDescriptor = DeviceDescriptor {
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
		err = fmt.Errorf("%s: %v (%d)", getFunctionInfo(), err, rc)
	}

	return err
}

// getConfigDescriptor retrieves the raw active config descriptor.
func (d *Device) getConfigDescriptor() (err error) {

	data := make([]byte, BufferSizeConfigDescriptor)

	rc, err := d.Control(
		RequestTypeStandardDeviceIn,
		RequestGetDescriptor,
		TypeConfigDescriptor,
		InterfaceNumber,
		data)

	if err == nil {

		*d.ConfigDescriptor = ConfigDescriptor {
			data[0],
			data[1],
			uint16(data[2]) + (uint16(data[3]) << 8),
			data[4],
			data[5],
			data[6],
			data[7],
			data[8]}
	} else {
		return fmt.Errorf("%s: %v (%d)", getFunctionInfo(), err, rc)
	}

	return err
}

// findBufferSize uses trial and error to find the control transfer data
// buffer size of the device. Failure to use the correct size for control
// transfers carrying vendor commands will result in a LIBUSB_ERROR_PIPE
// error.
func (d *Device) findBufferSize() (err error) {

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

// getProperty retrieves a property from device NVRAM using low-level commands.
func (d *Device) getProperty(id uint8) (prop string, err error) {

	data := make([]byte, d.BufferSize)
	copy(data[0:], []byte{CommandGetProp, 0x01, id})

	rc, err := d.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return prop, fmt.Errorf("%s: %v (%d)", getFunctionInfo(), err, rc)
	}

	data = make([]byte, d.BufferSize)

	rc, err = d.Control(
		RequestTypeVendorDeviceIn,
		RequestGetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return prop, fmt.Errorf("%s: %v (%d)", getFunctionInfo(), err, rc)
	}

	if data[0] > 0x00 {
		return prop, fmt.Errorf("%s: command error: %d",
			getFunctionInfo(), int(data[0]))
	}

	if data[1] > 0x00 {
		prop = string(data[2:int(data[1])+2])
	}

	return prop, err
}

// setProperty configures a property in device NVRAM using low-level commands.
func (d *Device) setProperty(id uint8, prop string) (err error) {

	if len(prop) > d.BufferSize - 3 {
		return fmt.Errorf("%s: property length > data buffer", getFunctionInfo())
	}

	data := make([]byte, d.BufferSize)
	copy(data[0:], []byte{CommandSetProp, uint8(len(prop)+1), id})
	copy(data[3:], prop)

	rc, err := d.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return fmt.Errorf("%s: %v (%d)", getFunctionInfo(), err, rc)
	}

	data = make([]byte, d.BufferSize)

	rc, err = d.Control(
		RequestTypeVendorDeviceIn,
		RequestGetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return fmt.Errorf("%s: %v (%d)", getFunctionInfo(), err, rc)
	}

	if data[0] > 0x00 {
		err = fmt.Errorf("%s: command error: %d",
			getFunctionInfo(), int(data[0]))
	}

	return err
}
