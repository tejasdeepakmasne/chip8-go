package main

import (
	"log"
	"math/rand"
	"time"
)

type Chip8 struct {
	Memory [4096]uint8 // memory 4KiB
	Stack  [16]uint16  //stack
	VX     [16]uint8   // 16 general purpose registers V0 - VF
	ST     uint8       // sound timer
	DT     uint8       // delay timer

	PC uint16 //program counter
	SP uint8  //stack pointer
	I  uint16 // used for pointing to memory addresses

	Display [64][42]uint8 // basically 0 for off and 1 for on for each pixel

	Key       uint8 // Key that is pressed
	FontStart uint
}

var font [90]uint8 = [90]uint8{
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

func (c *Chip8) loadFonts(start uint16) {
	for i := uint16(0); i < 90; i++ {
		c.Memory[start+i] = font[i]
	}
}

// stack functions
func (c *Chip8) stackPush(value uint16) {
	if c.SP == 15 {
		log.Fatalf("Stack Overflow")
	} else {
		c.Stack[c.SP] = value
		c.SP++
	}
}

func (c *Chip8) stackPop() uint16 {
	stackTop := c.Stack[c.SP]
	if c.SP != 0 {
		c.SP--
	}
	return stackTop
}

func (c *Chip8) Fetch() uint16 {
	instruction := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
	c.PC += 2
	return instruction
}

func (c *Chip8) Decode(instruction uint16) {
	// temp values which can be used depending on instruction
	vx := instruction & 0x0F00 >> 8
	vy := instruction & 0x00F0 >> 4
	n := instruction & 0x000F
	nn := instruction & 0x00FF
	nnn := instruction & 0x0FFF

	code := instruction & 0xF000
	switch code {
	case 0x0000:
		switch nn {
		case 0xE0:
			c.cls()
		case 0xEE:
			c.PC = c.stackPop()
		}

	case 0x1000:
		c.PC = nnn

	case 0x2000:
		c.stackPush(c.PC)
		c.PC = nnn

	case 0x3000:
		if c.VX[vx] == uint8(nn) {
			c.PC += 2
		}

	case 0x4000:
		if c.VX[vx] != uint8(nn) {
			c.PC += 2
		}

	case 0x5000:
		if c.VX[vx] == c.VX[vy] {
			c.PC += 2
		}

	case 0x6000:
		c.VX[vx] = uint8(nn)

	case 0x7000:
		c.VX[vx] += uint8(nn)

	case 0x8000:
		switch n {
		case 0x0:
			c.VX[vx] = c.VX[vy]
		case 0x1:
			c.VX[vx] = c.VX[vx] | c.VX[vy]
		case 0x2:
			c.VX[vx] = c.VX[vx] & c.VX[vy]
		case 0x3:
			c.VX[vx] = c.VX[vx] ^ c.VX[vy]
		case 0x4:
			temp := uint(c.VX[vx]) + uint(c.VX[vy])
			if temp > 255 {
				c.VX[0xF] = 1
			} else {
				c.VX[0xF] = 0
			}
			c.VX[vx] += c.VX[vy]
		case 0x5:
			if c.VX[vx] > c.VX[vy] {
				c.VX[0xF] = 1
			} else {
				c.VX[0xF] = 0
			}
			c.VX[vx] -= c.VX[vy]
		case 0x6:
			if c.VX[vx]&0x1 == 0x1 {
				c.VX[0xF] = 1
			} else {
				c.VX[vx] = 0
			}
			c.VX[vx] = c.VX[vx] / 2
		case 0x7:
			if c.VX[vy] > c.VX[vx] {
				c.VX[0xF] = 1
			} else {
				c.VX[0xF] = 0
			}
			c.VX[vx] = c.VX[vy] - c.VX[vx]
		case 0xE:
			if c.VX[vx]&0b10000000 == 1 {
				c.VX[0xF] = 1
			} else {
				c.VX[0xF] = 0
			}
			c.VX[vx] = c.VX[vx] * 2
		}
	case 0x9000:
		if c.VX[vx] != c.VX[vy] {
			c.PC += 2
		}
	case 0xA000:
		c.I = nnn
	case 0xB000:
		c.PC = nnn + uint16(c.VX[0x0])
	case 0xC000:
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rngn := r.Uint32() % 255
		c.VX[vx] = uint8(rngn) & uint8(nn)
	case 0xD000:
		//implement draw
	case 0xE000:
		//implement keyboard instructions
	case 0xF000:
		switch nn {
		case 0x07:
			c.VX[vx] = c.DT
		case 0x0A:
			//store key press
		case 0x15:
			c.DT = c.VX[vx]
		case 0x18:
			c.ST = c.VX[vx]
		case 0x1E:
			c.I += uint16(c.VX[vx])
		case 0x29:
			c.I = uint16(c.FontStart) + (5 * vx)
		case 0x33:
			//BCD vallues??
		case 0x55:
			for i := uint16(0); i <= vx; i++ {
				c.Memory[c.I+i] = c.VX[vx+i]
			}
		case 0x65:
			for i := uint16(0); i <= vx; i++ {
				c.VX[vx+i] = c.Memory[c.I+i]
			}
		}
	}
}

// executes the Fetch, Decode and Execute cycles
func (c *Chip8) Cycle() {
	c.Fetch()
}

// opcode functions
func (c *Chip8) cls() {

}
