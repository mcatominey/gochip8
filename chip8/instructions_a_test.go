package chip8

import (
	"testing"
)

func Test0xAnnn(t *testing.T) {
	c := New([]byte{
		0xAF, 0xFF,
	})

	c.Step()

	if c.i != 0xFFF {
		t.Error("i was not set to 0xFFF as expected")
	}
}
