package chip8

import (
	"math/rand"
	"testing"
)

func Test0xFx07(t *testing.T) {
	c := New([]byte{
		0xF0, 0x07,
	})
	c.delay = 0xFF

	c.Step()

	if c.v[0] != c.delay {
		t.Error("V0 not set to delay timer value")
	}
}

func Test0xFx0A(t *testing.T) {
	t.Skip("Not implemented")
}

func Test0xFx15(t *testing.T) {
	c := New([]byte{
		0xF0, 0x15,
	})
	c.v[0] = 0xFF

	c.Step()

	if c.delay != c.v[0] {
		t.Error("delay timer not set to V0")
	}
}

func Test0xFx18(t *testing.T) {
	c := New([]byte{
		0xF0, 0x18,
	})
	c.v[0] = 0xFF

	c.Step()

	if c.sound != c.v[0] {
		t.Error("sound timer not set to V0")
	}
}

func Test0xFx1E(t *testing.T) {
	c := New([]byte{
		0xF0, 0x1E,
	})
	c.v[0] = 0x0A
	c.i = 0x0A

	c.Step()

	if c.i != (0x0A + 0x0A) {
		t.Error("i not set to correct sum result")
	}
}

func Test0xFx29(t *testing.T) {
	hex := byte(rand.Int31n(0xF))

	c := New([]byte{
		0xF5, 0x29, // Hex character in V5
	})
	c.v[0x5] = hex

	c.Step()

	if c.i != uint16(fontStartAddress+(hex*5)) {
		t.Error("incorrect font memory address in I")
	}
}

func Test0xFx33(t *testing.T) {
	var startMemoryAddress uint16 = 0x300
	var bcdTest byte = 251

	c := New([]byte{
		0xF0, 0x33,
	})

	c.i = startMemoryAddress
	c.v[0] = bcdTest // Vx

	c.Step()

	// Hundreds
	if c.memory[startMemoryAddress] != 2 {
		t.Error("BCD hundreds digit incorrect")
	}

	// Tens
	if c.memory[startMemoryAddress+1] != 5 {
		t.Error("BCD tens digit incorrect")
	}

	// Ones
	if c.memory[startMemoryAddress+2] != 1 {
		t.Error("BCD ones digit incorrect")
	}
}

var (
	memoryStartAddress uint16 = 0x300
	memoryValues              = []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
		0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B,
	}
)

func Test0xFx55(t *testing.T) {
	c := New([]byte{
		0xFB, 0x55,
	})

	// Setup register values
	for i := 0; i < len(memoryValues); i++ {
		c.v[i] = memoryValues[i]
	}
	c.i = memoryStartAddress

	c.Step()

	for i := 0; i < len(memoryValues); i++ {
		if c.memory[memoryStartAddress+uint16(i)] != c.v[i] {
			t.Errorf("memory at %#x does not contain correct value", memoryStartAddress+uint16(i))
		}
	}
}

func Test0xFx65(t *testing.T) {
	c := New([]byte{
		0xFB, 0x65,
	})

	// Setup register values
	for i := 0; i < len(memoryValues); i++ {
		c.memory[memoryStartAddress+uint16(i)] = memoryValues[i]
	}
	c.i = memoryStartAddress

	c.Step()

	for i := 0; i < len(memoryValues); i++ {
		if c.memory[memoryStartAddress+uint16(i)] != c.v[i] {
			t.Errorf("memory at %#x does not contain correct value", memoryStartAddress+uint16(i))
		}
	}
}
