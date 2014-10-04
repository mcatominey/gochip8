package chip8

func getOpcode(left, right byte) uint16 {
	l := uint16(left)
	r := uint16(right)

	return l<<8 | r
}

func getX(opcode uint16) byte {
	return byte((opcode & 0x0F00) >> 8)
}

func getY(opcode uint16) byte {
	return byte((opcode & 0x00F0) >> 4)
}
