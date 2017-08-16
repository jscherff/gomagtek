package main

import "github.com/google/gousb"
import "../../gomagtek"
import "log"
import "fmt"
import "os"

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

		serialNumDesc, err := magtek.GetSerialNumberDesc()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("Before:\t%s,%s,%s,%s,(%s),%s\n", vendorID, productID,
			softwareID, serialNum, serialNumDesc, hostName)

		err = magtek.EraseSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		magtek.Reset()

		serialNum, err = magtek.GetSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		serialNumDesc, err = magtek.GetSerialNumberDesc()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("After:\t%s,%s,%s,%s,(%s),%s\n", vendorID, productID,
			softwareID, serialNum, serialNumDesc, hostName)

	}
}
