package iputils

import (
	"fmt"
	"math/bits"
	"net"
	"sync"
)

const CHUNK_SIZE = 8192

// 速度优化
var ErrEmptyBuff error = fmt.Errorf("empty buf")
var ErrIPExhausted error = fmt.Errorf("ip exhausted")

type IPPool struct {
	sync.Mutex
	ipRange       IPRange
	buf           [][]byte
	bufIPCount    int
	chunkIPCounts []uint32 // 方便chunk取完后设nil
	chunkN        int
	loop          bool
	offset        uint32
}

func NewIPPool(ipRange IPRange) (*IPPool, error) {
	if ipRange.beginNum > ipRange.endNum {
		return nil, fmt.Errorf("ipRange begin > end")
	}

	chunkN := int((ipRange.endNum-ipRange.beginNum)/CHUNK_SIZE) + 1
	return &IPPool{
		ipRange:       ipRange,
		buf:           make([][]byte, chunkN),
		bufIPCount:    0,
		chunkIPCounts: make([]uint32, chunkN),
		chunkN:        chunkN,
	}, nil
}

func (a *IPPool) getNumFromBuf() (num uint32, e error) {
	e = ErrEmptyBuff

	if a.bufIPCount == 0 {
		return
	}

	var offset, segment int
	var count uint32
	var byte_ byte
	for segment, count = range a.chunkIPCounts {
		if count > 0 {
			chunk := a.buf[segment]
			for offset, byte_ = range chunk {
				if i := bits.Len8(byte_) - 1; i > -1 {
					a.chunkIPCounts[segment] -= 1
					a.bufIPCount -= 1
					if a.chunkIPCounts[segment] == 0 {
						a.buf[segment] = nil
					}
					chunk[offset] &^= 1 << uint32(i)
					return a.ipRange.beginNum + uint32(segment*CHUNK_SIZE+offset*8+i), nil
				}
			}
		}
	}
	return
}

func (a *IPPool) GetNum() (num uint32, e error) {
	a.Lock()

	if num, e = a.getNumFromBuf(); e != nil {
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

	if !a.ipRange.Has(num) {
		a.Unlock()
		return fmt.Errorf("ip not in range")
	}

	if a.ipRange.beginNum+a.offset <= num && num <= a.ipRange.endNum {
		a.Unlock()
		return fmt.Errorf("ip still exist")
	}

	d4 := num % 256
	if d4 == 0 || d4 == 255 {
		a.Unlock()
		return fmt.Errorf("ip not valid")
	}

	var segment, offset uint32
	var byte_ byte
	segment = (num - a.ipRange.beginNum) / CHUNK_SIZE
	offset = (num - a.ipRange.beginNum) % CHUNK_SIZE

	chunk := a.buf[segment]
	if chunk == nil {
		chunk = make([]byte, CHUNK_SIZE)
		a.buf[segment] = chunk
	}
	byte_ = chunk[offset>>3]
	numByte := byte(offset % 8)
	bit := byte(1 << numByte)
	if (byte_ & bit) == bit {
		a.Unlock()
		return fmt.Errorf("num already in pool")
	}

	chunk[offset>>3] = byte_ | bit
	a.bufIPCount += 1
	a.chunkIPCounts[segment] += 1
	a.Unlock()
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
