package main

import (
	"github.com/jscherff/gomagtek"
	"github.com/google/gousb"
	"flag"
	"log"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "You must specify a mode of operation.\n")
		fsMode.Usage()
		os.Exit(1)
	}

	fsMode.Parse(os.Args[1:2])

	var flagset *flag.FlagSet

	switch {

	case *fModeReport:
		flagset = fsReport

	case *fModeConfig:
		flagset = fsConfig

	case *fModeReset:
		flagset = fsReset
	}

	if flagset.Parse(os.Args[2:]); flagset.NFlag() == 0 {
		fmt.Fprintf(os.Stderr, "You must specify at least one option.\n")
		flagset.Usage()
		os.Exit(1)
	}

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

		di, errs := gomagtek.NewDeviceInfo(device)

		if len(errs) > 0 {
			log.Fatalf("Errors encountered"); continue
		}

		fmt.Println(di)
		fmt.Println(di.JSON())
		fmt.Println(di.XML())
		fmt.Println(di.FXML())
		os.Exit(0)

		switch {

		case *fModeReport:
			err = report(device)

		case *fModeConfig:
			err = config(device)

		case *fModeReset:
			err = reset(device)
		}
	}
}
