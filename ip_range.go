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
	endNum := beginNum + uint32(count-1)

	return NewIPRange(beginNum, endNum)
}

func NewIPRangeFromCIDR(cidr string) (r *IPRange, e error) {
	var ipnet *net.IPNet
	if _, ipnet, e = net.ParseCIDR(cidr); e == nil {
		return NewIPRangeFromIPNet(ipnet)
	}
	return
}

func (r IPRange) String() string {
	return fmt.Sprintf("[%s - %s]", r.At(0), r.At(-1))
}

func (r IPRange) Size() uint32 {
	return r.endNum - r.beginNum + 1
}

func (r IPRange) BeginNum() uint32 {
	return r.beginNum
}

func (r IPRange) EndNum() uint32 {
	return r.endNum
}

func (r IPRange) BeginIP() net.IP {
	return Uint32ToIPv4(r.beginNum)
}

func (r IPRange) EndIP() net.IP {
	return Uint32ToIPv4(r.endNum)
}

func (r IPRange) At(i int) net.IP {
	if i >= 0 {
		return Uint32ToIPv4(r.beginNum + uint32(i))
	} else {
		return Uint32ToIPv4(r.endNum + uint32(-i-1))
	}
}

func (r IPRange) Has(ip interface{}) bool {
	n := AsIPv4Uint32(ip)
	return r.beginNum <= n && n <= r.endNum
}

func (r IPRange) HasOverlap(r2 *IPRange) bool {
	return !(r2.beginNum > r.endNum || r2.endNum < r.beginNum)
}

func (r IPRange) Split(n uint32) (r1 *IPRange, r2 *IPRange, e error) {
	if r1, e = NewIPRange(r.beginNum, r.beginNum+n-1); e == nil {
		r2, e = NewIPRange(r1.endNum+1, r.endNum)
	}
	return
}

func (r IPRange) TrimCopy(firstN uint32, lastN uint32) (ipRange *IPRange, e error) {
	beginNum := r.beginNum + firstN
	endNum := r.endNum - lastN
	ipRange, e = NewIPRange(beginNum, endNum)
	return
}
