package iputils

import (
	"fmt"
	"math/bits"
	"net"
	"sync"
)

const ChunkSize = 2 ^ 32/8

var ErrEmptyBuff error = fmt.Errorf("empty bins")
var ErrIPExhausted error = fmt.Errorf("ip exhausted")

type IPPool struct {
	sync.Mutex

	ipRange *IPRange
	offset  uint32 // ipRange中已经分配掉的数量

	bins            [][]byte // 记录放回的IP，以bit为单位，每个[]byte长度为ChunkSize
	binIPCounts     []uint32 // 每个bin的可用IP数量
	totalBinIPCount int      // 所有bin的IP数量

	loop bool // 如果为true，用完IP后会从头使用
}

func NewIPPool(ipRange IPRange) (*IPPool, error) {
	if ipRange.beginNum > ipRange.endNum {
		return nil, fmt.Errorf("ipRange begin > end")
	}

	chunkN := int((ipRange.endNum-ipRange.beginNum)/ChunkSize) + 1
	return &IPPool{
		ipRange:         &ipRange,
		bins:            make([][]byte, chunkN),
		totalBinIPCount: 0,
		binIPCounts:     make([]uint32, chunkN),
	}, nil
}

func (a *IPPool) getNumFromBins() (num uint32, e error) {
	e = ErrEmptyBuff

	if a.totalBinIPCount == 0 {
		return
	}

	var offset, segment int
	var count uint32
	var byte_ byte
	for segment, count = range a.binIPCounts {
		if count > 0 {
			chunk := a.bins[segment]
			for offset, byte_ = range chunk {
				if i := bits.Len8(byte_) - 1; i > -1 {
					a.binIPCounts[segment] -= 1
					a.totalBinIPCount -= 1
					if a.binIPCounts[segment] == 0 {
						a.bins[segment] = nil
					}
					chunk[offset] &^= 1 << uint32(i)
					return a.ipRange.beginNum + uint32(segment*ChunkSize+offset*8+i), nil
				}
			}
		}
	}
	return
}

func (a *IPPool) GetNum() (num uint32, e error) {
	a.Lock()

	if num, e = a.getNumFromBins(); e != nil {
		for {
			num = a.ipRange.beginNum + a.offset
			if num > a.ipRange.endNum {
				if a.loop == false {
					e = ErrIPExhausted
					break
				}
				a.offset = 0
				a.Unlock()
				return a.GetNum()
			}

			a.offset += 1

			d4 := num % 256
			if d4 == 0 || d4 == 255 {
				continue
			} else {
				e = nil
				break
			}
		}
	}
	a.Unlock()
	return
}

func (a *IPPool) PutNum(num uint32) error {
	a.Lock()
	defer a.Unlock()

	if !a.ipRange.Has(num) {
		return fmt.Errorf("ip not in range")
	}

	if a.ipRange.beginNum+a.offset <= num && num <= a.ipRange.endNum {
		return fmt.Errorf("ip still unused")
	}

	d4 := num % 256
	if d4 == 0 || d4 == 255 {
		return fmt.Errorf("ip not valid")
	}

	var binIndex, bitOffset uint32
	binIndex = (num - a.ipRange.beginNum) / ChunkSize
	bitOffset = (num - a.ipRange.beginNum) % ChunkSize

	bin := a.bins[binIndex]
	if bin == nil {
		bin = make([]byte, ChunkSize)
		a.bins[binIndex] = bin
	}
	bit := byte(1 << (bitOffset % 8))
	if (bin[bitOffset>>3] & bit) == bit {
		return fmt.Errorf("num already in pool")
	}

	bin[bitOffset>>3] |= bit
	a.totalBinIPCount += 1
	a.binIPCounts[binIndex] += 1
	return nil
}

func (a *IPPool) GetIP() (ipv4 net.IP, e error) {
	var num uint32
	if num, e = a.GetNum(); e == nil {
		return Uint32ToIPv4(num), nil
	}
	return
}

func (a *IPPool) PutIP(ip net.IP) (e error) {
	return a.PutNum(AsIPv4Uint32(ip))
}

func (a *IPPool) SetLoop(loop bool) {
	a.loop = loop
}
