package CPU

import (
	"fmt"

	"github.com/zoul0813/go6502/pkg/Memory"
)

type StatusFlag uint8

/*
	Status Register
	--------------------------------------------------
	7  bit  0
	---- ----
	NVss DIZC
	|||| ||||
	|||| |||+- Carry
	|||| ||+-- Zero
	|||| |+--- Interrupt Disable
	|||| +---- Decimal
	||++------ No CPU effect, see: the B flag
	|+-------- Overflow
	+--------- Negative
*/

// long and short form versions of the status register
const (
	Carry     = 0b00000001
	C         = 0b00000001
	Zero      = 0b00000010
	Z         = 0b00000010
	Interrupt = 0b00000100
	I         = 0b00000100
	Decimal   = 0b00001000
	D         = 0b00001000
	B         = 0b00010000
	Reserved  = 0b00100000
	Overflow  = 0b01000000
	V         = 0b01000000
	Negative  = 0b10000000
	N         = 0b10000000
)

type CPU struct {
	PC     uint16
	SP     uint8
	A      uint8
	X      uint8
	Y      uint8
	Status uint8

	SingleStep bool

	LastAddress uint16
}

// Utility Functions

const (
	Bit0 = 0b00000001
	Bit1 = 0b00000010
	Bit2 = 0b00000100
	Bit3 = 0b00001000
	Bit4 = 0b00010000
	Bit5 = 0b00100000
	Bit6 = 0b01000000
	Bit7 = 0b10000000
)

func BitSet(bit byte, value uint16) bool {
	b := value & uint16(bit)
	return b != 0
}

func IsNegative(value uint8) bool {
	n := value & 0b10000000
	return n != 0
}

func IsOverflow(prev uint8, current uint8) bool {
	p7 := BitSet(Bit7, uint16(prev))
	c7 := BitSet(Bit7, uint16(current))
	return p7 == c7
}

func New(
	PC uint16,
	SP uint8,
	A uint8,
	X uint8,
	Y uint8,
	Status uint8,
	SingleStep bool,
) *CPU {
	return &CPU{
		PC:         PC,
		SP:         SP,
		A:          A,
		X:          X,
		Y:          Y,
		Status:     Status,
		SingleStep: SingleStep,
	}
}

func (o *CPU) SetStatus(flag uint8, value bool) {
	var status byte
	if value {
		// cpu.Status |= 0b00000001 // bitset
		status = o.Status | flag
	} else {
		// cpu.Status &^ 0b01000000 // bitclear
		status = o.Status &^ flag
	}

	// fmt.Printf(" Status: %08b > %08b | %08b %v", o.Status, status, flag, value)

	o.Status = status
}

func (o *CPU) Debug() {
	// fmt.Printf("State: %v", o)
	fmt.Print("\n\nPC    SP    A    X    Y    Status   \n")
	fmt.Print("---------------------------NV-BDIZC-\n")
	//          PC    SP    A      X      Y      Status
	fmt.Printf("%04x  %04x  %02x   %02x   %02x   %08b\n\n",
		o.PC,
		o.SP,
		o.A,
		o.X,
		o.Y,
		o.Status,
	)
}

func (o *CPU) DebugBits() {
	fmt.Printf("      PC: %016b\n", o.PC)
	fmt.Printf("      SP: %08b\n", o.SP)
	fmt.Printf("       A: %08b\n", o.A)
	fmt.Printf("       X: %08b\n", o.X)
	fmt.Printf("       Y: %08b\n", o.Y)
	fmt.Printf(" Status : %08b\n", o.Status)
	fmt.Print("          NV-BDIZC\n\n")
}

func (o *CPU) Write(rom *Memory.Memory, b byte) error {
	err := rom.Set(o.LastAddress, b)
	return err
}

/*
  Addressing Modes
*/

func (o *CPU) Immediate(rom *Memory.Memory) (uint16, error) {
	o.LastAddress = o.PC
	o.PC += 1
	return o.LastAddress, nil
}

func (o *CPU) ZeroPage(rom *Memory.Memory) (uint16, error) {
	zp, err := rom.Get(o.PC)
	o.PC += 1
	o.LastAddress = uint16(zp)
	return o.LastAddress, err
}

func (o *CPU) ZeroPageX(rom *Memory.Memory) (uint16, error) {
	zp, err := rom.Get(o.PC)
	o.PC += 1
	o.LastAddress = uint16(zp + o.X)
	return o.LastAddress, err
}

func (o *CPU) ZeroPageY(rom *Memory.Memory) (uint16, error) {
	zp, err := rom.Get(o.PC)
	o.PC += 1
	o.LastAddress = uint16(zp + o.Y)
	return o.LastAddress, err
}

func (o *CPU) Absolute(rom *Memory.Memory) (uint16, error) {
	addr, err := rom.GetWord(o.PC)
	o.PC += 2
	o.LastAddress = addr
	return o.LastAddress, err
}

func (o *CPU) AbsoluteX(rom *Memory.Memory) (uint16, error) {
	addr, err := rom.GetWord(o.PC)
	o.PC += 2
	o.LastAddress = addr + uint16(o.X)
	return o.LastAddress, err
}

func (o *CPU) AbsoluteY(rom *Memory.Memory) (uint16, error) {
	addr, err := rom.GetWord(o.PC)
	o.PC += 2
	o.LastAddress = addr + uint16(o.Y)
	return o.LastAddress, err
}

func (o *CPU) Indirect(rom *Memory.Memory) (uint16, error) {
	from, _ := rom.GetWord(o.PC)
	o.PC += 2
	addr, err := rom.GetWord(from)
	o.LastAddress = addr
	return o.LastAddress, err
}

func (o *CPU) IndirectX(rom *Memory.Memory) (uint16, error) {
	zp, _ := rom.Get(o.PC)
	o.PC += 1
	addr, err := rom.GetWord(uint16(zp + o.X))
	o.LastAddress = addr
	return o.LastAddress, err
}

func (o *CPU) IndirectY(rom *Memory.Memory) (uint16, error) {
	zp, _ := rom.Get(o.PC)
	o.PC += 1
	addr1, err := rom.GetWord(uint16(zp))
	addr := addr1 + uint16(o.Y)
	o.LastAddress = addr
	return o.LastAddress, err
}
