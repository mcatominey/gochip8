package main

import (
	"flag"
	"github.com/mcatominey/gochip8/chip8"
	"github.com/veandco/go-sdl2/sdl"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Config
var (
	// If this flag is present, the rom shold be loaded then output to stdout line by line with
	// an explanation of each opcode
	disassemble = flag.Bool("disassemble", false, "disassemble to stdout")

	// Scale factor for upscaling the display from Chip 8 resolution of 64*32
	scaleFactor = flag.Int("scaling", 10, "scale factor to multiply Chip 8 resolution (64*32) by")

	// Controls execution speed, useful since some roms play at mad speeds compared to others
	stepsPerSecond = flag.Int64("sps", 500, "instructions to emulate per second")

	// Default key bindings
	defaultKeyBindings = map[sdl.Keycode]chip8.Key{
		sdl.K_1: chip8.Key1, sdl.K_2: chip8.Key2, sdl.K_3: chip8.Key3, sdl.K_4: chip8.KeyC, // 1 2 3 4
		sdl.K_q: chip8.Key4, sdl.K_w: chip8.Key5, sdl.K_e: chip8.Key6, sdl.K_r: chip8.KeyD, // Q W E R
		sdl.K_a: chip8.Key7, sdl.K_s: chip8.Key8, sdl.K_d: chip8.Key9, sdl.K_f: chip8.KeyE, // A S D F
		sdl.K_z: chip8.KeyA, sdl.K_x: chip8.Key0, sdl.K_c: chip8.KeyB, sdl.K_v: chip8.KeyF, // Z X C V
	}

	keyBindings = defaultKeyBindings
)

// State
var (
	c8 *chip8.Chip8

	window   *sdl.Window
	renderer *sdl.Renderer

	// Reuse pixel for drawing
	pixel = &sdl.Rect{
		W: int32(*scaleFactor),
		H: int32(*scaleFactor),
	}

	// Main loop control
	running bool
)

func main() {
	log.SetPrefix("[gochip8]")
	flag.Parse()

	romFile := flag.Arg(0)
	if len(romFile) == 0 {
		log.Fatalln("no rom file specified")
	}

	program := readProgram(romFile)

	if *disassemble {
		err := chip8.DisassembleProgram(program, ";", os.Stdout)
		if err != nil {
			log.Fatalln("error disassembling program:", err.Error())
		}
	} else {
		runROM(program)
	}
}

func readProgram(filename string) []byte {
	if program, err := ioutil.ReadFile(filename); err != nil {
		log.Fatalln("error reading from rom file:", err.Error())
		return nil
	} else {
		return program
	}
}

func setupSDL() {
	sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)

	w := *scaleFactor * chip8.DisplayWidth
	h := *scaleFactor * chip8.DisplayHeight

	romName := filepath.Base(flag.Arg(0))
	window = sdl.CreateWindow("gochip8 - "+romName, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, w, h, 0)
	renderer = sdl.CreateRenderer(window, -1, 0)
}

func cleanUpSDL() {
	renderer.Destroy()
	window.Destroy()
	sdl.Quit()
}

func runROM(rom []byte) {
	c8 = chip8.New(rom)

	// Lock goroutine to main thread
	runtime.LockOSThread()

	setupSDL()
	defer cleanUpSDL()

	// Update timers at 60Hz
	go func() {
		c := time.Tick(16 * time.Millisecond)
		for _ = range c {
			c8.UpdateTimers()
		}
	}()

	running = true
	c := time.Tick(time.Millisecond * time.Duration(1000/(*stepsPerSecond)))
	for _ = range c {
		if !running {
			break
		}

		processInput()

		// Run an instruction
		c8.Step()

		// Draw if needed
		if c8.DrawFlag {
			c8.DrawFlag = false
			draw()
		}
	}
}

// processInput handles polls and handles SDL events
func processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			running = false
		case *sdl.KeyDownEvent:
			if e.Keysym.Sym == sdl.K_ESCAPE {
				running = false
			} else {
				k, ok := keyBindings[e.Keysym.Sym]
				if ok {
					c8.PressKey(k)
				}
			}
		case *sdl.KeyUpEvent:
			k, ok := keyBindings[e.Keysym.Sym]
			if ok {
				c8.DePressKey(k)
			}
		}
	}
}

func draw() {
	display := c8.GetDisplay()
	renderer.SetDrawColor(0, 0, 0, 1)
	renderer.Clear()

	renderer.SetDrawColor(255, 255, 255, 1)
	for y := 0; y < chip8.DisplayHeight; y++ {
		for x := 0; x < chip8.DisplayWidth; x++ {
			// Only draw if set to 1
			if display[x][y] == 1 {
				pixel.X = int32(*scaleFactor * x)
				pixel.Y = int32(*scaleFactor * y)
				renderer.FillRect(pixel)
			}

		}
	}

	renderer.Present()
}
