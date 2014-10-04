package chip8

import (
	"testing"
)

func Test0xBnnn(t *testing.T) {
	c := New([]byte{
		0xBA, 0xAA,
	})
	c.v[0] = 0xBB

	c.Step()

	if c.pc != (0xAAA + 0xBB) {
		t.Error("pc was not set to correct result")
	}
}
