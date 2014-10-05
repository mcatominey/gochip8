package chip8

type Instruction struct {
	// original Opcode that was decoded
	Opcode uint16

	// explanation of decoded instruction
	Description string

	// Execute the Instruction
	// pc should be incremented before calling to ensure any Jump/Calls are correct
	implementation func(op uint16, c *Chip8)
}

var (
	unknown = Instruction{
		0xFFFF,
		"Unknown opcode",
		func(op uint16, c *Chip8) {
		},
	}
)

func DecodeOpcode(op uint16) Instruction {
	// First narrow down by first nibble
	switch op >> 12 {
	case 0x0:
		// 0x0???
		return decodeOpcode0(op)
	case 0x1:
		// 0x1???
		return decodeOpcode1(op)
	case 0x2:
		// 0x2???
		return decodeOpcode2(op)
	case 0x3:
		// 0x3???
		return decodeOpcode3(op)
	case 0x4:
		// 0x4???
		return decodeOpcode4(op)
	case 0x5:
		// 0x5???
		return decodeOpcode5(op)
	case 0x6:
		// 0x6???
		return decodeOpcode6(op)
	case 0x7:
		//0x7???
		return decodeOpcode7(op)
	case 0x8:
		//0x8???
		return decodeOpcode8(op)
	case 0x9:
		//0x9???
		return decodeOpcode9(op)
	case 0xA:
		//0xA???
		return decodeOpcodeA(op)
	case 0xB:
		//0xB???
		return decodeOpcodeB(op)
	case 0xC:
		//0xC???
		return decodeOpcodeC(op)
	case 0xD:
		//0xD???
		return decodeOpcodeD(op)
	case 0xE:
		//0xE???
		return decodeOpcodeE(op)
	case 0xF:
		//0xF???
		return decodeOpcodeF(op)
	}

	return unknown
}

func decodeOpcode0(op uint16) Instruction {
	if op == 0xE0 {
		return Instruction{
			op,
			"Clear the screen",
			func(op uint16, c *Chip8) {
				for x := 0; x < DisplayWidth; x++ {
					for y := 0; y < DisplayHeight; y++ {
						c.display[x][y] = 0
					}
				}
				c.DrawFlag = true
			},
		}
	} else if op == 0xEE {
		return Instruction{
			op,
			"Return from a subroutine",
			func(op uint16, c *Chip8) {
				c.sp--
				c.pc = c.stack[c.sp]
			},
		}
	} else {
		return Instruction{
			op,
			"Jump to subroutine in lowest 12 bits [IGNORED]",
			func(op uint16, c *Chip8) {
				// According to the reference this should be ignored by interpreters
			},
		}
	}
}

func decodeOpcode1(op uint16) Instruction {
	// Only 1 opcode where hightest nibble is 1
	return Instruction{
		op,
		"JUMP to location in lowest 12 bits",
		func(op uint16, c *Chip8) {
			// Set PC to lowest 12 bits
			c.pc = op & 0x0FFF
		},
	}
}

func decodeOpcode2(op uint16) Instruction {
	// Only 1 opcode where highest nibble is 2
	return Instruction{
		op,
		"CALL subroutine in lowest 12 bits",
		func(op uint16, c *Chip8) {
			c.stack[c.sp] = c.pc
			c.sp++
			c.pc = op & 0x0FFF
		},
	}
}

func decodeOpcode3(op uint16) Instruction {
	// Only 1 opcode where highest nibble is 3
	return Instruction{
		op,
		"Skip next instruction if Vx == kk when Opcode is 0x3xkk",
		func(op uint16, c *Chip8) {
			// If value in specified V register is equal to value in lowest byte
			// skip next instruction by incrementing PC by 2
			if c.v[getX(op)] == byte(op&0xFF) {
				c.pc += 2
			}
		},
	}
}

func decodeOpcode4(op uint16) Instruction {
	// Only 1 opcode where highest nibble is 4
	return Instruction{
		op,
		"Skip next instruction is Vx != kk when Opcode is 0x4xkk",
		func(op uint16, c *Chip8) {
			// If value in specified V register is not equal to value in lowest byte
			// skip next instruction by incrementing PC by 2
			if c.v[getX(op)] != byte(op&0xFF) {
				c.pc += 2
			}
		},
	}
}

func decodeOpcode5(op uint16) Instruction {
	// Only 1 opcode where highest nibble is 5
	return Instruction{
		op,
		"Skip next instruction if Vx == Vy when Opcode is 0x5xy0",
		func(op uint16, c *Chip8) {
			// If value in Vx is equal to value in Vy
			// skip next instruction by incrementing PC by 2
			if c.v[getX(op)] == c.v[getY(op)] {
				c.pc += 2
			}
		},
	}
}

func decodeOpcode6(op uint16) Instruction {
	// Only 1 opcode where highest nibble is 6
	return Instruction{
		op,
		"LOAD kk into Vx when Opcode is 0x6xkk",
		func(op uint16, c *Chip8) {
			// Load value kk into Vx
			c.v[getX(op)] = byte(op & 0xFF)
		},
	}
}

