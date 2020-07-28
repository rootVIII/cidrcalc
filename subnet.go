package cidrcalc

/*

	rootVIII 2020

*/

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
	IP                     [4]byte
	BroadcastAddress       [4]byte
	SubnetMask             [4]byte
	NetworkAddress         [4]byte
	IPUINT32               uint32
	BroadcastAddressUINT32 uint32
	SubnetMaskUINT32       uint32
	NetworkAddressUINT32   uint32
	SubnetBitmap           []byte
	HostsMAX               uint32
}

func (s *Subnet) toBytes(src uint32) [4]byte {
	tmp := [4]byte{}
	tmp[0] = uint8((src & 0xFF000000) >> 24)
	tmp[1] = uint8((src & 0x00FF0000) >> 16)
	tmp[2] = uint8((src & 0x0000FF00) >> 8)
	tmp[3] = uint8(src & 0x000000FF)
	return tmp
}

func (s *Subnet) setSubnetMask(out chan<- struct{}) {
	s.SubnetMask = s.toBytes(s.SubnetMaskUINT32)
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
	s.NetworkAddress = s.toBytes(s.NetworkAddressUINT32)
	out <- struct{}{}
}

func (s *Subnet) mask() {
	trailing := 32 - s.CIDR
	s.SubnetMaskUINT32 = (0xFFFFFFFF >> trailing) << trailing

	nOn := bits.OnesCount32(s.SubnetMaskUINT32)
	nOff := bits.TrailingZeros32(s.SubnetMaskUINT32)

	netCH := make(chan struct{})
	go s.setSubnetBitmap(nOn, nOff, netCH)
	go s.setSubnetMask(netCH)
	go s.setMaxHosts(nOff, netCH)
	go s.setNetworkID(netCH)
	for range [4]uint8{} {
		<-netCH
	}

	s.BroadcastAddressUINT32 = (s.NetworkAddressUINT32 + s.HostsMAX) - 1
	s.BroadcastAddress = s.toBytes(s.BroadcastAddressUINT32)
	if s.CIDR != 32 {
		s.HostsMAX -= 2
	}
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
