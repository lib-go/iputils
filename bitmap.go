package iputils

type bitmap struct {
	bits      []byte // 存放具体的bitmap
	totalBits uint32 // bitmap总共有多少位
	flags     uint32 // bitmap中有多少个1
}

func newBitmap(totalBits uint32) *bitmap {
	return &bitmap{
		bits:      make([]byte, totalBits/8+1), // 根据totalBits来分配bitmap的空间
		totalBits: totalBits,
	}
}

func (b *bitmap) SetBit(i uint32) {
	if i < b.totalBits {
		offset := i / 8
		// 当该位为0时，才设置为1
		if b.bits[offset]&(1<<(i%8)) == 0 {
			b.bits[offset] |= 1 << (i % 8)
			b.flags++
		}
	}
}

func (b *bitmap) UnsetBit(i uint32) {
	if i < b.totalBits {
		offset := i / 8
		// 当该位为1时，才设置为0
		if b.bits[offset]&(1<<(i%8)) > 0 {
			b.bits[offset] &^= 1 << (i % 8)
			b.flags--
		}
	}
}

func (b *bitmap) GetBit(i uint32) bool {
	if i < b.totalBits {
		return b.bits[i/8]&(1<<(i%8)) > 0
	}
	return false
}

// FirstFlagOffset 返回第一个为1的位的偏移量，如果没有为1的位，则返回-1
func (b *bitmap) FirstFlagOffset() int {
	if b.flags == 0 {
		return -1
	}

	// 先以字节为单位，找到第一个不为0的字节
	for i, v := range b.bits {
		if v > 0 {
			// 找到第一个不为0的字节后，再在该字节中找到第一个为1的位
			for j := 0; j < 8; j++ {
				if v&(1<<uint(j)) > 0 {
					return i*8 + j
				}
			}
		}
	}
	return -1
}

func (b *bitmap) Ones() uint32 {
	return b.flags
}