func decodeOpcode7(op uint16) Instruction {
	// Only 1 opcode where highest nibble is 7
	return Instruction{
		op,
		"ADD Vx to kk, store result in Vx when Opcode is 0x7xkk",
		func(op uint16, c *Chip8) {
			i := getX(op)
			c.v[i] = c.v[i] + byte(op&0xFF)
		},
	}
}

func decodeOpcode8(op uint16) Instruction {
	// Opcodes with highest nibble 8 are identified by the lowest nibble
	switch op & 0xF {
	case 0x0:
		return Instruction{
			op,
			"Stores the value of Vy in Vx",
			func(op uint16, c *Chip8) {
				c.v[getX(op)] = c.v[getY(op)]
			},
		}
	case 0x1:
		return Instruction{
			op,
			"Store bitwise OR result of Vx and Vy in Vx",
			func(op uint16, c *Chip8) {
				c.v[getX(op)] = c.v[getX(op)] | c.v[getY(op)]
			},
		}
	case 0x2:
		return Instruction{
			op,
			"Store bitwise AND result of Vx and Vy in Vx",
			func(op uint16, c *Chip8) {
				c.v[getX(op)] = c.v[getX(op)] & c.v[getY(op)]
			},
		}
	case 0x3:
		return Instruction{
			op,
			"Store bitwise XOR result of Vx and Vy in Vx",
			func(op uint16, c *Chip8) {
				c.v[getX(op)] = c.v[getX(op)] ^ c.v[getY(op)]
			},
		}
	case 0x4:
		return Instruction{
			op,
			"Add Vx and Vy, set VF (flag) if result carries (> 255) lowest 8 bits stored in Vx",
			func(op uint16, c *Chip8) {
				// Perform sum
				var result uint16 = uint16(c.v[getX(op)]) + uint16(c.v[op>>4&0xF])
				// Store only lowest 8 bits
				c.v[getX(op)] = byte(result & 0xFF)
				// Set VF for carry
				if result > 255 {
					c.v[0xF] = 1
				} else {
					c.v[0xF] = 0
				}
			},
		}
	case 0x5:
		return Instruction{
			op,
			"Subtract Vy from Vx, store result in Vx, VF = Vx > Vy ? 1 : 0",
			func(op uint16, c *Chip8) {
				// Set VF
				if c.v[getX(op)] > c.v[getY(op)] {
					c.v[0xF] = 1
				} else {
					c.v[0xF] = 0
				}
				// Subtract
				c.v[getX(op)] = c.v[getX(op)] - c.v[getY(op)]
			},
		}
	case 0x6:
		return Instruction{
			op,
			"Divide Vx by 2, if LSB of Vx is 1 set VF to 1 else 0",
			func(op uint16, c *Chip8) {
				// Set VF
				if c.v[getX(op)]&1 == 1 {
					c.v[0xF] = 1
				} else {
					c.v[0xF] = 0
				}
				// Half Vx by shifting right one
				c.v[getX(op)] >>= 1
			},
		}
	case 0x7:
		return Instruction{
			op,
			"Subtract Vx from Vy, store result in Vx, VF = Vy > Vx ? 1 : 0",
			func(op uint16, c *Chip8) {
				// Set VF
				if c.v[getY(op)] > c.v[getX(op)] {
					c.v[0xF] = 1
				} else {
					c.v[0xF] = 0
				}
				// Subtract
				c.v[getX(op)] = c.v[getY(op)] - c.v[getX(op)]
			},
		}
	case 0xE:
		// return unknown
		return Instruction{
			op,
			"Multiply Vx by 2. If MSB of Vx is 1 set VF to 1 else 0",
			func(op uint16, c *Chip8) {
				if ((c.v[getX(op)] & (1 << 7)) >> 7) == 1 {
					c.v[0xF] = 1
				} else {
					c.v[0xF] = 0
				}
				// Multiply
				c.v[getX(op)] *= 2
			},
		}
	}

	return unknown
}

func decodeOpcode9(op uint16) Instruction {
	// Only 1 Opcode where highest nibble is 9
	return Instruction{
		op,
		"Skip next instruction if Vx != Vy",
		func(op uint16, c *Chip8) {
			if c.v[getX(op)] != c.v[getY(op)] {
				c.pc += 2
			}
		},
	}
}

func decodeOpcodeA(op uint16) Instruction {
	// Only 1 Opcode where highest nibble is A
	return Instruction{
		op,
		"Set I to nnn (0xAnnn)",
		func(op uint16, c *Chip8) {
			c.i = op & 0xFFF
		},
	}
}

func decodeOpcodeB(op uint16) Instruction {
	// Only 1 Opcode where highest nibble is B
	return Instruction{
		op,
		"JUMP to location nnn + V0",
		func(op uint16, c *Chip8) {
			c.pc = op&0xFFF + uint16(c.v[0])
		},
	}
}

