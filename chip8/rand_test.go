package chip8

import (
	"testing"
)

// This test is in no way comprehensive or well designed it only
// serves to ensure that Rand doesn't return the same result constantly
func TestRand(t *testing.T) {
	iterations := 100
	r := Rand{}

	last := r.Byte()
	for i := 0; i < iterations; i++ {
		if r.Byte() != last {
			return
		}
	}

	t.Errorf("random byte generator returned same value %d times", iterations)
}

// Same as above applies
func TestMockRand(t *testing.T) {
	iterations := 100
	r := MockRand{}

	last := r.Byte()
	for i := 0; i < iterations; i++ {
		if r.Byte() != last {
			t.Error("mock random byte generator returned differing value")
		}
	}
}
