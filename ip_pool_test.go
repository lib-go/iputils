package iputils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestIPool(t *testing.T) {
	ipRange := MustNewIPRange("0.0.0.0", "0.0.0.2")
	a, e := NewIPPool(*ipRange)

	// pop不应该影响下面的测试
	_, _ = ipRange.PopFirst()
	_, _ = ipRange.PopFirst()

	ip, e := a.GetIP()
	if ip != nil {
		assert.True(t, net.IPv4(0, 0, 0, 1).Equal(ip))
	}
	assert.Nil(t, e)

	// 无法放入还存在的ip
	e = a.PutIP(net.IPv4(0, 0, 0, 2))
	assert.NotNil(t, e)

	// 正常放入
	e = a.PutIP(net.IPv4(0, 0, 0, 1))
	assert.Nil(t, e)

	// 正常取出
	ip, e = a.GetIP()
	assert.True(t, net.IPv4(0, 0, 0, 1).Equal(ip))
	assert.Nil(t, e)

	// 取完最后一个非池子的ip
	ip, e = a.GetIP()
	assert.True(t, net.IPv4(0, 0, 0, 2).Equal(ip))
	assert.Nil(t, e)

	// ip耗尽
	_, e = a.GetIP()
	assert.NotNil(t, e)

	// loop后，不再耗尽，从头开始
	a.SetLoop(true)
	ip, e = a.GetIP()
	assert.True(t, net.IPv4(0, 0, 0, 1).Equal(ip))
	assert.Nil(t, e)

	// 两次重复put
	e = a.PutIP(net.IPv4(0, 0, 0, 1))
	assert.Nil(t, e)

	e = a.PutIP(net.IPv4(0, 0, 0, 1))
	assert.NotNil(t, e)

	// range外
	e = a.PutIP(net.IPv4(0, 0, 0, 3))
	assert.NotNil(t, e)

	// 跨chunk取
	a, e = NewIPPool(*MustNewIPRange("0.0.0.0", "0.1.0.1"))
	a.offset = 0x10000

	e = a.PutIP(net.IPv4(0, 0, 0, 1))
	assert.Nil(t, e)

	e = a.PutIP(net.IPv4(0, 0, 0, 10))
	assert.Nil(t, e)

	ip, e = a.GetIP()
	assert.True(t, net.IPv4(0, 0, 0, 1).Equal(ip))

	ip, e = a.GetIP()
	assert.True(t, net.IPv4(0, 0, 0, 10).Equal(ip))

	// 避开0和255结尾的IP
	a, e = NewIPPool(*MustNewIPRange("0.0.0.254", "0.0.1.1"))

	ip, e = a.GetIP()
	assert.True(t, net.IPv4(0, 0, 0, 254).Equal(ip))
	assert.Nil(t, e)

	ip, e = a.GetIP()
	assert.True(t, net.IPv4(0, 0, 1, 1).Equal(ip))
	assert.Nil(t, e)
}

func BenchmarkIPPool_GetNum(b *testing.B) {
	a, _ := NewIPPool(*MustNewIPRange(0, 0x0000ffff))
	for i := 0; i < b.N; i++ {
		a.GetNum()
	}
}

func BenchmarkIPPool_PutNum(b *testing.B) {
	a, _ := NewIPPool(*MustNewIPRange(0, 0x0000ffff))
	a.offset = 0x0000ffff

	for i := 0; i < b.N; i++ {
		a.PutNum(uint32(i))
	}
}

func BenchmarkIPPool_PutGetNum(b *testing.B) {
	a, _ := NewIPPool(*MustNewIPRange(0, 0x0000ffff))
	a.offset = 0x0000ffff

	for i := 0; i < b.N; i++ {
		a.PutNum(uint32(i))
	}

	b.Log("ip count after put", a.totalBinIPCount)

	for i := 0; i < b.N; i++ {
		a.GetNum()
	}
}

func TestPlay(t *testing.T) {
	a := make([]byte, 10)

	a[0] = 0b00000001
	a[0] |= 0b10000000
	fmt.Println(a, a[0])
}
