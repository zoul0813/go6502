package CPU

import (
	"fmt"
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
}

func New(
	PC uint16,
	SP uint8,
	A uint8,
	X uint8,
	Y uint8,
	Status uint8,
) *CPU {
	return &CPU{
		PC:     PC,
		SP:     SP,
		A:      A,
		X:      X,
		Y:      Y,
		Status: Status,
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

func (o *CPU) IsNegative(value uint8) bool {
	n := value & 0b10000000
	return n != 0
}

func (o *CPU) Debug() {
	// fmt.Printf("State: %v", o)
	fmt.Print("\n\nPC    SP    A    X    Y    Status   \n")
	fmt.Print("---------------------------NV-BDIZC-\n")
	//          PC    SP    A      X      Y      Status
	fmt.Printf("%04x  %04x  %02x   %02x   %02x   %08b\n\n\n",
		o.PC,
		o.SP,
		o.A,
		o.X,
		o.Y,
		o.Status,
	)
}
