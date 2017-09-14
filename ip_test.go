package iputils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestAsIPv4(t *testing.T) {
	assert.Equal(t, net.IPv4(192, 168, 31, 123).To4(), AsIPv4("192.168.31.123"))
	assert.Equal(t, net.IPv4(192, 168, 31, 123).To4(), AsIPv4(net.IPv4(192, 168, 31, 123)))
}

func TestIsIPv4(t *testing.T) {
	assert.True(t, IsIPv4("192.168.0.1"))
	assert.False(t, IsIPv4("baidu.com"))
	assert.False(t, IsIPv4("2001:0DB8:0000:0000:0000:0000:1428:0000"))
	assert.False(t, IsIPv4("l13hj41lkh34l1hj3"))
}

func TestIsPrivateIPv4(t *testing.T) {
	assert.True(t, IsPrivateIPv4("127.0.0.1"))
	assert.False(t, IsPrivateIPv4("128.0.0.1"))
	assert.True(t, IsPrivateIPv4("10.0.0.1"))
	assert.False(t, IsPrivateIPv4("11.0.0.0"))
	assert.True(t, IsPrivateIPv4("172.16.1.1"))
	assert.False(t, IsPrivateIPv4("172.32.0.0"))
	assert.True(t, IsPrivateIPv4("192.168.0.1"))
	assert.False(t, IsPrivateIPv4("192.169.0.0"))

	assert.True(t, IsPrivateIPv4(net.IPv4(127, 0, 0, 1)))
	assert.False(t, IsPrivateIPv4(net.IPv4(128, 0, 0, 1)))
}

func TestIPv4ToUint32(t *testing.T) {
	ip := net.IPv4(1, 1, 1, 1)
	n := IPv4ToUint32(ip)
	assert.Equal(t, 0x1010101, int(n))
}

func TestUint32ToIPv4(t *testing.T) {
	var n uint32 = 0x1010101
	ip := Uint32ToIPv4(n)
	assert.True(t, ip.Equal(net.IPv4(1, 1, 1, 1)))
}

func TestIPv4Inc(t *testing.T) {
	ip := net.IPv4(0, 0, 0, 0)

	ip2 := IPv4Inc(ip, 10)
	assert.True(t, ip2.Equal(net.IPv4(0, 0, 0, 10)))
}

func TestIPEnd(t *testing.T) {
	ip, ipnet, _ := net.ParseCIDR("172.244.0.1/13")
	ip = IPv4End(ipnet)
	fmt.Println(ip)
	assert.Equal(t, ip, net.IPv4(172, 247, 255, 255).To4())
}
