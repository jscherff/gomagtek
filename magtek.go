package main

import "github.com/google/gousb"
import "github.com/jscherff/gomagtek"
import "log"
import "fmt"
import "os"

func main() {

	context := gousb.NewContext()
	defer context.Close()

	// Open devices that report a Magtek vendor ID, 0x0801.
	// We omit error checking on OpenDevices() because this
	// function terminates with 'libusb: not found [code -5]'
	// on Windows systems.

	devices, _ := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == gousb.ID(MagtekVendorID)
	})

	if len(devices) == 0 {
		log.Fatalf("No Magtek devices found")
	}

	for _, device := range devices {

		defer device.Close()

		// Determine the data buffer size for vendor commands.
		// Failure to use the correct size value for control
		// transfers carrying vendor commands will result in
		// a LIBUSB_ERROR_PIPE error.

		bufSize, err := findDeviceBufferSize(device)

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		// Information from the Device Descriptor

		vendorID := getVendorID(device)
		productID := getProductID(device)

		// Information obtained from the device using control
		// transfer. Serial number can also be obtained using
		// the devices getStringDescriptor command with the
		// serial number index obtained from the device des-
		// criptor; however, this value is not refreshed until
		// the device is power-cycled.

		softwareID, _ := getDeviceSoftwareID(device, bufSize)

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		serialNum, err := getDeviceSerialNumber(device, bufSize)

		// Host name as reported by the operating system.

		hostName, _ := os.Hostname()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("Before: %s,%s,%s,%s,%s\n", vendorID, productID,
			softwareID, serialNum, hostName)

		if len(serialNum) == 0 {

			//TODO Phone Home routine to get serial number
			serialNum = "24FA12C"
			err = setDeviceSerialNumber(device, bufSize, serialNum)

			if err != nil {
				log.Fatalf("Error: %v", err); continue
			}

			serialNum, err = getDeviceSerialNumber(device, bufSize)

			if err != nil {
				log.Fatalf("Error: %v", err); continue
			}
		}


		fmt.Printf("After:  %s,%s,%s,%s,%s\n", vendorID, productID,
			softwareID, serialNum, hostName)
	}
}
