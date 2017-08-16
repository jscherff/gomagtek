package main

import "github.com/google/gousb"
import "github.com/jscherff/gomagtek"
import "log"
import "fmt"

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

		defer device.Close()

		magtek, err := gomagtek.NewDevice(device)

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}


		magnesafeVersion, err := magtek.GetMagnesafeVersion()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("MagneSafe Version: %s\n", magnesafeVersion)


		devSerialNum, err:= magtek.GetDevSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("Device Serial Number: %s\n", devSerialNum)


		usbSerialNum, err:= magtek.GetUsbSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("USB Serial Number: %s\n", usbSerialNum)


		fmt.Println("Erasing serial number...")

		err = magtek.EraseUsbSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}


		fmt.Println("Checking serial number...")

		usbSerialNum, err = magtek.GetUsbSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("USB Serial Number: %s\n", usbSerialNum)


		fmt.Println("Copying serial number...")

		err = magtek.CopyDevSerialNumber(7)

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}


		fmt.Println("Checking serial number...")

		usbSerialNum, err = magtek.GetUsbSerialNumber()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("USB Serial Number: %s\n", usbSerialNum)
	}
}
