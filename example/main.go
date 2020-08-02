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
		os.Exit(2)
	}

	fmt.Printf("CIDR: %d\n", net.CIDR)
	fmt.Printf("IP Address uint32: 0x%08X\n", net.IPUINT32)
	fmt.Printf("IP Address: %v\n", net.IP)
	fmt.Printf("Broadcast Address uint32: 0x%08X\n", net.BroadcastAddressUINT32)
	fmt.Printf("Broadcast Address: %v\n", net.BroadcastAddress)
	fmt.Printf("Network Address uint32: 0x%08X\n", net.NetworkAddressUINT32)
	fmt.Printf("Network Address: %v\n", net.NetworkAddress)
	fmt.Printf("Subnet Mask uint32: 0x%08X\n", net.SubnetMaskUINT32)
	fmt.Printf("Subnet Mask: %v\n", net.SubnetMask)
	fmt.Printf("Wildcard uint32: 0x%08X\n", net.WildcardUINT32)
	fmt.Printf("Wildcard: %v\n", net.Wildcard)
	fmt.Printf("Subnet Bitmap: %q\n", net.SubnetBitmap)
	fmt.Printf("Number of Hosts: %d\n", net.HostsMAX)
}
