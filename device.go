package gomagtek

import (
	"github.com/google/gousb"
	"strconv"
	"math"
	"time"
	"fmt"
)

// Device represents a USB device. The Device struct Desc field contains all
// information about the device. It includes the raw device descriptor, the
// config descriptor of the active config, and the size of the data buffer
// required by the device for vendor commands sent via control transfer.
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

// GetBusNumber retrieves the USB bus number of the device.
func (d *Device) GetBusNumber() string {
	return strconv.Itoa(d.Desc.Bus)
}

// GetBusAddress retrieves address of the device on the USB bus.
func (d *Device) GetBusAddress() string {
	return strconv.Itoa(d.Desc.Address)
}

// GetDeviceSpeed retrieves the negotiated operating speed of the device.
func (d *Device) GetDeviceSpeed() string {
	return d.Desc.Speed.String()
}

// GetUSBSpec retrieves the USB specification release number of the device.
func (d *Device) GetUSBSpec() string {
	return d.Desc.Spec.String()
}

// GetDeviceVer retrieves the major/minor version number ofthe device.
func (d *Device) GetDeviceVer() string {
	return d.Desc.Device.String()
}

// GetVendorId retrieves the USB vendor ID of the device.
func (d *Device) GetVendorID() string {
	return d.Desc.Vendor.String()
}

// GetProductID retrieves the USB product ID of the device.
func (d *Device) GetProductID() string {
	return d.Desc.Product.String()
}

// GetUSBClass retrieves the USB class of the device.
func (d *Device) GetUSBClass() string {
	return d.Desc.Class.String()
}

// GetUSBSubclass retrieves the USB subclass of the device.
func (d *Device) GetUSBSubclass() string {
	return d.Desc.SubClass.String()
}

// GetUSBProtocol retrieves the USB protocol of the device.
func (d *Device) GetUSBProtocol() string {
	return d.Desc.Protocol.String()
}

// GetMaxPktSize retrieves the maximum size of the control transfer.
func (d *Device) GetMaxPktSize() string {
	return strconv.Itoa(d.Desc.MaxControlPacketSize)
}

// GetBufferSize retrieves the size of the device data buffer.
func (d *Device) GetBufferSize() (string, error) {
	return strconv.Itoa(d.BufferSize), nil
}

// GetSoftwareID retrieves the software ID from the device NVRAM.
func (d *Device) GetSoftwareID() (string, error) {
	return d.getProperty(PropSoftwareID)
}

// GetProductVer retrieves the MagneSafe version from device NVRAM.
func (d *Device) GetProductVer() (value string, err error) {
	value, err = d.getProperty(PropProductVer)
	if len(value) <= 1 {value = ""}
	return value, err
}

// GetDeviceSN retrieves the configurable serial number from device NVRAM.
func (d *Device) GetDeviceSN() (string, error) {
	return d.getProperty(PropDeviceSN)
}

// SetDeviceSN sets the configurable serial number in device NVRAM.
func (d *Device) SetDeviceSN(value string) (error) {
	return d.setProperty(PropDeviceSN, value)
}

// EraseDeviceSN removes the configurable serial number from device NVRAM.
func (d *Device) EraseDeviceSN() (error) {
	return d.setProperty(PropDeviceSN, "")
}

// GetFactorySN retrieves the factory serial number from device NVRAM.
func (d *Device) GetFactorySN() (value string, err error) {
	value, err = d.getProperty(PropFactorySN)
	if len(value) <= 1 {value = ""}
	return value, err
}

// SetFactorySN sets the factory serial number in device NVRAM. This command
// will fail with result code 07 if the serial number is already configured.
func (d *Device) SetFactorySN(value string) (error) {
	return d.setProperty(PropFactorySN, value)
}

