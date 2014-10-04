package chip8

import (
	"testing"
)

// Skips next instruction if Vx == kk
func Test0x3xkk(t *testing.T) {
	c := New([]byte{
		0x30, 0xFF,
		0x1F, 0xFF, // Jump, should be skipped
		0x00, 0xE0, // Clear screen
	})
	c.v[0] = 0xFF

	c.Step()
	c.Step()

	// Check if instruction was skipped
	if c.pc == 0xFFF {
		t.Error("Test instruction was not skipped as expected")
	}
}
