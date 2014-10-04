package chip8

import (
	"testing"
)

func Test0x2nnn(t *testing.T) {
	c := New([]byte{
		0x2A, 0xAA,
	})

	c.Step()

	// Check for PC jump
	if c.pc != 0x0AAA {
		t.Error("expected pc to be set to 0x0AAA")
	}

	// Check pc was placed onto stack
	// Since there is only one instruction pc should be at start address + 2
	if c.stack[c.sp-1] != programStartAddress+2 {
		t.Error("pc was not added to stack")
	}
}