func decodeOpcodeC(op uint16) Instruction {
	// Only 1 Opcode where highest nibble is C
	return Instruction{
		op,
		"Bitwise AND result of Vx and Rand(0-255) stored in Vx",
		func(op uint16, c *Chip8) {
			rnd := c.rand.Byte()
			c.v[getX(op)] = rnd & byte(op&0xFF)
		},
	}
}

func decodeOpcodeD(op uint16) Instruction {
	// Only 1 Opcode where highest nibble is D
	return Instruction{
		op,
		"Draw to screen (too long to describe)",
		func(op uint16, c *Chip8) {
			startX := uint16(c.v[getX(op)])
			startY := uint16(c.v[getY(op)])

			// Set VF to 0
			c.v[0xF] = 0
			// Read sprite from memory
			var x, y uint16
			for y = 0; y < op&0xF; y++ {
				line := c.memory[c.i+y]

				// Sprites are 8 wide
				for x = 0; x < 8; x++ {
					xCoord := startX + x
					yCoord := startY + y

					// Use modulo to wrap if necessary
					if xCoord >= DisplayWidth {
						xCoord %= (DisplayWidth - 1)
					}
					if yCoord >= DisplayHeight {
						yCoord %= (DisplayHeight - 1)
					}

					if line&(1<<(7-x)) != 0 {
						if c.display[xCoord][yCoord] == 1 {
							// Coliision, set VF
							c.v[0xF] = 1
						}

						c.display[xCoord][yCoord] ^= 1
					}
				}
			}

			c.DrawFlag = true
		},
	}
}

func decodeOpcodeE(op uint16) Instruction {
	// Opcodes with highest nibble E are identified by the lowest byte
	switch op & 0xFF {
	case 0x9E:
		return Instruction{
			op,
			"Skip next instruction if if key with value in Vx is pressed",
			func(op uint16, c *Chip8) {
				if c.keys[c.v[getX(op)]] {
					c.pc += 2
				}
			},
		}
	case 0xA1:
		return Instruction{
			op,
			"Skip next instruction if if key with value in Vx is NOT pressed",
			func(op uint16, c *Chip8) {
				if !c.keys[c.v[getX(op)]] {
					c.pc += 2
				}
			},
		}
	}

	return unknown
}

func decodeOpcodeF(op uint16) Instruction {
	// Opcodes with highest nibble F are identified by the lowest byte
	switch op & 0xFF {
	case 0x07:
		return Instruction{
			op,
			"Set Vx to value of delay timer",
			func(op uint16, c *Chip8) {
				c.v[getX(op)] = c.delay
			},
		}
	case 0x0A:
		return Instruction{
			op,
			"Wait for a key press, store key value in Vx",
			func(op uint16, c *Chip8) {
				c.waitingForKey = true
				c.waitingKeyRegister = int8(getX(op))
			},
		}
	case 0x15:
		return Instruction{
			op,
			"Set delay timer to Vx",
			func(op uint16, c *Chip8) {
				c.delay = c.v[getX(op)]
			},
		}
	case 0x18:
		return Instruction{
			op,
			"Set sound timer to Vx",
			func(op uint16, c *Chip8) {
				c.sound = c.v[getX(op)]
			},
		}
	case 0x1E:
		return Instruction{
			op,
			"Add I and Vx, store result in I",
			func(op uint16, c *Chip8) {
				c.i = uint16(c.v[getX(op)]) + c.i
			},
		}
	case 0x29:
		return Instruction{
			op,
			"Set I to memory address of the sprite data for character in VX",
			func(op uint16, c *Chip8) {
				// 5 bytes per character, starting at address 0x000
				// makes calculating address simple
				c.i = 5 * uint16(c.v[getX(op)])
			},
		}
	case 0x33:
		return Instruction{
			op,
			"Store BCD of Vx in memory, hundreds at I, tens at I+1, ones at I+2",
			func(op uint16, c *Chip8) {
				val := c.v[getX(op)]
				c.memory[c.i] = (val / 100) % 10  // hundreds
				c.memory[c.i+1] = (val / 10) % 10 // tens
				c.memory[c.i+2] = val % 10        // ones
			},
		}
	case 0x55:
		return Instruction{
			op,
			"Store V0 to Vx in memory starting at location I",
			func(op uint16, c *Chip8) {
				var reg uint16
				end := uint16(getX(op))
				for reg = 0; reg <= end; reg++ {
					c.memory[c.i+reg] = c.v[reg]
				}
			},
		}
	case 0x65:
		return Instruction{
			op,
			"Load into V0 to Vx from memory starting at location I",
			func(op uint16, c *Chip8) {
				var reg uint16
				end := uint16(getX(op))
				for reg = 0; reg <= end; reg++ {
					c.v[reg] = c.memory[c.i+reg]
				}
			},
		}
	}

	return unknown
}
