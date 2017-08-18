package gomagtek

const (
	MagtekVendorID uint16 = 0x0801

	SureswipeKbPID uint16 = 0x0001
	SureswipeHidPID uint16 = 0x0002

	MagnesafeSwipeKbPID uint16 = 0x0001
	MagnesafeInsertKbPID uint16 = 0x0001
	MagnesafeSwipeHidPID uint16 = 0x0011
	MagnesafeInsertHidPID uint16 = 0x0013
	MagnesafeWirelessHidPID uint16 = 0x0014

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

	BufferSizeDeviceDescriptor int = 18
	BufferSizeConfigDescriptor int = 9
	BufferSizeSureswipe int = 24
	BufferSizeMagnesafe int = 60

	CommandGetProp uint8 = 0x00
	CommandSetProp uint8 = 0x01
	CommandResetDevice uint8 = 0x02

	ResultCodeSuccess uint8 = 0x00
	ResultCodeFailure uint8 = 0x01
	ResultCodeBadParam uint8 = 0x02

	PropSoftwareID uint8 = 0x00
	PropSerialNum uint8 = 0x01
	PropFactorySerialNum uint8 = 0x03
	PropMagnesafeVer uint8 = 0x04

	DefaultSerialNumLength int = 7
)

var (
	vendorBufferSizes = []int {24, 60}
)
