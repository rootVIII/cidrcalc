package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rootVIII/cidrcalc"
)

func main() {
	IPCIDR := flag.String("i", "", "IP address")
	flag.Parse()

	var net = &cidrcalc.Subnet{}
	err := net.Calculate(*IPCIDR)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Printf("IP uint32: 0%X\n", net.IPUINT32)
	fmt.Printf("CIDR: %d\n", net.CIDR)
	fmt.Printf("Broadcast Address uint32: 0x%08X\n", net.BroadcastAddressUINT32)
	fmt.Printf("Broadcast Address bytes: %v\n", net.BroadcastAddress)
	fmt.Printf("Network Address uint32: 0x%X\n", net.NetworkAddressUINT32)
	fmt.Printf("Network Address: %v\n", net.NetworkAddress)
	fmt.Printf("Subnet Mask: %v\n", net.SubnetMask)
	fmt.Printf("SubnetMask uint32: 0x%X\n", net.SubnetMaskUINT32)
	fmt.Printf("Subnet Bitmap: %q\n", net.SubnetBitmap)
	fmt.Printf("Number of Hosts: %d\n", net.HostsMAX)

}
