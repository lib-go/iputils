package iputils

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestIPPool_GetIPUint32(t *testing.T) {
	p, e := NewIPPool(mustNewIPRange("0.0.0.0", "0.0.0.2"))
	assert.Nil(t, e)

	ipu32, e := p.GetIPUint32()
	assert.Nil(t, e)
	assert.Equal(t, ipu32, uint32(0))

	ipu32, e = p.GetIPUint32()
	assert.Nil(t, e)
	assert.Equal(t, ipu32, uint32(1))

	ipu32, e = p.GetIPUint32()
	assert.Nil(t, e)
	assert.Equal(t, ipu32, uint32(2))

	// get次数超过总IP数 => ErrIPExhausted
	ipu32, e = p.GetIPUint32()
	assert.NotNil(t, e)
	assert.Equal(t, e, ErrIPExhausted)
}

func TestIPPool_PutIPUint32(t *testing.T) {
	p, e := NewIPPool(mustNewIPRange("0.0.0.0", "0.0.0.2"))
	assert.Nil(t, e)

	// put一个超过ipRange的ip => ErrIPOutOfRange
	e = p.PutIPUint32(3)
	assert.NotNil(t, e)
	assert.Equal(t, e, ErrIPOutOfRange)

	// put一个未分配的ip => ErrIPNotAllocated
	e = p.PutIPUint32(0)
	assert.NotNil(t, e)
	assert.Equal(t, e, ErrIPNotAllocated)

	// put一个已分配的ipu32 => 正常
	n, e := p.GetIPUint32()
	assert.Nil(t, e)
	assert.Equal(t, n, uint32(0))
	e = p.PutIPUint32(n)
	assert.Nil(t, e)
	// 二次put同一个ipu32 => ErrIPAlreadyRecycled
	e = p.PutIPUint32(n)
	assert.NotNil(t, e)
	assert.Equal(t, e, ErrIPAlreadyRecycled)

	// put后get => 从recycled中取
	n, e = p.GetIPUint32()
	assert.Nil(t, e)
	assert.Equal(t, n, uint32(0))
	// 再次get  => 从ipRange中取
	n, e = p.GetIPUint32()
	assert.Nil(t, e)
	assert.Equal(t, n, uint32(1))
}

func TestIPPool_GetIP(t *testing.T) {
	p, e := NewIPPool(mustNewIPRange("0.0.0.0", "0.0.0.2"))
	assert.Nil(t, e)

	// 就测试接口能否正确转换
	ip, e := p.GetIP()
	assert.Nil(t, e)
	assert.Equal(t, ip, net.IPv4(0, 0, 0, 0).To4())
}
