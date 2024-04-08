package main

import "log"

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

	Key uint8 // Key that is pressed

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
	instruction := uint16((c.Memory[c.PC] << 8) | c.Memory[c.PC+1])
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
			c.ret()
		}
	}
}

// executes the Fetch, Decode and Execute cycles
func (c *Chip8) Cycle() {

}

// opcode functions
func (c *Chip8) cls() {

}

func (c *Chip8) ret() {

}
