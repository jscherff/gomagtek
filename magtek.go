package main

import "github.com/google/gousb"
import "log"
import "fmt"
import "os"

func main() {

	context := gousb.NewContext()
	defer context.Close()

	// Open devices that report a Magtek vendor ID, 0x0801

	devices, err := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == gousb.ID(MagtekVendorID)
	})

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if len(devices) == 0 {
		log.Fatalf("No Magtek devices found")
	}

	for _, device := range devices {

		defer device.Close()

		bufSize, err := findDeviceBufferSize(device, BufferSizes)

		if err != nil {
			log.Fatalf("Error: %v", err)
			continue
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
		serialNum, _ := getDeviceSerialNumber(device, bufSize)

		// Host name as reported by the operating system.

		hostName, _ := os.Hostname()

		fmt.Printf("Before: %s,%s,%s,%s,%s\n", vendorID, productID,
			softwareID, serialNum, hostName)

		if len(serialNum) == 0 {

			//TODO Phone Home routine to get serial number
			err = setDeviceSerialNumber(device, bufSize, "24FA12C")

			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			serialNum, _ = getDeviceSerialNumber(device, bufSize)
		}


		fmt.Printf("After:  %s,%s,%s,%s,%s\n", vendorID, productID,
			softwareID, serialNum, hostName)
	}
}
