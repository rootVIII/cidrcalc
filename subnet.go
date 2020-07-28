package cidrcalc

import (
	"bytes"
	"fmt"
	"math/bits"
	"net"
	"strconv"
	"strings"
)

// Subnet represents network attributes for the given IP address.
type Subnet struct {
	CIDR                   byte
	BroadcastAddress       [4]byte
	BroadcastAddressUINT32 uint32
	IP                     [4]byte
	IPUINT32               uint32
	SubnetMask             [4]byte
	SubnetMaskUINT32       uint32
	NetworkAddress         [4]byte
	NetworkAddressUINT32   uint32
	SubnetBitmap           []byte
	HostsMAX               uint32
}

func (s *Subnet) setSubnetMask(out chan<- struct{}) {
	s.SubnetMask[0] = uint8((s.SubnetMaskUINT32 & 0xFF000000) >> 24)
	s.SubnetMask[1] = uint8((s.SubnetMaskUINT32 & 0x00FF0000) >> 16)
	s.SubnetMask[2] = uint8((s.SubnetMaskUINT32 & 0x0000FF00) >> 8)
	s.SubnetMask[3] = uint8(s.SubnetMaskUINT32 & 0x000000FF)

	fmt.Printf("Subnet Mask: %v\n", s.SubnetMask)
	fmt.Printf("SubnetMaskUINT32: 0x%X\n", s.SubnetMaskUINT32)

	out <- struct{}{}
}

func (s *Subnet) setSubnetBitmap(onBits int, offBits int, out chan<- struct{}) {
	var netmap bytes.Buffer
	var count int
	for count = 0; count < onBits; count++ {
		netmap.WriteByte(0x6E)
	}
	for count = 0; count < offBits; count++ {
		netmap.WriteByte(0x68)
	}
	s.SubnetBitmap = netmap.Bytes()

	fmt.Printf("Subnet Bitmap: %q\n", s.SubnetBitmap)

	out <- struct{}{}
}

func (s *Subnet) setMaxHosts(offBits int, out chan<- struct{}) {
	s.HostsMAX = 1 << offBits
	out <- struct{}{}
}

func (s *Subnet) setNetworkID(out chan<- struct{}) {
	s.IPUINT32 = uint32(s.IP[0]) << 24
	s.IPUINT32 += uint32(s.IP[1]) << 16
	s.IPUINT32 += uint32(s.IP[2]) << 8
	s.IPUINT32 += uint32(s.IP[3])
	s.NetworkAddressUINT32 = s.IPUINT32 & s.SubnetMaskUINT32

	s.NetworkAddress[0] = uint8((s.NetworkAddressUINT32 & 0xFF000000) >> 24)
	s.NetworkAddress[1] = uint8((s.NetworkAddressUINT32 & 0x00FF0000) >> 16)
	s.NetworkAddress[2] = uint8((s.NetworkAddressUINT32 & 0x0000FF00) >> 8)
	s.NetworkAddress[3] = uint8(s.NetworkAddressUINT32 & 0x000000FF)

	fmt.Printf("IP32: 0%X\n", s.IPUINT32)
	fmt.Printf("Network AddressUINT32: 0x%X\n", s.NetworkAddressUINT32)
	fmt.Printf("Network Address: %v\n", s.NetworkAddress)

	out <- struct{}{}
}

func (s *Subnet) mask() {
	trailing := 32 - s.CIDR
	s.SubnetMaskUINT32 = (0xFFFFFFFF >> trailing) << trailing
	netCH := make(chan struct{})
	nOn := bits.OnesCount32(s.SubnetMaskUINT32)
	nOff := bits.TrailingZeros32(s.SubnetMaskUINT32)
	go s.setSubnetBitmap(nOn, nOff, netCH)
	go s.setSubnetMask(netCH)
	go s.setMaxHosts(nOff, netCH)
	go s.setNetworkID(netCH)
	for range [4]uint8{} {
		<-netCH
	}

	s.BroadcastAddressUINT32 = (s.NetworkAddressUINT32 + s.HostsMAX) - 1
	s.BroadcastAddress[0] = uint8((s.BroadcastAddressUINT32 & 0xFF000000) >> 24)
	s.BroadcastAddress[1] = uint8((s.BroadcastAddressUINT32 & 0x00FF0000) >> 16)
	s.BroadcastAddress[2] = uint8((s.BroadcastAddressUINT32 & 0x0000FF00) >> 8)
	s.BroadcastAddress[3] = uint8(s.BroadcastAddressUINT32 & 0x000000FF)

	fmt.Printf("Broadcast AddressUINT32: 0x%08X\n", s.BroadcastAddressUINT32)
	fmt.Printf("Broadcast Address: %v\n", s.BroadcastAddress)

	fmt.Printf("CIDR: %d\n", s.CIDR)

	if s.CIDR != 32 {
		s.HostsMAX -= 2
	}
	fmt.Printf("Number of Hosts: %d\n", s.HostsMAX)
}

// Calculate is the public method used to set all type Subnet attributes.
func (s *Subnet) Calculate(IPCIDR string) error {
	_, _, err := net.ParseCIDR(IPCIDR)
	if err != nil {
		return err
	}

	if strings.Count(IPCIDR, ".") != 3 || strings.Contains(IPCIDR, ":") {
		return fmt.Errorf("invalid IPV4 address: %s", IPCIDR)
	}

	ipv4Arr := strings.Split(IPCIDR, "/")
	CIDR, _ := strconv.Atoi(ipv4Arr[1])
	for index, octet := range strings.Split(ipv4Arr[0], ".") {
		val, _ := strconv.Atoi(octet)
		s.IP[index] = uint8(val)
	}
	s.CIDR = uint8(CIDR)
	s.mask()
	return err
}
