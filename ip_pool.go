package iputils

import (
	"fmt"
	"net"
	"sync"
)

// 参考 bitmap.go

// todo: check 返回的时候是否跳过.0和.255
const maxIPCount = 2 << 24 // 10.x.x.x 有2^24个IP，占用大概4MB内存
var ErrIPExhausted = fmt.Errorf("ip exhausted")
var ErrIPOutOfRange = fmt.Errorf("ip out of range")
var ErrIPNotAllocated = fmt.Errorf("ip not allocated")
var ErrIPAlreadyRecycled = fmt.Errorf("ip already recycled")

type IPPool struct {
	sync.Mutex

	ipRange       *IPRange
	nextGetOffset uint32 // 下一次从ipRange中分配IP的偏移

	recycleIPs *bitmap // PutIP会将IP放入recycleIPs中，GetIP时优先从reuseIPs中取
}

func NewIPPool(ipRange *IPRange) (p *IPPool, e error) {
	if ipRange.Size() > maxIPCount {
		return nil, fmt.Errorf("ipRange.Size too large, > %d", maxIPCount)
	}

	return &IPPool{
		ipRange:    ipRange,
		recycleIPs: newBitmap(ipRange.Size()),
	}, nil
}

func (a *IPPool) GetIPUint32() (ipu32 uint32, e error) {
	a.Lock()
	defer a.Unlock()
	if i := a.recycleIPs.FirstFlagOffset(); i < 0 {
		// 没有可重用IP，从ipRange中取
		if a.nextGetOffset < a.ipRange.Size() {
			ipu32 = a.ipRange.beginNum + a.nextGetOffset
			a.nextGetOffset++
			return
		} else {
			return 0, ErrIPExhausted
		}
	} else {
		a.recycleIPs.UnsetBit(uint32(i))
		return uint32(i), nil
	}
}

func (a *IPPool) PutIPUint32(ipu32 uint32) error {
	a.Lock()
	defer a.Unlock()

	if a.ipRange.Has(ipu32) {
		// 如果ipu32没有被分配，那么报错
		if ipu32 >= a.ipRange.beginNum+a.nextGetOffset {
			return ErrIPNotAllocated
		} else {
			var i uint32 = ipu32 - a.ipRange.beginNum
			if a.recycleIPs.GetBit(i) {
				return ErrIPAlreadyRecycled
			} else {
				a.recycleIPs.SetBit(ipu32 - a.ipRange.beginNum)
				return nil
			}
		}
	} else {
		return ErrIPOutOfRange
	}
}

func (a *IPPool) GetIP() (ipv4 net.IP, e error) {
	var num uint32
	if num, e = a.GetIPUint32(); e == nil {
		return Uint32ToIPv4(num), nil
	}
	return
}

func (a *IPPool) PutIP(ip net.IP) (e error) {
	return a.PutIPUint32(AsIPv4Uint32(ip))
}
