package iputils

import (
	"encoding/binary"
	"net"
)

var privateNets = []*net.IPNet{
	{IP: net.IPv4(127, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
	{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
	{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
	{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
}

func IsPrivateIPv4(ip interface{}) (b bool) {
	if _ip := AsIPv4(ip); _ip != nil {
		for _, ipnet := range privateNets {
			if ipnet.Contains(_ip) {
				return true
			}
		}
	}
	return false
}

func AsIPv4(ip interface{}) (_ip net.IP) {
	switch x := ip.(type) {
	case net.IP:
		_ip = x.To4()
	case []byte:
		if len(x) == 4 {
			_ip = net.IP(x)
		}
	case string:
		_ip = net.ParseIP(x).To4()
	case int:
		_ip = Uint32ToIPv4(uint32(x))
	case uint32:
		_ip = Uint32ToIPv4(x)
	}
	return
}

func AsIPv4Uint32(ip interface{}) uint32 {
	switch x := ip.(type) {
	case uint32:
		return x
	case net.IP:
		return IPv4ToUint32(x)
	case []byte:
		return IPv4ToUint32(net.IP(x))
	case string:
		return IPv4ToUint32(net.ParseIP(x))
	case int:
		return uint32(x)
	}
	return 0
}

func IsIPv4(ip interface{}) bool {
	_ip := AsIPv4(ip)
	return _ip != nil && !_ip.IsUnspecified()
}

func IPv4ToUint32(ip net.IP) uint32 {
	if len(ip) == 0 {
		return 0
	}
	return binary.BigEndian.Uint32(ip.To4())
}

func Uint32ToIPv4(n uint32) net.IP {
	bytes := make([]byte, 4)

	binary.BigEndian.PutUint32(bytes, n)
	return net.IP(bytes)
}

func IPv4Inc(ip net.IP, inc int) net.IP {
	n := int(IPv4ToUint32(ip))
	return Uint32ToIPv4(uint32(int(n) + inc))
}

func IPv4End(ipnet *net.IPNet) net.IP {
	ip := ipnet.IP.To4()
	mask := ipnet.Mask

	ret := net.IPv4(0, 0, 0, 0).To4()
	for i, d := range mask {
		ret[i] = ip[i] | (255 &^ d)
	}

	return ret
}
