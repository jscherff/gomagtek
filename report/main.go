package main

import (
	"github.com/jscherff/gomagtek"
	"github.com/google/gousb"
	"flag"
	"log"
	"fmt"
	"os"
)

var (
	fsMode = flag.NewFlagSet("mode", flag.ExitOnError)
	fReportMode = fsMode.Bool("report", false, "Report mode")
	fConfigMode = fsMode.Bool("config", false, "Config mode")
)

var (
	fsReportOptions = flag.NewFlagSet("report_fields", flag.ExitOnError)
	fHostname = fsReportOptions.Bool("hostname", true, "Include host name")
	fVendorId = fsReportOptions.Bool("vendor_id", false, "Include device vendor ID")
	fVendorName = fsReportOptions.Bool("vendor_name", false, "Include device vendor name")
	fProductId = fsReportOptions.Bool("product_id", false, "Include device product ID")
	fProductName = fsReportOptions.Bool("product_name", false, "Include device product name")
	fProductVer = fsReportOptions.Bool("product_ver", false, "Include device product version")
	fSoftwareId = fsReportOptions.Bool("software_id", false, "Include device software")
	fBufferSize = fsReportOptions.Bool("buffer_size", false, "Include device buffer size")
	fDeviceSerial = fsReportOptions.Bool("device_serial", true, "Include device serial number")
	fFactorySerial = fsReportOptions.Bool("factory_serial", false, "Include device factory serial number")
	fDescripSerial = fsReportOptions.Bool("descriptor_serial", false, "Include descriptor serial number")
	fOutputFormatCSV = fsReportOptions.Bool("output_format_csv", true, "Write output in CSV format")
	fOutputFormatNVP = fsReportOptions.Bool("output_format_nvp", false, "Write output as name/value pairs")
	fOutputScreen = fsReportOptions.Bool("output_screen", false, "Write output to screen")
	fOutputFile = fsReportOptions.String ("output_file", "", "Write output to `file`")
)

var (
	fsConfigOptions = flag.NewFlagSet("config_options", flag.ExitOnError)
	fEraseSerial = fsConfigOptions.Bool("erase_device_serial", false, "Erase device serial number")
	fConfigSerial = fsConfigOptions.String("config_device_serial", "", "Set device serial number to `string`")
	fCopySerial = fsConfigOptions.Int("copy_factory_serial", 7, "Copy `n` bytes from factory to device serial number")
	fPerformUsbReset = fsConfigOptions.Bool("perform_usb_reset", false, "Perform a USB reset")
	fPerformDevReset = fsConfigOptions.Bool("perform_dev_reset", false, "Perform a device reset")
)

var printFormat string = "\tVendor ID:\t%s\n\tProduct ID:\t%s\n\t" +
	"Software ID:\t%s\n\tSerial Num:\t%s\n\tHost Name:\t%s\n\n"

func main() {

	fsMode.Parse(os.Args[1:2])

	if fsMode.NFlag() != 1 {
		fmt.Fprintf(os.Stderr, "You must specify one and only one option\n\n")
	}

	switch {
		case *fReportMode:
			fsReportOptions.Parse(os.Args[2:])
		case *fConfigMode:
			fsConfigOptions.Parse(os.Args[2:])
	}


	//fmt.Println(flag.NFlag())
	//fmt.Println(*flagProductId)
	//fmt.Println(*flagVendorId)
	//flag.PrintDefaults()

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
