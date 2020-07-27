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
	CIDR                 byte
	BroadcastAddress     [4]byte
	IP                   [4]byte
	IPUINT32             uint32
	SubnetMask           [4]byte
	SubnetMaskUINT32     uint32
	NetworkAddress       [4]byte
	NetworkAddressUINT32 uint32
	SubnetBitmap         []byte
	HostsMAX             uint
	SubnetsMAX           uint
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

	fmt.Printf("Number of Hosts: %d\n", s.HostsMAX)

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
	fmt.Printf("Network AddressUINT32: 0%X\n", s.NetworkAddressUINT32)
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

	// TODO: fix this. get number of subnets and also
	// broadcast address: add two uint32s and convert back?
	sep, byteCount := 0, 0
	for index, sbyte := range s.SubnetMask {
		byteCount = index
		if bits.TrailingZeros8(sbyte) != 0 {
			fmt.Printf("%b\n", sbyte)
			sep = bits.OnesCount8(sbyte)
			break
		}
	}
	fmt.Printf("\nonbits %d bytecount: %d\n\n", sep, byteCount)
}

// Calculate is the public method used to set all type Subnet attributes.
func (s *Subnet) Calculate(IPCIDR string) error {
	_, _, err := net.ParseCIDR(IPCIDR)
	if err != nil {
		return err
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
