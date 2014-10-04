package chip8

import (
	"testing"
)

func Test0xCxkk(t *testing.T) {
	c := New([]byte{
		0xC0, 0xBB,
	})

	c.rand = MockRand{}

	c.Step()

	if c.v[0] != mockRandByte&0xBB {
		t.Error("incorrect result of random byte & lower op byte")
	}
}
