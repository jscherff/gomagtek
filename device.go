package gomagtek

import "github.com/google/gousb"
import "math"
import "fmt"

// ============================================================================
// Device Object.
// ============================================================================

/*
 * The gomagtek Device struct represents a USB device. The Device struct
 * Desc field contains all information about the device from gousb.Device.
 * The gomagtek Device extends the gousb Device by adding the raw device
 * descriptor, the config descriptor of the active config, and the size
 * of the data buffer required by the device for vendor commands sent via
 * control transfer.
 */
type Device struct {
	*gousb.Device
	DeviceDescriptor *DeviceDescriptor
	ConfigDescriptor *ConfigDescriptor
	BufferSize int
}

/*
 * Construct a new gomagtek Device from a gousb Device.
 */
func NewDevice(d *gousb.Device) (nd *Device, err error) {

	nd = &Device{d, new(DeviceDescriptor), new(ConfigDescriptor), 0}

	err = nd.getDeviceDescriptor()
	err = nd.getConfigDescriptor()
	err = nd.findBufferSize()

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return nd, err
}

// ============================================================================
// Public Methods.
// ============================================================================

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
 * Get the size of the device data buffer for vendor commands.
 */
func (d *Device) GetBufferSize() (int) {
	return d.BufferSize
}

/*
 * Get the software ID from device NVRAM using vendor control transfer.
 */
func (d *Device) GetSoftwareID() (string, error) {
	return d.getProperty(PropertySoftwareID)
}

/*
 * Get the MagneSafe version from device NVRAM using vendor control transfer.
 */
func (d *Device) GetMagnesafeVersion() (string, error) {
	return d.getProperty(PropertyMagnesafeVersion)
}

/*
 * Get the USB serial number from device NVRAM using vendor control transfer.
 */
func (d *Device) GetSerialNum() (string, error) {
	return d.getProperty(PropertySerialNum)
}

/*
 * Set the USB serial number in device NVRAM using vendor control transfer.
 */
func (d *Device) SetSerialNum(prop string) (error) {
	return d.setProperty(PropertySerialNum, prop)
}

/*
 * Erase the USB serial number from device NVRAM using vendor control transfer.
 */
func (d *Device) EraseSerialNum() (error) {
	return d.setProperty(PropertySerialNum, "")
}

/*
 * Get the device serial number from device NVRAM using vendor control transfer.
 */
func (d *Device) GetFactorySerialNum() (string, error) {
	return d.getProperty(PropertyFactorySerialNum)
}

/*
 * Set the device serial number in device NVRAM using vendor control transfer.
 * On Dynamag devices, this command will fail with result code 07 if the serial
 * number has already been configured.
 */
func (d *Device) SetFactorySerialNum(prop string) (error) {
	return d.setProperty(PropertyFactorySerialNum, prop)
}

/*
 * Erase the serial number from device NVRAM using vendor control transfer.
 * On Dynamag devices, this command will fail with result code 07 if the serial
 * number has already been configured.
 */
func (d *Device) EraseFactorySerialNum() (error) {
	return d.setProperty(PropertyFactorySerialNum, "")
}

/*
 * Copy 'length' characters from the device serial number to the USB serial
 * number in device NVRAM using vendor control transfer.
 */
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

/*
 * Get the manufacturer name of the device from the device descriptor.
 */
func (d *Device) GetManufacturerName() (prop string, err error) {

	if d.DeviceDescriptor.ManufacturerIndex > 0 {
		prop, err = d.GetStringDescriptor(int(d.DeviceDescriptor.ManufacturerIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

/*
 * Get the product name of the device from the device descriptor.
 */
func (d *Device) GetProductName() (prop string, err error) {

	if d.DeviceDescriptor.ProductIndex > 0 {
		prop, err = d.GetStringDescriptor(int(d.DeviceDescriptor.ProductIndex))
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
func (d *Device) GetDescriptorSerialNum() (prop string, err error) {

	if d.DeviceDescriptor.SerialNumIndex > 0 {
		prop, err = d.GetStringDescriptor(int(d.DeviceDescriptor.SerialNumIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return prop, err
}

/*
 * Reset the device using vendor commands in control transfer.
 */
func (d *Device) VendorReset() (rc int, err error) {

	data := make([]byte, d.BufferSize)
	data[0] = 0x02

	rc, err = d.Control(
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

// ============================================================================
// Private Methods.
// ============================================================================

/*
 * Get the device descriptor of the device.
 */
func (d *Device) getDeviceDescriptor() (err error) {

	data := make([]byte, BufferSizeDeviceDescriptor)

	_, err = d.Control(
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
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return err
}

/*
 * Get the config descriptor of the device.
 */
func (d *Device) getConfigDescriptor() (err error) {

	data = make([]byte, BufferSizeConfigDescriptor)

	_, err = d.Control(
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
		return fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return err
}

/*
 * Use trial and error to find the control transfer data buffer size.
 * Failure to use the correct size for control transfers carrying
 * vendor commands will result in a LIBUSB_ERROR_PIPE error.
 */
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

/*
 * Get a property from the device NVRAM using vendor commands in control transfer.
 */
func (d *Device) getProperty(id uint8) (prop string, err error) {

	data := make([]byte, d.BufferSize)
	copy(data[0:], []byte{CommandGetProperty, 0x01, id})

	_, err = d.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return prop, fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	data = make([]byte, d.BufferSize)

	_, err = d.Control(
		RequestTypeVendorDeviceIn,
		RequestGetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return prop, fmt.Errorf("%s: %v", getFunctionInfo(), err)
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

/*
 * Set a property on the device NVRAM using vendor commands in control transfer.
 */
func (d *Device) setProperty(id uint8, prop string) (err error) {

	if len(prop) > d.BufferSize - 3 {
		return fmt.Errorf("%s: property length > data buffer", getFunctionInfo())
	}

	data := make([]byte, d.BufferSize)
	copy(data[0:], []byte{CommandSetProperty, uint8(len(prop)+1), id})
	copy(data[3:], prop)

	_, err = d.Control(
		RequestTypeVendorDeviceOut,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	data = make([]byte, d.BufferSize)

	_, err = d.Control(
		RequestTypeVendorDeviceIn,
		RequestGetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	if data[0] > 0x00 {
		err = fmt.Errorf("%s: command error: %d",
			getFunctionInfo(), int(data[0]))
	}

	return err
}
