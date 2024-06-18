package models

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type Network struct {
	address string
	mask    string
}

func NewNetwork(addrmask string) (Network, error) {
	parts := strings.Split(addrmask, "/")
	if len(parts) < 2 {
		return Network{}, errors.New("неправильно задана сеть")
	}
	return Network{address: parts[0], mask: parts[1]}, nil
}

func (n Network) String() string {
	return fmt.Sprintf("%s/%s", n.address, n.mask)
}

func (n Network) GetIPs() []string {
	var IPs []string
	// convert string to IPNet struct
	_, ipv4Net, err := net.ParseCIDR(n.String())
	if err != nil {
		log.Fatal(err)
	}

	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	// loop through addresses as uint32
	for i := start; i <= finish; i++ {
		// convert back to net.IP
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		IPs = append(IPs, ip.String())
	}
	return IPs
}
