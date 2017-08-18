package main

import "github.com/jscherff/gomagtek"
import "github.com/google/gousb"
import "log"
import "fmt"
import "os"

var printFormat string = "\tVendor ID:\t%s\n\tProduct ID:\t%s\n\t" +
	"Software ID:\t%s\n\tSerial Num:\t%s\n\tHost Name:\t%s\n\n"

func main() {

	context := gousb.NewContext()
	defer context.Close()

	// Open devices that report a Magtek vendor ID, 0x0801.
	// We omit error checking on OpenDevices() because this
	// function terminates with 'libusb: not found [code -5]'
	// on Windows systems.

	devices, _ := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return uint16(desc.Vendor) == gomagtek.MagtekVendorID
	})

	if len(devices) == 0 {
		log.Fatalf("No Magtek devices found")
	}

	for _, device := range devices {

		defer device.Close()

		device, err := gomagtek.NewDevice(device)

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		vendorID := device.GetVendorID()
		productID := device.GetProductID()

		// Information obtained from the device using control
		// transfer. Serial number can also be obtained using
		// the devices getStringDescriptor command with the
		// serial number index obtained from the device des-
		// criptor; however, this value is not refreshed until
		// the device is power-cycled.

		hostName, err := os.Hostname()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		softwareID, err := device.GetSoftwareID()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		serialNum, err := device.GetSerialNum()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("BEFORE\n" + printFormat, vendorID, productID,
			softwareID, serialNum, hostName)

		if len(serialNum) == 0 {

			serialNum = "24FA12C" //TODO: obtain from server
			err = device.SetSerialNum(serialNum)

			if err != nil {
				log.Fatalf("Error: %v", err); continue
			}

			serialNum, err = device.GetSerialNum()

			if err != nil {
				log.Fatalf("Error: %v", err); continue
			}
		}

		fmt.Printf("AFTER\n" + printFormat, vendorID, productID,
			softwareID, serialNum, hostName)
	}
}
