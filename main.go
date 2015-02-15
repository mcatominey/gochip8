package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pmcatominey/gochip8/chip8"
	"github.com/veandco/go-sdl2/sdl"
)

// Config
var (
	// If this flag is present, the rom shold be loaded then output to stdout line by line with
	// an explanation of each opcode
	disassemble = flag.Bool("disassemble", false, "disassemble to stdout")

	// Scale factor for upscaling the display from Chip 8 resolution of 64*32
	scaleFactor = flag.Int("scaling", 10, "scale factor to multiply Chip 8 resolution (64*32) by")

	// Controls execution speed, useful since some roms play at mad speeds compared to others
	cyclesPerLoop = flag.Int("cycles", 10, "steps to emulate per loop")

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

	exitChan = make(chan bool, 1) // true sent this channel to exit main loop
)

func main() {
	flag.Parse()

	romFile := flag.Arg(0)
	if len(romFile) == 0 {
		fmt.Println("no rom file specified")
		os.Exit(1)
	}

	program := readProgram(romFile)
	if *disassemble {
		fmt.Println("Disassembling ROM to stdout")
		disassembleROM(program)
	} else {
		fmt.Println("Running ROM")
		runROM(program)
	}
}

func readProgram(filename string) []byte {
	if program, err := ioutil.ReadFile(filename); err != nil {
		fmt.Println("error reading from rom file:", err.Error())
		os.Exit(1)
		return nil
	} else {
		return program
	}
}

func setupSDL() {
	var (
		w = *scaleFactor * chip8.DisplayWidth
		h = *scaleFactor * chip8.DisplayHeight

		romName = filepath.Base(flag.Arg(0))
		err     error
	)

	sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	window, err = sdl.CreateWindow("gochip8 - "+romName, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, w, h, 0)
	if err != nil {
		panic(err)
	}

	renderer, err = sdl.CreateRenderer(window, -1, 0)
	if err != nil {
		panic(err)
	}
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

	// Setup and ensure cleanup of SDL
	setupSDL()
	defer cleanUpSDL()

	// Run main loop at 60Hz
	c := time.Tick(time.Second / 60)
	for _ = range c {
		select {
		case <-exitChan:
			return
		default:
		}

		c8.UpdateTimers()
		processInput()

		for i := 0; i < *cyclesPerLoop; i++ {
			c8.Step()
		}

		// Draw if needed
		if c8.DrawFlag {
			c8.DrawFlag = false
			draw()
		}
	}
}

func disassembleROM(rom []byte) {
	for i := 0; i < len(rom); i += 2 {
		op := chip8.GetOpcode(rom[i], rom[i+1])
		inst := chip8.DecodeOpcode(op)

		fmt.Printf("%#x ; %s\n", op, inst.Description)
	}
}

// processInput polls and handles SDL events
func processInput() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			exitChan <- true
		case *sdl.KeyDownEvent:
			if e.Keysym.Sym == sdl.K_ESCAPE {
				exitChan <- true
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
