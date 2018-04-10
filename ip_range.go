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

func (r *IPRange) HasIntersection(r2 *IPRange) bool {
	return !(r2.beginNum > r.endNum || r2.endNum < r.beginNum)
}

func (r *IPRange) Has(ip interface{}) bool {
	n := AsIPv4Uint32(ip)
	return r.beginNum <= n && n <= r.endNum
}

func (r *IPRange) Split(n uint32) (r1 *IPRange, r2 *IPRange, e error) {
	if r1, e = NewIPRange(r.beginNum, r.beginNum+n-1); e == nil {
		r2, e = NewIPRange(r1.endNum+1, r.endNum)
	}
	return
}

func (r *IPRange) PopLeft() (IP net.IP, e error) {
	if r.beginNum >= r.endNum {
		return nil, fmt.Errorf("ipRange begin > end")
	}

	IP = r.FirstIP()
	r.beginNum += 1
	return
}

func (r *IPRange) PopRight() (IP net.IP, e error) {
	if r.beginNum >= r.endNum {
		return nil, fmt.Errorf("ipRange begin > end")
	}

	IP = r.LastIP()
	r.endNum -= 1
	return
}

func (r *IPRange) TrimLeft(count uint32) (ipRange *IPRange, e error) {
	if r.beginNum+count > r.endNum {
		return nil, fmt.Errorf("ipRange begin + count > end")
	}
	if ipRange, e = NewIPRange(r.beginNum, r.beginNum+count-1); e == nil {
		r.beginNum += count
	}
	return
}

func (r *IPRange) TrimRight(count uint32) (ipRange *IPRange, e error) {
	if r.beginNum > r.endNum-count {
		return nil, fmt.Errorf("ipRange begin > end - count")
	}
	if ipRange, e = NewIPRange(r.endNum-count+1, r.endNum); e == nil {
		r.endNum -= count
	}
	return
}

func (r *IPRange) Trim(leftCount uint32, rightCount uint32) (ipRange *IPRange, e error) {
	beginNum := r.beginNum + leftCount
	endNum := r.endNum - rightCount
	ipRange, e = NewIPRange(beginNum, endNum)
	return
}
