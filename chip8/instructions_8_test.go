package chip8

import (
	"testing"
)

func Test0x8xy0(t *testing.T) {
	c := New([]byte{
		0x80, 0x10,
	})

	c.v[1] = 0xFF

	c.Step()

	if c.v[1] != 0xFF {
		t.Error("V1 did not contain expected value")
	}
}

func Test0x8xy1(t *testing.T) {
	c := New([]byte{
		0x80, 0x11,
	})

	c.v[0] = 0x0F
	c.v[1] = 0xF0

	c.Step()

	if c.v[0] != (0x0F | 0xF0) {
		t.Error("V0 did not contain correct Bitwise OR result")
	}
}

func Test0x8xy2(t *testing.T) {
	c := New([]byte{
		0x80, 0x12,
	})

	c.v[0] = 0x0F
	c.v[1] = 0xF0

	c.Step()

	if c.v[0] != (0x0F & 0xF0) {
		t.Error("V0 did not contain correct Bitwise AND result")
	}
}

func Test0x8xy3(t *testing.T) {
	c := New([]byte{
		0x80, 0x13,
	})

	c.v[0] = 0x0F
	c.v[1] = 0xF0

	c.Step()

	if c.v[0] != (0x0F ^ 0xF0) {
		t.Error("V0 did not contain correct Bitwise XOR result")
	}
}

func Test0x8xy4(t *testing.T) {
	c := New([]byte{
		0x80, 0x14,
	})

	c.v[0] = 0x0F
	c.v[1] = 0xF0

	c.Step()

	if c.v[0] != (0x0F + 0xF0) {
		t.Error("V0 did not contain correct result of ADD")
	}
}

func Test0x8xy4Carry(t *testing.T) {
	c := New([]byte{
		0x80, 0x14,
	})

	c.v[0] = 0xFF
	c.v[1] = 0x01

	c.Step()

	// Check carry is set
	if c.v[0xF] != 1 {
		t.Error("VF was not set to 1 as expected after carry")
	}
}

func Test0x8xy4NoCarry(t *testing.T) {
	c := New([]byte{
		0x80, 0x14,
	})

	c.v[0] = 0xFF
	c.v[1] = 0x00

	c.Step()

	// Check carry is set
	if c.v[0xF] != 0 {
		t.Error("VF was not set to 0 as expected after no carry")
	}
}

func Test0x8xy5NotBorrow(t *testing.T) {
	c := New([]byte{
		0x80, 0x15,
	})

	c.v[0] = 0x0F // Vx
	c.v[1] = 0x0A // Vy

	c.Step()

	// Check VF
	if c.v[0xF] != 1 {
		t.Error("VF should be 1")
	}

	if c.v[0] != (0x0F - 0x0A) {
		t.Error("subtraction result is incorrect")
	}
}

func Test0x8xy5Borrow(t *testing.T) {
	c := New([]byte{
		0x80, 0x15,
	})

	c.v[0] = 0x0A // Vx
	c.v[1] = 0x0F // Vy
	c.v[0xF] = 1  // Force VF to 1 to check for 0 later

	// introduce variable since constants do not wrap around
	expectedResult := c.v[0] - c.v[1]

	c.Step()

	if c.v[0xF] != 0 {
		t.Error("VF should be 0")
	}

	if c.v[0] != expectedResult {
		t.Error("subtraction result is incorrect")
	}
}

func Test0x8xy6LSB0(t *testing.T) {
	c := New([]byte{
		0x80, 0x16,
	})

	c.v[0] = 0xAA
	c.v[0xF] = 1 // Force VF to 1 to check for 0 later

	c.Step()

	// VF should be 0 since LSB is 0
	if c.v[0xF] != 0 {
		t.Error("VF should be 0")
	}

	if c.v[0] != (0xAA / 2) {
		t.Error("division (>> 1) result is incorrect")
	}
}

func Test0x8xy6LSB1(t *testing.T) {
	c := New([]byte{
		0x80, 0x16,
	})

	c.v[0] = 0xFF // Vx

	c.Step()

	// VF should be 1 since LSB is 1
	if c.v[0xF] != 1 {
		t.Error("VF should be 1")
	}

	if c.v[0] != (0xFF / 2) {
		t.Error("division (>> 1) result is incorrect")
	}
}

func Test0x8xy7NotBorrow(t *testing.T) {
	c := New([]byte{
		0x80, 0x17,
	})

	c.v[0] = 0x0A // Vx
	c.v[1] = 0x0B // Vy

	c.Step()

	if c.v[0xF] != 1 {
		t.Error("VF should be 1")
	}

	if c.v[0] != (0x0B - 0x0A) {
		t.Error("subtraction result is incorrect")
	}
}

func Test0x8xy7Borrow(t *testing.T) {
	c := New([]byte{
		0x80, 0x17,
	})

	c.v[0] = 0x0B // Vx
	c.v[1] = 0x0A // Vy
	c.v[0xF] = 1  // Force VF to 1 to check for 0 later

	// introduce variable since constants do not wrap around
	expectedResult := c.v[1] - c.v[0]

	c.Step()

	if c.v[0xF] != 0 {
		t.Error("VF should be 0")
	}

	if c.v[0] != expectedResult {
		t.Error("subtraction result is incorrect")
	}
}

func Test0x8xyEMSB0(t *testing.T) {
	c := New([]byte{
		0x80, 0x1E,
	})

	c.v[0] = 0x7A
	c.v[0xF] = 1 // Force VF to 1 to check for 0 later

	c.Step()

	// VF should be 0 since MSB is 0
	if c.v[0xF] != 0 {
		t.Error("VF should be 0")
	}

	if c.v[0] != (0x7A << 1) {
		t.Error("multiplication (<< 1) result is incorrect")
	}
}

func Test0x8xyEMSB1(t *testing.T) {
	c := New([]byte{
		0x80, 0x1E,
	})

	c.v[0] = 0xF0 // Vx
	// introduce variable since constants do not wrap around
	expectedResult := 0xF0 << 1

	c.Step()

	// VF should be 1 since MSB is 1
	if c.v[0xF] != 1 {
		t.Error("VF should be 1")
	}

	if c.v[0] != byte(expectedResult) {
		t.Error("multiplication (<< 1) result is incorrect")
	}
}
