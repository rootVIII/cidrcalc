### CIDR Calculator

###### USAGE AS A PACKAGE:
The example main.go demonstrates how to use the module in your own code.<br>
The Calculate() method expects any IP/CIDR as a string, ex: 192.168.1.1/24
<pre>
  <code>
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
        fmt.Printf("&#37;v\n", err)
        os.Exit(1)
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
  </code>
</pre>


###### EXAMPLE OUTPUT:
<pre>
  <code>
CIDR: 24
IP Address uint32: 0xC0A80101
IP Address: [192 168 1 1]
Broadcast Address uint32: 0xC0A801FF
Broadcast Address: [192 168 1 255]
Network Address uint32: 0xC0A80100
Network Address: [192 168 1 0]
Subnet Mask uint32: 0xFFFFFF00
Subnet Mask: [255 255 255 0]
Wildcard uint32: 0x000000FF
Wildcard: [0 0 0 255]
Subnet Bitmap: "nnnnnnnnnnnnnnnnnnnnnnnnhhhhhhhh"
Number of Hosts: 254
  </code>
</pre>


###### TEST THE EXAMPLE:

<code>go get github.com/rootVIII/cidrcalc</code><br>
<code>go build example/main.go</code><br>
<code>./main -i 192.168.1.1/24  // example CIDR</code>


This was developed on Ubuntu 18.04 LTS.
<hr>
<b>Author: rootVIII  2020</b>
<br><br>
