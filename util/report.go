package main

import(
	"strings"
	"fmt"
)

type ReportItem struct {
	Name string
	Value string
}

type Report []ReportItem

func report() {
	fields := strings.Split(",", *fReportInc)
	for field := range fields {

		switch field {
		case "hn":
			fmt.Println("Host Name")
			fallthrough
		case "vid":
			fmt.Println("Vendor ID")
			fallthrough
		case "vn":
			fmt.Println("Vendor Name")
			fallthrough
		case "pid":
			fmt.Println("Product ID")
			fallthrough
		case "pn":
			fmt.Println("Product Name")
			fallthrough
		case "pv":
			fmt.Println("Product Version")
			fallthrough
		case "sid":
			fmt.Println("Software ID")
			fallthrough
		case "bs":
			fmt.Println("Buffer Size")
			fallthrough
		case "sn":
			fmt.Println("Device Serial Number")
			fallthrough
		case "fsn":
			fmt.Println("Factory Serial Number")
			fallthrough
		case "dsn":
			fmt.Println("Descriptor Serial Number")
		default:
			fmt.Println("ERROR")
		}
	}
}
