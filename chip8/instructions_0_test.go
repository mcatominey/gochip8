package chip8

import (
	"testing"
)

func Test0x00E0(t *testing.T) {
	c := New([]byte{
		0x00, 0xE0,
	})

	// Set all pixels to on
	for x := 0; x < DisplayWidth; x++ {
		for y := 0; y < DisplayHeight; y++ {
			c.display[x][y] = 1
		}
	}

	// Run one instruction
	c.Step()

	// Draw flag should be set
	if !c.DrawFlag {
		t.Error("DrawFlag is false, expected true after call to clear")
	}

	// Check screen has been cleared
	for x := 0; x < DisplayWidth; x++ {
		for y := 0; y < DisplayHeight; y++ {
			if c.display[x][y] != 0 {
				// Fatal since it would print for every pixel otherwise
				t.Fatalf("Pixel at x:%d y:%d is not 0", x, y)

			}
		}
	}
}

func Test0x00EE(t *testing.T) {
	var stackAddressTest uint16 = 0xBEEF

	c := New([]byte{
		0x00, 0xEE,
	})

	// Add address to stack
	c.stack[c.sp] = stackAddressTest
	c.sp++

	c.Step()

	// PC is incremented in Step for next instruction
	if c.pc != stackAddressTest {
		t.Errorf("Expected PC to be set to %#X, actually %#X", stackAddressTest, c.pc)
	}
}

func Test0x0nnn(t *testing.T) {
	c := New([]byte{
		0x01, 0x11,
	})

	c.Step()

	// Ensure PC did not Jump as this instruction is ignored
	if c.pc != programStartAddress+2 {
		t.Error("expected instruction to be ignored and have no effect")
	}
}
