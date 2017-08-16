package main

import "github.com/google/gousb"
import "github.com/jscherff/gomagtek"
import "log"
import "fmt"
import "os"

var printFormat string = "\tVendor ID:\t%s\n\tProduct ID:\t%s\n\t" +
	"Software ID:\t%s\n\tSerial Num:\t%s\n\tHost Name:\t%s\n\n"

func main() {

	context := gousb.NewContext()
	defer context.Close()

	devices, _ := context.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == gousb.ID(gomagtek.MagtekVendorID)
	})

	if len(devices) == 0 {
		log.Fatalf("No Magtek devices found")
	}

	for _, device := range devices {

		magtek, err := gomagtek.NewDevice(device)

		defer magtek.Close()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		vendorID := magtek.GetVendorID()
		productID := magtek.GetProductID()

		hostName, err := os.Hostname()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		softwareID, err := magtek.GetSoftwareID()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		serialNum, err := magtek.GetSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("BEFORE\n" + printFormat, vendorID, productID,
			softwareID, serialNum, hostName)

		err = magtek.EraseSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		serialNum, err = magtek.GetSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("AFTER\n" + printFormat, vendorID, productID,
			softwareID, serialNum, hostName)
	}
}
