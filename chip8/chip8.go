package chip8

const (
	MemorySize     = 4096
	StackSize      = 16
	VRegisterCount = 16
	KeyCount       = 16

	// Original Chip 8 Resolution
	DisplayWidth  = 64
	DisplayHeight = 32

	// Starting memory address where fonts are loaded
	fontStartAddress = 0x000

	// Starting memory address where roms are loaded
	programStartAddress = 0x200 // (512)
)

var (
	// Sprite data for hexidecimal character set
	fontData = [...]byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
)

type Chip8 struct {
	// memory
	// 0x000 - 0x050 => Fonts
	// 0x051 - 0x1FF => Reserved (Not used)
	// 0x200 - 0xFFF => Program and General memory
	memory [MemorySize]byte

	// Program Counter and Stack Pointer
	pc, sp uint16

	// Stack
	stack [StackSize]uint16

	// V Registers
	v [VRegisterCount]byte

	// I Register (16 bits wide)
	i uint16

	// display, on or off
	display [DisplayWidth][DisplayHeight]byte

	DrawFlag bool

	// Keyboard, down = true
	keys [KeyCount]bool
	// true if waiting for a key
	waitingForKey bool
	// index of V register where value waiting key should be stored
	// -1 while not waiting or not set
	waitingKeyRegister int8

	// Timers 60Hz counting down
	sound, delay byte

	rand RandSource
}

// New creates a new Chip 8, reading the program file from file at filename
func New(program []byte) *Chip8 {
	c8 := &Chip8{
		rand: Rand{},
	}

	c8.Reset()

	c8.LoadProgram(program)

	return c8
}

// Clear the state of the Chip8
func (c *Chip8) Reset() {
	// CLear memory
	for i := 0; i < len(c.memory); i++ {
		c.memory[i] = 0
	}

	// Clear registers
	c.pc = programStartAddress
	c.sp = 0
	c.i = 0
	for i := 0; i < len(c.v); i++ {
		c.v[i] = 0
	}

	// Clear stack
	for i := 0; i < len(c.stack); i++ {
		c.stack[i] = 0
	}

	// Clear display
	for x := 0; x < DisplayWidth; x++ {
		for y := 0; y < DisplayHeight; y++ {
			c.display[x][y] = 0
		}
	}
	c.DrawFlag = false

	// Keyboard
	for i := 0; i < len(c.keys); i++ {
		c.keys[i] = false
	}
	c.waitingKeyRegister = -1
	c.waitingForKey = false

	// Timers
	c.delay = 0
	c.sound = 0
}

func (c *Chip8) LoadProgram(program []byte) {
	// Load program into memory
	for i := 0; i < len(program); i++ {
		c.memory[programStartAddress+i] = program[i]
	}
	// Load font into memory
	for i := 0; i < len(fontData); i++ {
		c.memory[fontStartAddress+i] = fontData[i]
	}
}

func (c *Chip8) PressKey(key Key) {
	// Only Hex keys
	if key > 0xF {
		return
	}

	if c.waitingForKey {
		c.waitingForKey = false
		c.v[c.waitingKeyRegister] = byte(key)
		c.waitingKeyRegister = -1
	}

	c.keys[key] = true
}

func (c *Chip8) DePressKey(key Key) {
	c.keys[key] = false
}

// Step emulates the execution of a single instruction
// returns true if an instruction was actually executed
func (c *Chip8) Step() bool {
	// Check if PC is in bounds
	if c.pc < programStartAddress || c.pc > MemorySize-1 {
		panic("pc is at invalid address")
	}

	// Check if waiting for a key
	if c.waitingForKey {
		return false
	}

	// Get the opcode from program memory
	op := getOpcode(c.memory[c.pc], c.memory[c.pc+1])

	c.pc += 2

	instruction := DecodeOpcode(op)

	instruction.implementation(op, c)

	return true
}

// UpdateTimers will decrement both the sound and delay timers given
// that the Chip 8 is not currently waiting for a key press
func (c *Chip8) UpdateTimers() {
	// Check if execution is paused for key press
	if c.waitingForKey {
		return
	}

	if c.delay > 0 {
		c.delay--
	}
	if c.sound > 0 {
		c.sound--
	}
}

func (c *Chip8) GetDisplay() [DisplayWidth][DisplayHeight]byte {
	return c.display
}
