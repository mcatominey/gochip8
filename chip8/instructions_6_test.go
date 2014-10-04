package chip8

import (
	"testing"
)

func Test0x6xkk(t *testing.T) {
	c := New([]byte{
		0x60, 0xFF,
	})

	c.Step()

	if c.v[0] != 0xFF {
		t.Error("V0 was not set as expected")
	}
}
