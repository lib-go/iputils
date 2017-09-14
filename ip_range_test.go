package iputils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestIPRange(t *testing.T) {
	r, _ := NewIPRange(1, 2)

	// intersection
	r2 := &IPRange{beginNum: 10, endNum: 11}
	assert.False(t, r.HasIntersection(r2))
	r3 := &IPRange{beginNum: 2, endNum: 3}
	assert.True(t, r.HasIntersection(r3))
	fmt.Println(r, r2, r3)

	// InRange
	assert.True(t, r.Has(net.IPv4(0, 0, 0, 1).To4()))
	assert.False(t, r.Has(net.IPv4(0, 0, 0, 3).To4()))

	// pop / remove
	r4 := &IPRange{beginNum: 1, endNum: 10}
	ip, e := r4.PopLeft()
	assert.Equal(t, ip, net.IPv4(0, 0, 0, 1).To4())
	assert.Equal(t, r4.FirstIP(), net.IPv4(0, 0, 0, 2).To4())
	assert.Nil(t, e)

	ip, e = r4.PopRight()
	assert.Equal(t, ip, net.IPv4(0, 0, 0, 10).To4())
	assert.Equal(t, r4.LastIP(), net.IPv4(0, 0, 0, 9).To4())
	assert.Nil(t, e)

	r, e = r4.RemoveLeft(3)
	assert.Equal(t, r.FirstIP(), net.IPv4(0, 0, 0, 2).To4())
	assert.Equal(t, r.LastIP(), net.IPv4(0, 0, 0, 4).To4())
	assert.Equal(t, r4.FirstIP(), net.IPv4(0, 0, 0, 5).To4())
	assert.Nil(t, e)

	r, e = r4.RemoveRight(3)
	assert.Equal(t, r.FirstIP(), net.IPv4(0, 0, 0, 7).To4())
	assert.Equal(t, r.LastIP(), net.IPv4(0, 0, 0, 9).To4())
	assert.Equal(t, r4.LastIP(), net.IPv4(0, 0, 0, 6).To4())
	assert.Nil(t, e)
}

func TestIPRange_Exclude(t *testing.T) {
	r, _ := NewIPRange("192.168.31.0", "192.168.31.255")
	r2, e := r.Exclude(10, 10)
	t.Log(r2, e)

	assert.Equal(t, r2.FirstIP(), net.IPv4(192, 168, 31, 10).To4())
	assert.Equal(t, r2.LastIP(), net.IPv4(192, 168, 31, 245).To4())
}

func TestNewIPRangeFromIPNet(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.31.0/24")

	t.Log(ipnet)

	r, e := NewIPRangeFromIPNet(ipnet)
	t.Log(r, e)

	assert.Equal(t, r.FirstIP(), net.IPv4(192, 168, 31, 0).To4())
	assert.Equal(t, r.LastIP(), net.IPv4(192, 168, 31, 255).To4())
}

func TestNewIPRangeFromCIDR(t *testing.T) {
	r, e := NewIPRangeFromCIDR("192.168.31.123/24")
	t.Log(r, e)
	assert.Equal(t, r.FirstIP(), net.IPv4(192, 168, 31, 0).To4())
	assert.Equal(t, r.LastIP(), net.IPv4(192, 168, 31, 255).To4())
}
