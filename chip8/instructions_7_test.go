package chip8

import (
	"testing"
)

func Test0x7xkk(t *testing.T) {
	c := New([]byte{
		0x70, 0x02,
	})
	c.v[0] = 0x02

	c.Step()

	if c.v[0] != (0x02 + 0x02) {
		t.Error("V0 does not contain correct sum")
	}
}
