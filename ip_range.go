package iputils

import (
	"fmt"
	"math"
	"net"
)

type IPRange struct {
	beginNum uint32
	endNum   uint32
}

func NewIPRange(begin, end interface{}) (r *IPRange, e error) {
	ipv4Begin := AsIPv4(begin)
	ipv4End := AsIPv4(end)

	beginNum := IPv4ToUint32(ipv4Begin)
	endNum := IPv4ToUint32(ipv4End)
	if beginNum <= endNum {
		r = &IPRange{
			beginNum: beginNum,
			endNum:   endNum,
		}
	} else {
		e = fmt.Errorf("ipRange begin > end")
	}
	return
}

func NewIPRangeFromIPNet(ipnet *net.IPNet) (r *IPRange, e error) {
	beginNum := AsIPv4Uint32(ipnet.IP)
	maskOnes, _ := ipnet.Mask.Size()
	count := math.Pow(float64(2), float64(32-maskOnes))
	endNum := beginNum + uint32(count) - 1

	return NewIPRange(beginNum, endNum)
}

func NewIPRangeFromCIDR(cidr string) (r *IPRange, e error) {
	var ipnet *net.IPNet
	if _, ipnet, e = net.ParseCIDR(cidr); e == nil {
		return NewIPRangeFromIPNet(ipnet)
	}
	return
}

func MustNewIPRange(begin, end interface{}) (r *IPRange) {
	var e error
	r, e = NewIPRange(begin, end)
	if e != nil {
		panic(e)
	}
	return
}

func (r IPRange) FirstIP() net.IP {
	return Uint32ToIPv4(r.beginNum)
}

func (r IPRange) LastIP() net.IP {
	return Uint32ToIPv4(r.endNum)
}

func (r IPRange) String() string {
	return fmt.Sprintf("[%s - %s]", r.FirstIP(), r.LastIP())
}

func (r *IPRange) Has(ip interface{}) bool {
	n := AsIPv4Uint32(ip)
	return r.beginNum <= n && n <= r.endNum
}

func (r *IPRange) HasOverlap(r2 *IPRange) bool {
	return !(r2.beginNum > r.endNum || r2.endNum < r.beginNum)
}

func (r *IPRange) Split(n uint32) (r1 *IPRange, r2 *IPRange, e error) {
	if r1, e = NewIPRange(r.beginNum, r.beginNum+n-1); e == nil {
		r2, e = NewIPRange(r1.endNum+1, r.endNum)
	}
	return
}

func (r *IPRange) PopFirst() (IP net.IP, e error) {
	if r.beginNum >= r.endNum {
		return nil, fmt.Errorf("ipRange begin >= end")
	}

	IP = r.FirstIP()
	r.beginNum += 1
	return
}

func (r *IPRange) PopLast() (IP net.IP, e error) {
	if r.beginNum >= r.endNum {
		return nil, fmt.Errorf("ipRange begin >= end")
	}

	IP = r.LastIP()
	r.endNum -= 1
	return
}

func (r *IPRange) TrimFirstN(n uint32) (ipRange *IPRange, e error) {
	if r.beginNum+n > r.endNum {
		return nil, fmt.Errorf("ipRange begin + n > end")
	}
	if ipRange, e = NewIPRange(r.beginNum, r.beginNum+n-1); e == nil {
		r.beginNum += n
	}
	return
}

func (r *IPRange) TrimLastN(n uint32) (ipRange *IPRange, e error) {
	if r.beginNum > r.endNum-n {
		return nil, fmt.Errorf("ipRange begin > end - n")
	}
	if ipRange, e = NewIPRange(r.endNum-n+1, r.endNum); e == nil {
		r.endNum -= n
	}
	return
}

func (r *IPRange) Trim(firstN uint32, lastN uint32) (ipRange *IPRange, e error) {
	beginNum := r.beginNum + firstN
	endNum := r.endNum - lastN
	ipRange, e = NewIPRange(beginNum, endNum)
	return
}
