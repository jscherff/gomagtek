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


		factorySerialNum, err:= magtek.GetFactorySN()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("Device Serial Number: %s\n", factorySerialNum)


		serialNum, err:= magtek.GetDeviceSN()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("USB Serial Number: %s\n", serialNum)


		fmt.Println("Erasing serial number...")

		err = magtek.EraseDeviceSN()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}


		fmt.Println("Checking serial number...")

		serialNum, err = magtek.GetDeviceSN()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("USB Serial Number: %s\n", serialNum)


		fmt.Println("Copying serial number...")

		err = magtek.CopyFactorySN(gomagtek.DefaultSNLength)

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}


		fmt.Println("Checking serial number...")

		serialNum, err = magtek.GetDeviceSN()

		if err != nil {
			log.Fatalf("Error: %v", err); continue
		}

		fmt.Printf("USB Serial Number: %s\n", serialNum)
	}
}
