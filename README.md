# gochip8

An emulator of the Chip 8 interpreted language in Go.

Under MIT license.

[![Build Status](https://travis-ci.org/pmcatominey/gochip8.svg?branch=master)](https://travis-ci.org/pmcatominey/gochip8)

## Usage

```gochip8 <flags> path/to/rom```

**Flags**

- ```-disassemble``` instead of running the rom, print an explanation of each opcode to stdout
- ```-scaling 10``` factor to scale from original Chip 8 resolution (64x32), defaults to 10 for a window size of 640x320
- ```-cycles 10``` number of steps to attempt to emulate per loop

A collection of games, understood to be in the public domain are in the ```games``` directory.

### Controls

The Chip 8 has a hexidecimal keyboard which is bound to these keys:

**1 2 3 4**

**Q W E R**

**A S D F**

**Z X C V**

## Building

**Go installation and C compiler required**

### 1. Get SDL2

Requries SDL2 library, usually available in your systems package manager,
	
**Ubuntu**
```apt-get install libsdl2```

**Mac OSX**
```brew install sdl2```

**Windows**

Install the MinGW development build of SDL2, ensure you have a 64bit MinGW toolchain.

### 2. go-sdl2

The only Go dependency is the SDL2 bindings:

```go get -v github.com/veandco/go-sdl2/sdl```

### 3. Build & Run

```go build``` will build an executable, to run on Windows you must have the SDL2.dll in the same path as the executable or in the system path.

## Testing

The aim for testing is to have at least one unit test per instruction to verify correctness
of the program. Some instructions still require a unit test.

Run ```go test -cover ./chip8```.

## Reference

Built using Reference.html found [here](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM),
a copy is also included.

Wikipedia also has a [good article](http://en.wikipedia.org/wiki/CHIP-8) on Chip 8.
