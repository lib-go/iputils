package iputils

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func mustNewIPRange(begin, end interface{}) (r *IPRange) {
	var e error
	r, e = NewIPRange(begin, end)
	if e != nil {
		panic(e)
	}
	return
}

func TestNewIPRange(t *testing.T) {
	r, e := NewIPRange("192.168.1.1", "192.168.31.125")
	t.Log(r, e)

	assert.Equal(t, r.At(0), net.IPv4(192, 168, 1, 1).To4())
	assert.Equal(t, r.At(-1), net.IPv4(192, 168, 31, 125).To4())
}

func TestNewIPRangeFromIPNet(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.31.0/24")

	t.Log(ipnet)

	r, e := NewIPRangeFromIPNet(ipnet)
	t.Log(r, e)

	assert.Equal(t, r.At(0), net.IPv4(192, 168, 31, 0).To4())
	assert.Equal(t, r.At(-1), net.IPv4(192, 168, 31, 255).To4())
}

func TestIPRange_Size(t *testing.T) {
	r := mustNewIPRange("10.0.0.1", "10.0.0.255")
	assert.Equal(t, r.Size(), uint32(255))
}

func TestIPRange_BeginNum(t *testing.T) {
	r := mustNewIPRange("0.0.0.1", "0.0.0.255")
	assert.Equal(t, r.BeginNum(), uint32(1))
}

func TestIPRange_EndNum(t *testing.T) {
	r := mustNewIPRange("0.0.0.1", "0.0.0.255")
	assert.Equal(t, r.EndNum(), uint32(255))
}

func TestNewIPRangeFromCIDR(t *testing.T) {
	r, e := NewIPRangeFromCIDR("192.168.31.123/24")
	t.Log(r, e)
	assert.Equal(t, r.At(0), net.IPv4(192, 168, 31, 0).To4())
	assert.Equal(t, r.At(-1), net.IPv4(192, 168, 31, 255).To4())
}

func TestIPRange_Has(t *testing.T) {
	r := mustNewIPRange("10.0.0.1", "10.0.0.255")
	assert.True(t, r.Has(net.IPv4(10, 0, 0, 1).To4()))
	assert.True(t, r.Has(net.IPv4(10, 0, 0, 255).To4()))
	assert.False(t, r.Has(net.IPv4(10, 0, 0, 0).To4()))
	assert.False(t, r.Has(net.IPv4(10, 0, 1, 0).To4()))

	r = mustNewIPRange("10.0.0.1", "10.1.0.255")
	assert.True(t, r.Has(net.IPv4(10, 0, 0, 1).To4()))
	assert.True(t, r.Has(net.IPv4(10, 1, 0, 255).To4()))
	assert.True(t, r.Has(net.IPv4(10, 0, 1, 255).To4()))
	assert.False(t, r.Has(net.IPv4(10, 0, 0, 0).To4()))
	assert.False(t, r.Has(net.IPv4(10, 1, 1, 0).To4()))
}

func TestIPRange_HasOverlap(t *testing.T) {
	r1 := mustNewIPRange("10.0.1.1", "10.0.1.100")

	// 头相等，尾在内
	r2 := mustNewIPRange("10.0.1.1", "10.0.1.10")
	assert.True(t, r1.HasOverlap(r2))

	// 头在内, 尾相等
	r2 = mustNewIPRange("10.0.1.90", "10.0.1.100")
	assert.True(t, r1.HasOverlap(r2))

	// 头尾都在内
	r2 = mustNewIPRange("10.0.1.50", "10.0.1.60")
	assert.True(t, r1.HasOverlap(r2))

	// 头在外, 尾在内
	r2 = mustNewIPRange("10.0.0.1", "10.0.1.50")
	assert.True(t, r1.HasOverlap(r2))

	// 头在内，尾出界
	r2 = mustNewIPRange("10.0.1.50", "10.0.2.1")
	assert.True(t, r1.HasOverlap(r2))

	// 头尾都出界
	r2 = mustNewIPRange("10.0.0.1", "10.0.2.1")
	assert.True(t, r1.HasOverlap(r2))

	// r2 < r1 不沾
	r2 = mustNewIPRange("10.0.0.1", "10.0.0.100")
	assert.False(t, r1.HasOverlap(r2))

	// r2 > r1 不沾
	r2 = mustNewIPRange("10.0.2.1", "10.0.2.100")
	assert.False(t, r1.HasOverlap(r2))
}

func TestIPRange_Split(t *testing.T) {
	r := mustNewIPRange("10.0.0.1", "10.0.0.100")
	r1, r2, e := r.Split(10)
	t.Log(r1, r2, e)
	assert.Equal(t, r1.At(0), net.IPv4(10, 0, 0, 1).To4())
	assert.Equal(t, r1.At(-1), net.IPv4(10, 0, 0, 10).To4())
	assert.Equal(t, r2.At(0), net.IPv4(10, 0, 0, 11).To4())
	assert.Equal(t, r2.At(-1), net.IPv4(10, 0, 0, 100).To4())
}

func TestIPRange_TrimAndCopy(t *testing.T) {
	r, _ := NewIPRange("10.0.0.1", "10.0.0.100")
	r2, e := r.TrimCopy(10, 10)
	t.Log(r2, e)

	assert.Equal(t, r2.At(0), net.IPv4(10, 0, 0, 11).To4())
	assert.Equal(t, r2.At(-1), net.IPv4(10, 0, 0, 90).To4())
}
