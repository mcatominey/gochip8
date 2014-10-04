package chip8

import (
	"math/rand"
)

const (
	mockRandByte byte = 0xAB
)

type RandSource interface {
	Byte() byte
}

type Rand struct{}

func (r Rand) Byte() byte {
	return byte(rand.Int31n(256))
}

// MockRand implemented RandSource returning mockRandByte
// every time
type MockRand struct{}

func (m MockRand) Byte() byte {
	return mockRandByte
}
