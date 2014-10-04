package chip8

import (
	"testing"
)

func Test0xEx9E(t *testing.T) {
	c := New([]byte{
		0xE0, 0x9E,
		0x1F, 0xFF, // Jump, should be skipped
		0x00, 0xE0, // Clear screen
	})
	c.v[0] = byte(KeyF)

	c.PressKey(KeyF)
	c.Step()
	c.Step()

	// Check if instruction was not skipped
	if c.pc == 0xFFF {
		t.Error("instruction not skipped as expected")
	}
}

func Test0xExA1(t *testing.T) {
	c := New([]byte{
		0xE0, 0xA1,
		0x1F, 0xFF, // Jump, should be skipped
		0x00, 0xE0, // Clear screen
	})
	c.v[0] = byte(KeyF)

	c.PressKey(KeyF)
	c.DePressKey(KeyF)
	c.Step()
	c.Step()

	// Check if instruction was not skipped
	if c.pc == 0xFFF {
		t.Error("instruction not skipped as expected")
	}
}
