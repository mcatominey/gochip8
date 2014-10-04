package chip8

import (
	"testing"
)

func Test0x1nnn(t *testing.T) {
	c := New([]byte{
		0x1F, 0xED,
	})

	c.Step()

	if c.pc != 0x0FED {
		t.Error("expected pc to be set to 0x0FED")
	}
}
