### CIDR Calculator

###### USAGE:

<code>go get github.com/rootVIII/cidrcalc</code>

Build the example:
<code>go build example/main.go</code>
<code>./main -i 192.168.1.1/24  // example CIDR</code>


The example main.go demonstrates how to use the module in your own code.
The Calculate() method expects the IP/CIDR as a string: 192.168.1.1/24
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
    fmt.Printf("IP uint32: 0x%08X\n", net.IPUINT32)
    fmt.Printf("CIDR: %d\n", net.CIDR)
    fmt.Printf("Broadcast Address uint32: 0x%08X\n", net.BroadcastAddressUINT32)
    fmt.Printf("Broadcast Address bytes: %v\n", net.BroadcastAddress)
    fmt.Printf("Network Address uint32: 0x%08X\n", net.NetworkAddressUINT32)
    fmt.Printf("Network Address: %v\n", net.NetworkAddress)
    fmt.Printf("Subnet Mask: %v\n", net.SubnetMask)
    fmt.Printf("SubnetMask uint32: 0x%08X\n", net.SubnetMaskUINT32)
    fmt.Printf("Subnet Bitmap: %q\n", net.SubnetBitmap)
    fmt.Printf("Number of Hosts: %d\n", net.HostsMAX)
}
  </code>
</pre>


###### Example Output:
<pre>
  <code>
IP uint32: 0xC0A80101
CIDR: 24
Broadcast Address uint32: 0xC0A801FF
Broadcast Address bytes: [192 168 1 255]
Network Address uint32: 0xC0A80100
Network Address: [192 168 1 0]
Subnet Mask: [255 255 255 0]
SubnetMask uint32: 0xFFFFFF00
Subnet Bitmap: "nnnnnnnnnnnnnnnnnnnnnnnnhhhhhhhh"
Number of Hosts: 254
  </code>
</pre>



This was developed on Ubuntu 18.04 LTS.
<hr>
<b>Author: rootVIII  2020</b>
<br><br>
