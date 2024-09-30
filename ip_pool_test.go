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
	p, e := NewIPPool(mustNewIPRange("10.16.0.1", "10.16.0.2"))
	assert.Nil(t, e)
	var begin uint32 = 168820737
	t.Log(IPv4ToUint32(net.IPv4(10, 16, 0, 1)))

	// put一个超过ipRange的ip => ErrIPOutOfRange
	e = p.PutIPUint32(begin + 2)
	assert.NotNil(t, e)
	assert.Equal(t, e, ErrIPOutOfRange)

	// put一个未分配的ip => ErrIPNotAllocated
	e = p.PutIPUint32(begin)
	assert.NotNil(t, e)
	assert.Equal(t, e, ErrIPNotAllocated)

	// put一个已分配的ipu32 => 正常
	n, e := p.GetIPUint32()
	assert.Nil(t, e)
	assert.Equal(t, n, begin)
	e = p.PutIPUint32(n)
	assert.Nil(t, e)
	// 二次put同一个ipu32 => ErrIPAlreadyRecycled
	e = p.PutIPUint32(n)
	assert.NotNil(t, e)
	assert.Equal(t, e, ErrIPAlreadyRecycled)

	// put后get => 从recycled中取
	n, e = p.GetIPUint32()
	assert.Nil(t, e)
	assert.Equal(t, begin, n)
	// 再次get  => 从ipRange中取
	n, e = p.GetIPUint32()
	assert.Nil(t, e)
	assert.Equal(t, begin+1, n)
}

func TestIPPool_GetIP(t *testing.T) {
	p, e := NewIPPool(mustNewIPRange("0.0.0.0", "0.0.0.2"))
	assert.Nil(t, e)

	// 就测试接口能否正确转换
	ip, e := p.GetIP()
	assert.Nil(t, e)
	assert.Equal(t, ip, net.IPv4(0, 0, 0, 0).To4())
}

func TestIPPool_GetValidIP(t *testing.T) {
	p, e := NewIPPool(mustNewIPRange("0.0.0.254", "0.0.1.1"))
	assert.Nil(t, e)

	// 第1个IP合法，直接返回
	ip, e := p.GetValidIP()
	assert.Nil(t, e)
	assert.Equal(t, ip, net.IPv4(0, 0, 0, 254).To4())

	// 第2个IP尾0.0.0.255，不合法，返回下下个
	ip, e = p.GetValidIP()
	assert.Nil(t, e)
	assert.Equal(t, ip, net.IPv4(0, 0, 1, 1).To4())

}
