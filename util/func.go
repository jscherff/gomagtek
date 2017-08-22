package main

import (
	"github.com/jscherff/gomagtek"
	"strings"
	"fmt"
)

func reset(d *gomagtek.Device) (err error) {

	switch {

	case *fResetUsb:
		err = d.UsbReset()

	case *fResetDev:
		err = d.DeviceReset()
	}

	return err
}

func report(d *gomagtek.Device) (err error) {

	r, err := d.Report(strings.Split(*fReportInclude, ","))

	fmt.Println(r.CSV(false) + "\n")
	fmt.Println(r.CSV(true) + "\n")
	fmt.Println(r.NVP(false) + "\n")
	fmt.Println(r.NVP(true) + "\n")

	return err
}

func config(d *gomagtek.Device) (err error) {

	switch {

	case *fConfigErase:
		err = d.EraseDeviceSN()

	case *fConfigEmpty:
		if ln, err := d.GetDeviceSN(); len(ln) > 0 { return err }
	}

	switch {

	case *fConfigErase:
		err = d.EraseDeviceSN()

	case len(*fConfigSet) > 0:
		err = d.SetDeviceSN(*fConfigSet)

	case len(*fConfigUrl) > 0:
		//TODO: call function to fetch serail number from server
		err = d.SetDeviceSN("24F0000")

	case *fConfigCopy > 0:
		err = d.CopyFactorySN(*fConfigCopy)
	}

	return err
}