// CopyFactorySN copies 'length' characters from the factory serial
// number to the configurable serial number in device NVRAM.
func (d *Device) CopyFactorySN(length int) (error) {

	fs, err := d.GetFactorySN()

	if err != nil {
		return fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	if len(fs) == 0 {
		return fmt.Errorf("%s: factory serial number not present", getFunctionInfo())
	}

	limit := int(math.Min(float64(length), float64(len(fs))))
	err = d.SetDeviceSN(fs[:limit])

	return err
}

// GetVendorName retrieves the manufacturer name from device descriptor.
func (d *Device) GetVendorName() (value string, err error) {

	if d.DeviceDescriptor.ManufacturerIndex > 0 {
		value, err = d.GetStringDescriptor(int(d.DeviceDescriptor.ManufacturerIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return value, err
}

// GetProductName retrieves the product name from device descriptor.
func (d *Device) GetProductName() (value string, err error) {

	if d.DeviceDescriptor.ProductIndex > 0 {
		value, err = d.GetStringDescriptor(int(d.DeviceDescriptor.ProductIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return value, err
}

// GetDescriptSN retrieves the serial number of the device from the
// device descriptor. Changes made to the serial number on the device using a
// control transfer are not reflected in the device descriptor until the device
// is power-cycled (unplugged). The most current information is always stored
// on the device.
func (d *Device) GetDescriptSN() (value string, err error) {

	if d.DeviceDescriptor.SerialNumIndex > 0 {
		value, err = d.GetStringDescriptor(int(d.DeviceDescriptor.SerialNumIndex))
	}

	if err != nil {
		err = fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	return value, err
}

// UsbReset performs a USB port reset to reinitialize the device.
func (d *Device) UsbReset() (err error) {
	return d.Reset()
}


// DeviceReset resets the device using low-level vendor commands.
func (d *Device) DeviceReset() (err error) {

	data := make([]byte, d.BufferSize)
	data[0] = CommandResetDevice

	_, err = d.Control(
		RequestDirectionOut + RequestTypeClass + RequestRecipientDevice,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		err = fmt.Errorf("%s: %v)", getFunctionInfo(), err)
	}

	data = make([]byte, d.BufferSize)

	_, err = d.Control(
		RequestDirectionIn + RequestTypeClass + RequestRecipientDevice,
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

	time.Sleep(5 * time.Second)

	return err
}

// getDeviceDescriptor retrieves the raw device descriptor.
func (d *Device) getDeviceDescriptor() (err error) {

	data := make([]byte, BufferSizeDeviceDescriptor)

	_, err = d.Control(
		RequestDirectionIn + RequestTypeStandard + RequestRecipientDevice,
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

// getConfigDescriptor retrieves the raw active config descriptor.
func (d *Device) getConfigDescriptor() (err error) {

	data := make([]byte, BufferSizeConfigDescriptor)

	_, err = d.Control(
		RequestDirectionIn + RequestTypeStandard + RequestRecipientDevice,
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

// findBufferSize uses trial and error to find the control transfer data
// buffer size of the device. Failure to use the correct size for control
// transfers carrying vendor commands will result in a LIBUSB_ERROR_PIPE
// error.
func (d *Device) findBufferSize() (err error) {

	var rc, size int

	for _, size = range vendorBufferSizes {

		data := make([]byte, size)
		copy(data, []byte{CommandGetProp, 0x01, PropSoftwareID})

		rc, err = d.Control(
			RequestDirectionOut + RequestTypeClass + RequestRecipientDevice,
			RequestSetReport,
			TypeFeatureReport,
			InterfaceNumber,
			data)

		data = make([]byte, size)

		rc, err = d.Control(
			RequestDirectionIn + RequestTypeClass + RequestRecipientDevice,
			RequestGetReport,
			TypeFeatureReport,
			InterfaceNumber,
			data)

		if rc == size {
			d.BufferSize = size
			break
		}
	}

	if err != nil {
		err = fmt.Errorf("%s: unsupported device", getFunctionInfo())
	}

	return err
}

// getProperty retrieves a property from device NVRAM using low-level commands.
func (d *Device) getProperty(id uint8) (value string, err error) {

	data := make([]byte, d.BufferSize)
	copy(data, []byte{CommandGetProp, 0x01, id})

	_, err = d.Control(
		RequestDirectionOut + RequestTypeClass + RequestRecipientDevice,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return value, fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	data = make([]byte, d.BufferSize)

	_, err = d.Control(
		RequestDirectionIn + RequestTypeClass + RequestRecipientDevice,
		RequestGetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return value, fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	if data[0] > 0x00 {
		return value, fmt.Errorf("%s: command error: %d",
			getFunctionInfo(), int(data[0]))
	}

	if data[1] > 0x00 {
		value = string(data[2:int(data[1])+2])
	}

	return value, err
}

// setProperty configures a property in device NVRAM using low-level commands.
func (d *Device) setProperty(id uint8, value string) (err error) {

	if len(value) > d.BufferSize - 3 {
		return fmt.Errorf("%s: property length > data buffer", getFunctionInfo())
	}

	data := make([]byte, d.BufferSize)
	copy(data[0:], []byte{CommandSetProp, uint8(len(value)+1), id})
	copy(data[3:], value)

	_, err = d.Control(
		RequestDirectionOut + RequestTypeClass + RequestRecipientDevice,
		RequestSetReport,
		TypeFeatureReport,
		InterfaceNumber,
		data)

	if err != nil {
		return fmt.Errorf("%s: %v", getFunctionInfo(), err)
	}

	data = make([]byte, d.BufferSize)

	_, err = d.Control(
		RequestDirectionIn + RequestTypeClass + RequestRecipientDevice,
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
