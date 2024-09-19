package iputils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBitmap(t *testing.T) {
	b := newBitmap(8)
	assert.NotNil(t, b)
}

func TestBitmap_SetBit_UnsetBit_GetBit(t *testing.T) {
	b := newBitmap(8)
	b.SetBit(0)
	assert.True(t, b.GetBit(0))
	b.UnsetBit(0)
	assert.False(t, b.GetBit(0))

	b.SetBit(7)
	assert.True(t, b.GetBit(7))
	b.UnsetBit(7)
	assert.False(t, b.GetBit(7))

	b.SetBit(8)
	assert.False(t, b.GetBit(8))
	b.UnsetBit(8)
	assert.False(t, b.GetBit(8))
}

func TestBitmap_Ones(t *testing.T) {
	b := newBitmap(8)
	assert.Equal(t, b.Ones(), uint32(0))

	b.SetBit(0)
	assert.Equal(t, b.Ones(), uint32(1))

	b.SetBit(7)
	assert.Equal(t, b.Ones(), uint32(2))

	b.UnsetBit(0)
	assert.Equal(t, b.Ones(), uint32(1))

	b.UnsetBit(7)
	assert.Equal(t, b.Ones(), uint32(0))
}

func TestBitmap_FirstFlagOffset(t *testing.T) {
	b := newBitmap(20)
	assert.Equal(t, b.FirstFlagOffset(), -1)

	b.SetBit(0)
	assert.Equal(t, b.FirstFlagOffset(), 0)

	b.SetBit(7)
	assert.Equal(t, b.FirstFlagOffset(), 0)

	b.SetBit(8)
	assert.Equal(t, b.FirstFlagOffset(), 0)

	b.SetBit(19)
	assert.Equal(t, b.FirstFlagOffset(), 0)

	b.UnsetBit(0)
	assert.Equal(t, b.FirstFlagOffset(), 7)

	b.UnsetBit(7)
	assert.Equal(t, b.FirstFlagOffset(), 8)

	b.UnsetBit(8)
	assert.Equal(t, b.FirstFlagOffset(), 19)

	b.UnsetBit(19)
	assert.Equal(t, b.FirstFlagOffset(), -1)
}

func BenchmarkBitmap_FirstFlagOffset(b *testing.B) {
	bitmap := newBitmap(1000)
	bitmap.SetBit(999)
	for i := 0; i < b.N; i++ {
		bitmap.FirstFlagOffset()
	}
}
