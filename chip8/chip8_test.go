package chip8

import (
	"testing"
)

func TestPCInitialization(t *testing.T) {
	c8 := New([]byte{})

	if c8.pc != programStartAddress {
		t.Errorf("PC should be %#X at initialization", programStartAddress)
	}
}

func TestIncorrectKeyValue(t *testing.T) {
	c := New([]byte{})

	// Array bounds would panic if not caught by PressKey
	// Test by deferring recover
	defer func() {
		if err := recover(); err != nil {
			t.Error("PressKey caused panic with invalid key index")
		}
	}()

	c.PressKey(KeyF + 1)
}

func TestWaitingForKey(t *testing.T) {
	c := New([]byte{
		0xF2, 0x0A, // Wait for key opcode, store key pressed in V2
		0x00, 0xE0, // Clear Screen
		0x00, 0xE0, // Clear Screen
	})

	// Check not initially waiting for key
	if c.Step() == false {
		t.Error("step returned false when not waiting for key")
	}

	if c.waitingForKey != true {
		t.Error("expected waitingForKey to be true")
	}

	if c.Step() == true {
		t.Error("step returned true when waiting for key")
	}

	if c.waitingKeyRegister != 0x2 {
		t.Error("wrong register index stored in waitingKeyRegister")
	}

	c.PressKey(KeyA)

	if c.Step() == false {
		t.Error("step returned false when not waiting for key")
	}
}
