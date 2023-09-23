package CPU

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/zoul0813/go6502/pkg/IO"
	"golang.org/x/image/font"
)

type StatusFlag uint8

// type word = uint16

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
	PC         uint16
	SP         uint8
	A          uint8
	X          uint8
	Y          uint8
	Status     uint8
	SingleStep bool
	Address    uint16
	DebugMode  bool
	halted     bool
}

const ZP_HEAD = 0x000
const STACK_HEAD = 0x100

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

func BitTest(bit uint8, value uint8) bool {
	b := value & bit
	return b != 0
}

func IsNegative(value uint8) bool {
	n := value & 0b10000000
	return n != 0
}

func IsOverflow(prev uint8, current uint8) bool {
	v := ((prev ^ current) & 0x80)
	return v != 0
}

func New(
	PC uint16,
	SP uint8,
	A uint8,
	X uint8,
	Y uint8,
	Status uint8,
	SingleStep bool,
	DebugMode bool,
) *CPU {
	return &CPU{
		PC:         PC,
		SP:         SP,
		A:          A,
		X:          X,
		Y:          Y,
		Status:     Status,
		SingleStep: SingleStep,
		DebugMode:  DebugMode,
	}
}

func (o *CPU) Step(io IO.Memory) (bool, error) {
	halted := false
	b, _ := io.Get(o.PC)
	var instr OpCode = OpCode(b)
	o.PC++ // increment the stack pointer
	o.Log("Instruction: %02x\n", instr)
	switch instr {
	// Jump/Branch Instructions
	case JMP_A:
		// jump absolute
		o.Log("I: JMP ")
		addr, _ := o.Absolute(io)
		o.Log("%04x (ABS)", addr)
		o.PC = addr
	case JMP_IN:
		// jump indirect
		o.Log("I: JMP ")
		// from, _ := io.GetWord(o.PC)
		// addr, _ := io.GetWord(fio)
		addr, _ := o.Indirect(io)
		o.Log("%04x (Indirect)", addr)
		o.PC = addr
	case JSR_A:
		// jump to subroute, absolute
		o.Log("I: JSR ")
		// addr, _ := io.GetWord(o.PC) // jump to here
		// o.PC += 2
		addr, _ := o.Absolute(io)
		o.Log("%04x (ABS)", addr)

		pc := o.PC - 1
		var lo uint8 = uint8(pc & 0b0000000011111111)
		var hi uint8 = uint8(pc >> 8)

		o.SP--
		io.Set(STACK_HEAD+uint16(o.SP), hi)
		o.SP--
		io.Set(STACK_HEAD+uint16(o.SP), lo)

		o.PC = addr
	case RTS:
		o.Log("I: RTS ")
		o.Log("(Implied)")
		lo, _ := io.Get(STACK_HEAD + uint16(o.SP))
		o.SP++
		hi, _ := io.Get(STACK_HEAD + uint16(o.SP))
		o.SP++
		var addr uint16 = (uint16(hi) << 8) | uint16(lo)
		o.PC = addr + 1
	case RTI:
		o.Log("I: RTI ")
		o.Log("(Implied)")
		o.SP++
		status, _ := io.Get(STACK_HEAD + uint16(o.SP))
		o.Status = status
		o.SP++
		lo, _ := io.Get(STACK_HEAD + uint16(o.SP))
		o.SP++
		hi, _ := io.Get(STACK_HEAD + uint16(o.SP))
		var addr uint16 = (uint16(hi) << 8) | uint16(lo)
		o.Log(" RTI: %02x %02x %02x %04x \n", status, lo, hi, addr)
		o.PC = addr
	case BPL:
		// branch on plus
		o.Log("I: BPL ")
		rel, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (Rel)", rel)
		neg := BitTest(Negative, o.Status)
		if !neg {
			o.PC += uint16(rel)
		}
	case BMI:
		// branch on minus
		o.Log("I: BMI ")
		rel, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (Rel)", rel)
		neg := BitTest(Negative, o.Status)
		if neg {
			o.PC += uint16(rel)
		}
	case BVC:
		// branch on overflow clear
		o.Log("I: BVC ")
		rel, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (Rel)", rel)
		overflow := BitTest(Overflow, o.Status)
		if !overflow {
			o.PC += uint16(rel)
		}
	case BVS:
		// branch on overflow set
		o.Log("I: BVS ")
		rel, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (Rel)", rel)
		overflow := BitTest(Overflow, o.Status)
		if overflow {
			o.PC += uint16(rel)
		}
	case BCC:
		// branch on carry clear
		o.Log("I: BCC ")
		rel, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (Rel)", rel)
		carry := BitTest(Carry, o.Status)
		if !carry {
			o.PC += uint16(rel)
		}
	case BCS:
		// branch on carry set
		o.Log("I: BCS ")
		rel, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (Rel)", rel)
		carry := BitTest(Carry, o.Status)
		if carry {
			o.PC += uint16(rel)
		}
	case BNE:
		// branch on not equal
		o.Log("I: BNE ")
		rel, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (Rel)", rel)
		zero := BitTest(Zero, o.Status)
		if !zero {
			o.PC += uint16(rel)
		}
	case BEQ:
		// branch on equal
		o.Log("I: BEQ ")
		rel, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (Rel)", rel)
		zero := BitTest(Zero, o.Status)
		if zero {
			o.PC += uint16(rel)
		}

	// Misc
	case BRK: // TODO: NMI
		// break
		o.Log("I: BRK ")
		// TODO: non-maskable interrupt
		o.PC += 2
	case NOP:
		o.Log("I: NOP")

	// Add (ADC)
	case ADC_I:
		// add with carry, immediate
		o.Log("I: ADC ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		o.ADC(b)
		// N V Z C
	case ADC_ZP:
		// add with carry, zero page
		o.Log("I: ADC ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		o.ADC(b)
		// N V Z C
	case ADC_ZPX:
		// add with carry, zero page, x
		o.Log("I: ADC ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		o.ADC(b)
		// N V Z C
	case ADC_A:
		// add with carry, absolute
		o.Log("I: ADC ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		o.ADC(b)
		// N V Z C
	case ADC_AX:
		// add with carry, absolute, x
		o.Log("I: ADC ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		o.ADC(b)
		// N V Z C
	case ADC_AY:
		// add with carry, absolute, y
		o.Log("I: ADC ")
		addr, _ := o.AbsoluteY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, Y)", addr)

		o.ADC(b)
		// N V Z C
	case ADC_INX:
		// add with carry, indirect, x
		o.Log("I: ADC ")
		addr, _ := o.IndirectX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, X)", addr)

		o.ADC(b)
		// N V Z C
	case ADC_INY:
		// add with carry, indirect, y
		o.Log("I: ADC ")
		addr, _ := o.IndirectY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, Y)", addr)

		o.ADC(b)
		// N V Z C

	// Subtract (SBC)
	case SBC_I:
		// subtract with carry, immediate
		o.Log("I: SBC ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		o.SBC(b)
		// N V Z C
	case SBC_ZP:
		// subtract with carry, zero page
		o.Log("I: SBC ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		o.SBC(b)
		// N V Z C
	case SBC_ZPX:
		// subtract with carry, zero page, x
		o.Log("I: SBC ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		o.SBC(b)
		// N V Z C
	case SBC_A:
		// subtract with carry, absolute
		o.Log("I: SBC ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		o.SBC(b)
		// N V Z C
	case SBC_AX:
		// subtract with carry, absolute, x
		o.Log("I: SBC ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		o.SBC(b)
		// N V Z C
	case SBC_AY:
		// subtract with carry, absolute, y
		o.Log("I: SBC ")
		addr, _ := o.AbsoluteY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, Y)", addr)

		o.SBC(b)
		// N V Z C
	case SBC_INX:
		// subtract with carry, indirect, x
		o.Log("I: SBC ")
		addr, _ := o.IndirectX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, X)", addr)

		o.SBC(b)
		// N V Z C
	case SBC_INY:
		// subtract with carry, indirect, y
		o.Log("I: SBC ")
		addr, _ := o.IndirectY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, Y)", addr)

		o.SBC(b)
		// N V Z C

	// Compare (A, X, Y)
	case CMP_I:
		// compare accumulator
		o.Log("I: CMP ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		a := o.A
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.A == b)
		o.SetStatus(Carry, o.A >= b)
	case CMP_ZP:
		// compare accumulator, zero page
		o.Log("I: CMP ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		a := o.A
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.A == b)
		o.SetStatus(Carry, o.A >= b)
	case CMP_ZPX:
		// compare accumulator, zero page, x
		o.Log("I: CMP ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		a := o.A
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.A == b)
		o.SetStatus(Carry, o.A >= b)
	case CMP_A:
		// compare accumulator, absolute
		o.Log("I: CMP ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		a := o.A
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.A == b)
		o.SetStatus(Carry, o.A >= b)
	case CMP_AX:
		// compare accumulator, absolute, x
		o.Log("I: CMP ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		a := o.A
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.A == b)
		o.SetStatus(Carry, o.A >= b)
	case CMP_AY:
		// compare accumulator, absolute, y
		o.Log("I: CMP ")
		addr, _ := o.AbsoluteY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, Y)", addr)

		a := o.A
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.A == b)
		o.SetStatus(Carry, o.A >= b)
	case CMP_INX:
		// compare accumulator, indirect, x
		o.Log("I: CMP ")
		addr, _ := o.IndirectX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, X)", addr)

		a := o.A
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.A == b)
		o.SetStatus(Carry, o.A >= b)
	case CMP_INY:
		// compare accumulator, indirect y
		o.Log("I: CMP ")
		addr, _ := o.IndirectY(io)
		b, _ := io.Get(addr)

		o.Log("%04x (Indirect, Y)", addr)

		a := o.A
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.A == b)
		o.SetStatus(Carry, o.A >= b)
	case CPX:
		// compare x
		o.Log("I: CPX ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		a := o.X
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.X == b)
		o.SetStatus(Carry, o.X >= b)
	case CPX_ZP:
		// compare x, zero page
		o.Log("I: CPX ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		a := o.X
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.X == b)
		o.SetStatus(Carry, o.X >= b)
	case CPX_A:
		// compare x, absolute
		o.Log("I: CPX ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		a := o.X
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.X == b)
		o.SetStatus(Carry, o.X >= b)
	case CPY:
		// compare y
		o.Log("I: CPY ")
		v, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (Immediate)", v)
		a := o.Y
		r := a - v // actually do the math, so we can determine if it's negative
		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.Y == v)
		o.SetStatus(Carry, o.Y >= v)
	case CPY_ZP:
		// compare x, zero page
		o.Log("I: CPY ")
		zp, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (ZP)", zp)
		v, _ := io.Get(uint16(zp))
		a := o.Y
		r := a - v // actually do the math, so we can determine if it's negative
		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.Y == v)
		o.SetStatus(Carry, o.Y >= v)
	case CPY_A:
		// compare x, absolute
		o.Log("I: CPY ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		a := o.Y
		r := a - b // actually do the math, so we can determine if it's negative

		o.SetStatus(Negative, IsNegative(r))
		o.SetStatus(Zero, o.Y == b)
		o.SetStatus(Carry, o.Y >= b)

	// AND
	case AND_I:
		// and with a
		o.Log("I: AND ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		o.A &= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case AND_ZP:
		// and with a, zero page
		o.Log("I: AND ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		o.A &= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case AND_ZPX:
		// and with a, zero page, x
		o.Log("I: AND ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		o.A &= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case AND_A:
		// and with a, absolute
		o.Log("I: AND ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		o.A &= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case AND_AX:
		// and with a, absolute, x
		o.Log("I: AND ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		o.A &= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case AND_AY:
		// and with a, absolute, y
		o.Log("I: AND ")
		addr, _ := o.AbsoluteY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, Y)", addr)

		o.A &= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case AND_INX:
		// and with a, indirect, x
		o.Log("I: AND ")
		addr, _ := o.IndirectX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, X)", addr)

		o.A &= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case AND_INY:
		// and with a, indirect, y
		o.Log("I: AND ")
		addr, _ := o.IndirectY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, Y)", addr)

		o.A &= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)

	// EOR
	case EOR_I:
		// exclusive or, immediate
		o.Log("I: EOR ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		o.A ^= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case EOR_ZP:
		// exclusive or, zero page
		o.Log("I: EOR ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		o.A ^= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case EOR_ZPX:
		// exclusive or, zeor page, x
		o.Log("I: EOR ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		o.A ^= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case EOR_A:
		// exclusive or, absolute
		o.Log("I: EOR ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		o.A ^= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case EOR_AX:
		// exclusive or, absolute, x
		o.Log("I: EOR ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		o.A ^= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case EOR_AY:
		// exclusive or, absolute, y
		o.Log("I: EOR ")
		addr, _ := o.AbsoluteY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, Y)", addr)

		o.A ^= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case EOR_INX:
		// exclusive or, indirect, x
		o.Log("I: EOR ")
		addr, _ := o.IndirectX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, X)", addr)
		o.A ^= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case EOR_INY:
		// exclusive or, indirect, y
		o.Log("I: EOR ")
		addr, _ := o.IndirectY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, Y)", addr)
		o.A ^= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)

		// ORA
	case ORA_I:
		// exclusive or, immediate
		o.Log("I: ORA ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		o.A |= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case ORA_ZP:
		// exclusive or, zero page
		o.Log("I: ORA ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		o.A |= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case ORA_ZPX:
		// exclusive or, zeor page, x
		o.Log("I: ORA ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		o.A |= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case ORA_A:
		// exclusive or, absolute
		o.Log("I: ORA ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		o.A |= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case ORA_AX:
		// exclusive or, absolute, x
		o.Log("I: ORA ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		o.A |= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case ORA_AY:
		// exclusive or, absolute, y
		o.Log("I: ORA ")
		addr, _ := o.AbsoluteY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, Y)", addr)

		o.A |= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case ORA_INX:
		// exclusive or, indirect, x
		o.Log("I: ORA ")
		addr, _ := o.IndirectX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, X)", addr)
		o.A |= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
	case ORA_INY:
		// exclusive or, indirect, y
		o.Log("I: ORA ")
		addr, _ := o.IndirectY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, Y)", addr)

		o.A |= b

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)

	// Store Instructions (STA, STX, STY)
	case STA_ZP:
		// store a, zero page
		o.Log("I: STA ")
		addr, _ := o.ZeroPage(io)
		o.Log("%02x (ZP)", addr)

		io.Set(addr, o.A)
	case STA_ZPX:
		// store a, zero page, x
		o.Log("I: STA ")
		addr, _ := o.ZeroPageX(io)
		o.Log("%02x (ZP, X)", addr)

		io.Set(addr, o.A)
	case STA_A:
		// store a, zero page, x
		o.Log("I: STA ")
		addr, _ := o.Absolute(io)
		o.Log("%04x (ABS)", addr)

		io.Set(addr, o.A)
	case STA_AX:
		// store a, zero page, x
		o.Log("I: STA ")
		addr, _ := o.AbsoluteX(io)
		o.Log("%04x (ABS, X)", addr)

		io.Set(addr, o.A)
	case STA_AY:
		// store a, zero page, x
		o.Log("I: STA ")
		addr, _ := o.AbsoluteY(io)
		o.Log("%04x (ABS, Y)", addr)

		io.Set(addr, o.A)
	case STA_INX: // store a, indirect, x
		o.Log("I: STA ")
		addr, _ := o.IndirectX(io)
		o.Log("%04x (Indirect, X)", addr)

		io.Set(addr, o.A)
	case STA_INY: // store a, indirect, y
		o.Log("I: STA ")
		addr, _ := o.IndirectY(io)
		o.Log("%04x (Indirect, Y)", addr)

		io.Set(addr, o.A)
	case STX_ZP:
		// store a, zero page
		o.Log("I: STX ")
		addr, _ := o.ZeroPage(io)
		o.Log("%02x (ZP)", addr)

		io.Set(addr, o.X)
	case STX_ZPY:
		// store a, zero page, y
		o.Log("I: STY ")
		addr, _ := o.ZeroPageY(io)
		o.Log("%02x (ZP, Y)", addr)

		io.Set(addr, o.X)
	case STX_A:
		// store a, zero page, x
		o.Log("I: STX ")
		addr, _ := o.Absolute(io)
		o.Log("%04x (ABS)", addr)

		io.Set(addr, o.X)
	case STY_ZP:
		// store a, zero page
		o.Log("I: STY ")
		addr, _ := o.ZeroPage(io)
		o.Log("%02x (ZP)", addr)

		io.Set(addr, o.Y)
	case STY_ZPX:
		// store a, zero page, x
		o.Log("I: STY ")
		addr, _ := o.ZeroPageX(io)
		o.Log("%02x (ZP, X)", addr)

		io.Set(addr, o.Y)
	case STY_A:
		// store a, zero page, x
		o.Log("I: STY ")
		addr, _ := o.Absolute(io)
		o.Log("%04x (ABS)", addr)

		io.Set(addr, o.Y)

	// INC/DEC Instructions
	case INC_ZP:
		// increment zero page
		o.Log("I: INC ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		b++
		io.Set(addr, b)

		o.SetStatus(Zero, b == 0)
		o.SetStatus(Negative, IsNegative(b))
	case INC_ZPX:
		// increment zero page, x
		o.Log("I: INC ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		b++
		io.Set(addr, b)

		o.SetStatus(Zero, b == 0)
		o.SetStatus(Negative, IsNegative(b))
	case INC_A:
		// increment absolute
		o.Log("I: INC ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		b++
		io.Set(addr, b)

		o.SetStatus(Zero, b == 0)
		o.SetStatus(Negative, IsNegative(b))
	case INC_AX:
		// increment absolute, x
		o.Log("I: INC ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		b++
		io.Set(addr, b)

		o.SetStatus(Zero, b == 0)
		o.SetStatus(Negative, IsNegative(b))
	case DEC_ZP:
		// increment zero page
		o.Log("I: DEC ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		b--
		io.Set(addr, b)

		o.SetStatus(Zero, b == 0)
		o.SetStatus(Negative, IsNegative(b))
	case DEC_ZPX:
		// increment zero page, x
		o.Log("I: DEC ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		b--
		io.Set(addr, b)

		o.SetStatus(Zero, b == 0)
		o.SetStatus(Negative, IsNegative(b))
	case DEC_A:
		// increment absolute
		o.Log("I: DEC ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		b--
		io.Set(addr, b)

		o.SetStatus(Zero, b == 0)
		o.SetStatus(Negative, IsNegative(b))
	case DEC_AX:
		// increment absolute, x
		o.Log("I: DEC ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		b--
		io.Set(addr, b)

		o.SetStatus(Zero, b == 0)
		o.SetStatus(Negative, IsNegative(b))
	case INX:
		// increment x
		o.Log("I: INX ")
		o.X++
		o.SetStatus(Zero, o.X == 0)
		// o.SetStatus(Overflow, o.X < x)
		o.SetStatus(Negative, IsNegative(o.X))
	case DEX:
		// increment x
		o.Log("I: DEX ")
		o.X--
		o.SetStatus(Zero, o.X == 0)
		o.SetStatus(Negative, IsNegative(o.X))
	case INY:
		// increment y
		o.Log("I: INY ")
		o.Y++
		o.SetStatus(Zero, o.Y == 0)
		o.SetStatus(Negative, IsNegative(o.Y))
	case DEY:
		// increment y
		o.Log("I: INY ")
		o.Y--
		o.SetStatus(Zero, o.Y == 0)
		o.SetStatus(Negative, IsNegative(o.Y))

	// Load Register Instructions
	case LDA_I:
		// load A immediate
		o.Log("I: LDA ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		o.A = b

		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))
	case LDA_ZP:
		// load A zero page
		o.Log("I: LDA ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		o.A = b

		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))
	case LDA_ZPX:
		// load A zero page, x index
		o.Log("I: LDA ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		o.A = b

		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))
	case LDA_A:
		// load A absolute
		o.Log("I: LDA ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		o.A = b

		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))
	case LDA_AX:
		// load A absolute, x
		o.Log("I: LDA ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		o.A = b

		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))
	case LDA_AY:
		// load A absolute, y
		o.Log("I: LDA ")
		addr, _ := o.AbsoluteY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, Y)", addr)

		o.A = b

		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))
	case LDA_INX:
		// load A x index, indirect
		o.Log("I: LDA ")
		addr, _ := o.IndirectX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, X)", addr)

		o.A = b

		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))
	case LDA_INY:
		// load A indirect, y index
		o.Log("I: LDA ")
		addr, _ := o.IndirectY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (Indirect, Y)", addr)

		o.A = b

		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))
	case LDX_I:
		// load X immediate
		o.Log("I: LDX ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		o.X = b

		o.SetStatus(Zero, o.X == 0)
		o.SetStatus(Negative, IsNegative(o.X))
	case LDX_ZP:
		// load X zero page
		o.Log("I: LDX ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		o.X = b

		o.SetStatus(Zero, o.X == 0)
		o.SetStatus(Negative, IsNegative(o.X))
	case LDX_ZPY:
		// load X zero page, y index
		o.Log("I: LDX ")
		addr, _ := o.ZeroPageY(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, Y)", addr)

		o.X = b

		o.SetStatus(Zero, o.X == 0)
		o.SetStatus(Negative, IsNegative(o.X))
	case LDX_A:
		// load X absolute
		o.Log("I: LDX ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		o.X = b

		o.SetStatus(Zero, o.X == 0)
		o.SetStatus(Negative, IsNegative(o.X))
	case LDX_AY:
		// load X absolute, y index
		o.Log("I: LDX ")
		addr, _ := o.AbsoluteY(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, Y)", addr)

		o.X = b

		o.SetStatus(Zero, o.X == 0)
		o.SetStatus(Negative, IsNegative(o.X))
	case LDY_I:
		// load Y immediate
		o.Log("I: LDY ")
		addr, _ := o.Immediate(io)
		b, _ := io.Get(addr)
		o.Log("%02x (Immediate)", b)

		o.Y = b

		o.SetStatus(Zero, o.Y == 0)
		o.SetStatus(Negative, IsNegative(o.Y))
	case LDY_ZP:
		// load Y zero page
		o.Log("I: LDY ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		o.Y = b

		o.SetStatus(Zero, o.Y == 0)
		o.SetStatus(Negative, IsNegative(o.Y))
	case LDY_ZPX:
		// load Y zero page, x index
		o.Log("I: LDY ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, Y)", addr)

		o.Y = b

		o.SetStatus(Zero, o.Y == 0)
		o.SetStatus(Negative, IsNegative(o.Y))
	case LDY_A:
		// load Y absolute
		o.Log("I: LDY ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		o.Y = b

		o.SetStatus(Zero, o.Y == 0)
		o.SetStatus(Negative, IsNegative(o.Y))
	case LDY_AX:
		// load Y absolute, x index
		o.Log("I: LDY ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		o.Y = b

		o.SetStatus(Zero, o.Y == 0)
		o.SetStatus(Negative, IsNegative(o.Y))

	// Register Transfer Instructions
	case TAX:
		// transfer a to x
		o.Log("I: TAX ")
		o.X = o.A
		o.SetStatus(Zero, o.X == 0)
		o.SetStatus(Negative, IsNegative(o.X))
	case TXA:
		// transfer x to a
		o.Log("I: TXA ")
		o.A = o.X
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))
	case TAY:
		// transfer a to y
		o.Log("I: TAY ")
		o.Y = o.A
		o.SetStatus(Zero, o.Y == 0)
		o.SetStatus(Negative, IsNegative(o.Y))
	case TYA:
		// transfer y to a
		o.Log("I: TYA ")
		o.A = o.Y
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Negative, IsNegative(o.A))

	// Status Instructions
	case CLC:
		// clear carry
		o.Log("I: CLC ")
		o.SetStatus(Carry, false)
	case SEC:
		// set carry
		o.Log("I: SEC ")
		o.SetStatus(Carry, true)
	case CLI:
		// clear interrupt
		o.Log("I: CLI ")
		o.SetStatus(Interrupt, false)
	case SEI:
		// set interrupt
		o.Log("I: SEI ")
		o.SetStatus(Interrupt, true)
	case CLV:
		// clear overflow
		o.Log("I CLV ")
		o.SetStatus(Overflow, false)
	case CLD:
		// clear decimal
		o.Log("I CLD ")
		o.SetStatus(Decimal, false)
	case SED:
		// set decimal
		o.Log("I SED ")
		o.SetStatus(Decimal, true)
		o.Log("Not yet implemented ... ")
		o.SetStatus(Decimal, true)

	// Bit Shift Instructions
	case ROL:
		// rotate left, a
		o.Log("I: ROL ")
		o.Log("%02x (A)", o.A)
		carry := BitTest(Carry, o.Status)
		b7 := BitTest(o.A, Bit7)
		a := (o.A << 1)
		if carry {
			a++
		}
		o.A = a

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b7)
	case ROL_ZP:
		// rotate left, zerp page
		o.Log("I: ROL ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		carry := BitTest(Carry, o.Status)
		b7 := BitTest(b, Bit7)
		v := (b << 1)
		if carry {
			v++
		}
		io.Set(uint16(b), v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b7)
	case ROL_ZPX:
		// rotate left, zerp page, x
		o.Log("I: ROL ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		carry := BitTest(Carry, o.Status)
		b7 := BitTest(b, Bit7)
		v := (b << 1)
		if carry {
			v++
		}
		io.Set(addr, v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b7)
	case ROL_A:
		// rotate left, absolute
		o.Log("I: ROL ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		carry := BitTest(Carry, o.Status)
		b7 := BitTest(b, Bit7)
		v := (b << 1)
		if carry {
			v++
		}
		io.Set(addr, v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b7)
	case ROL_AX:
		// rotate left, absolute, x
		o.Log("I: ROL ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		carry := BitTest(Carry, o.Status)
		b7 := BitTest(b, Bit7)
		v := (b << 1)
		if carry {
			v++
		}
		io.Set(addr, v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b7)
	case ROR:
		// rotate left, a
		o.Log("I: ROR ")
		o.Log("%02x (A)", o.A)
		carry := BitTest(Carry, o.Status)
		b0 := BitTest(o.A, Bit0)
		v := (o.A >> 1)
		if carry {
			v |= 0b10000000 // bitset
		}
		o.A = v
		o.Log(" %08b", v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b0)
	case ROR_ZP:
		// rotate left, zerp page
		o.Log("I: ROR ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		carry := BitTest(Carry, o.Status)
		b0 := BitTest(b, Bit0)
		v := (b >> 1)
		if carry {
			v |= 0b10000000 // bitset
		}
		io.Set(uint16(b), v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b0)
	case ROR_ZPX:
		// rotate left, zerp page, x
		o.Log("I: ROR ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		carry := BitTest(Carry, o.Status)
		b0 := BitTest(b, Bit0)
		v := (b >> 1)
		if carry {
			v |= 0b10000000 // bitset
		}
		io.Set(addr, v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b0)
	case ROR_A:
		// rotate left, absolute
		o.Log("I: ROR ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)
		carry := BitTest(Carry, o.Status)
		b0 := BitTest(b, Bit0)
		v := b >> 1
		if carry {
			v |= 0b10000000 // bitset
		}
		io.Set(addr, v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b0)
	case ROR_AX:
		// rotate left, absolute, x
		o.Log("I: ROR ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		carry := BitTest(Carry, o.Status)
		b0 := BitTest(b, Bit0)
		v := b >> 1
		if carry {
			v |= 0b10000000 // bitset
		}
		io.Set(addr, v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b0)
	case ASL:
		// arithmetic shift left, a
		o.Log("I: ASL ")
		o.Log("%02x", o.A)

		b7 := BitTest(o.A, Bit7)
		a := o.A << 1
		o.A = a

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b7)
	case ASL_ZP:
		// arithmetic shift left, zero page
		o.Log("I: ASL ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		b7 := BitTest(b, Bit7)
		b = b << 1
		io.Set(addr, b)

		o.SetStatus(Negative, IsNegative(b))
		o.SetStatus(Zero, b == 0)
		o.SetStatus(Carry, b7)
	case ASL_ZPX:
		// arithmetic shift left, zero page, x
		o.Log("I: ASL ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		b7 := BitTest(b, Bit7)
		v := b << 1
		io.Set(addr, v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b7)
	case ASL_A:
		// arithmetic shift left, absolute
		o.Log("I: ROL ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		b7 := BitTest(b, Bit7)
		v := b << 1
		io.Set(addr, v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b7)
	case ASL_AX:
		// arithmetic shift left, absolute, x
		o.Log("I: ASL ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		b7 := BitTest(b, Bit7)
		v := b << 1
		io.Set(addr, v)

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, b7)
	case LSR:
		// logical shift right, a
		o.Log("I: LSR ")
		carry := BitTest(Bit0, o.A)
		o.A = o.A >> 1

		o.SetStatus(Negative, IsNegative(o.A))
		o.SetStatus(Zero, o.A == 0)
		o.SetStatus(Carry, carry)
	case LSR_ZP:
		// logical shift right, zero page
		o.Log("I: LSR ")
		// zp, _ := io.Get(o.PC)
		// o.PC++
		// b, _ := io.Get(uint16(zp))
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		carry := BitTest(Bit1, b)
		b = b >> 1
		io.Set(addr, b)

		o.SetStatus(Negative, IsNegative(b))
		o.SetStatus(Zero, b == 0)
		o.SetStatus(Carry, carry)
	case LSR_ZPX:
		// logical shift right, zero page
		o.Log("I: LSR ")
		addr, _ := o.ZeroPageX(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP, X)", addr)

		carry := BitTest(Bit1, b)
		b = b >> 1
		io.Set(addr, b)

		o.SetStatus(Negative, IsNegative(b))
		o.SetStatus(Zero, b == 0)
		o.SetStatus(Carry, carry)
	case LSR_A:
		// logical shift right, absolute
		o.Log("I: LSR ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		carry := BitTest(Bit1, b)
		b = b >> 1
		io.Set(addr, b)

		o.SetStatus(Negative, IsNegative(b))
		o.SetStatus(Zero, b == 0)
		o.SetStatus(Carry, carry)
	case LSR_AX:
		// logical shift right, absolute, x
		o.Log("I: LSR ")
		addr, _ := o.AbsoluteX(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS, X)", addr)

		carry := BitTest(Bit1, b)
		b = b >> 1
		io.Set(addr, b)

		o.SetStatus(Negative, IsNegative(b))
		o.SetStatus(Zero, b == 0)
		o.SetStatus(Carry, carry)
	case BIT_ZP:
		// bit zero page
		o.Log("I: BIT ")
		addr, _ := o.ZeroPage(io)
		b, _ := io.Get(addr)
		o.Log("%02x (ZP)", addr)

		a := o.A
		v := a & b

		o.SetStatus(Overflow, IsOverflow(b, v))
		o.SetStatus(Zero, v == 0)
		o.SetStatus(Negative, IsNegative(b))
	case BIT_A:
		// bit absolute
		o.Log("I: BIT ")
		addr, _ := o.Absolute(io)
		b, _ := io.Get(addr)
		o.Log("%04x (ABS)", addr)

		a := o.A
		v := a & b

		o.SetStatus(Overflow, IsOverflow(b, v))
		o.SetStatus(Zero, v == 0)
		o.SetStatus(Negative, IsNegative(b))

	// Stack Instructions
	case TXS:
		// transfer x to stack
		o.Log("I: TXS ")
		o.SP--
		io.Set(STACK_HEAD+uint16(o.SP), o.X)
	case TSX:
		// transfer stack to x
		o.Log("I: TSX ")
		x, _ := io.Get(STACK_HEAD + uint16(o.SP))
		o.X = x
		o.SP++
	case PHA:
		// push accumulater
		o.Log("I: PHA ")
		io.Set(STACK_HEAD+uint16(o.SP), o.A)
		o.SP--
	case PLA:
		// pull accumulater
		o.Log("I: PLA ")
		o.SP++
		a, _ := io.Get(STACK_HEAD + uint16(o.SP))
		o.A = a
	case PHP:
		// push status to stack
		o.Log("I: PHP ")
		io.Set(STACK_HEAD+uint16(o.SP), o.Status)
		o.SP--
	case PLP:
		// pull status from stack
		o.Log("I: PLP ")
		o.SP++
		status, _ := io.Get(STACK_HEAD + uint16(o.SP))
		o.Status = status

	case DEBUG:
		o.Log("I: DEBUG ")
		bp, _ := io.Get(o.PC)
		o.PC++
		o.Log("%02x (halted)", bp)
		halted = true
		fmt.Printf("HALTED\n")
	}

	o.Log("\n") // always end the instructions debug lines
	o.halted = halted
	return halted, nil // TODO: return errors?
}

func (o *CPU) IsHalted() bool {
	return o.halted
}

func (o *CPU) SetStatus(flag uint8, value bool) {
	var status uint8
	if value {
		// o.Status |= 0b00000001 // bitset
		status = o.Status | flag
	} else {
		// o.Status &^ 0b01000000 // bitclear
		status = o.Status &^ flag
	}

	//o.Log(" Status: %08b > %08b | %08b %v", o.Status, status, flag, value)

	o.Status = status
}

func (o *CPU) Log(format string, a ...any) {
	if !o.DebugMode && !o.SingleStep {
		return
	}
	fmt.Printf(format, a...)
}

func (o *CPU) Debug() {
	if !o.DebugMode {
		return
	}
	//o.Log("State: %v", o)
	fmt.Print("\n\nPC    SP  A    X    Y    Status     \n")
	fmt.Print("-------------------------NV-BDIZC- ($SS)\n")
	//          PC    SP    A      X      Y      Status
	fmt.Printf("%04x  %02x  %02x   %02x   %02x   %08b  ($%02x)\n\n",
		o.PC,
		o.SP,
		o.A,
		o.X,
		o.Y,
		o.Status,
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

func (o *CPU) DebugRegister(screen *ebiten.Image, font font.Face, bound image.Rectangle) {
	x := 0
	y := 0 + bound.Dy()

	s := ""
	s += "PC    SP  A    X    Y    Status     \n"
	s += "-------------------------NV-BDIZC- ($SS)\n"
	//          PC    SP    A      X      Y      Status
	s += fmt.Sprintf("%04x  %02x  %02x   %02x   %02x   %08b  ($%02x)\n\n",
		o.PC,
		o.SP,
		o.A,
		o.X,
		o.Y,
		o.Status,
		o.Status,
	)

	text.Draw(screen, s, font, x, y, color.White)
}

// func (o *CPU) Write(io *Memory.Memory, b uint8) error {
// 	err := io.Set(o.Address, b)
// 	return err
// }

/*
  Addressing Modes
*/

func (o *CPU) Immediate(io IO.Memory) (uint16, error) {
	o.Address = o.PC
	o.PC += 1
	return o.Address, nil
}

func (o *CPU) ZeroPage(io IO.Memory) (uint16, error) {
	zp, err := io.Get(o.PC)
	o.PC += 1
	o.Address = uint16(zp)
	return o.Address, err
}

func (o *CPU) ZeroPageX(io IO.Memory) (uint16, error) {
	zp, err := io.Get(o.PC)
	o.PC += 1
	o.Address = uint16(zp + o.X)
	return o.Address, err
}

func (o *CPU) ZeroPageY(io IO.Memory) (uint16, error) {
	zp, err := io.Get(o.PC)
	o.PC += 1
	o.Address = uint16(zp + o.Y)
	return o.Address, err
}

func (o *CPU) Absolute(io IO.Memory) (uint16, error) {
	addr, err := io.GetWord(o.PC)
	o.PC += 2
	o.Address = addr
	return o.Address, err
}

func (o *CPU) AbsoluteX(io IO.Memory) (uint16, error) {
	addr, err := io.GetWord(o.PC)
	o.PC += 2
	o.Address = addr + uint16(o.X)
	return o.Address, err
}

func (o *CPU) AbsoluteY(io IO.Memory) (uint16, error) {
	addr, err := io.GetWord(o.PC)
	o.PC += 2
	o.Address = addr + uint16(o.Y)
	return o.Address, err
}

func (o *CPU) Indirect(io IO.Memory) (uint16, error) {
	from, _ := io.GetWord(o.PC)
	o.PC += 2
	addr, err := io.GetWord(from)
	o.Address = addr
	return o.Address, err
}

func (o *CPU) IndirectX(io IO.Memory) (uint16, error) {
	zp, _ := io.Get(o.PC)
	o.PC += 1
	addr, err := io.GetWord(uint16(zp + o.X))
	o.Address = addr
	return o.Address, err
}

func (o *CPU) IndirectY(io IO.Memory) (uint16, error) {
	zp, _ := io.Get(o.PC)
	o.PC += 1
	addr1, err := io.GetWord(uint16(zp))
	addr := addr1 + uint16(o.Y)
	o.Address = addr
	return o.Address, err
}

/*
 * Math
 */

func (o *CPU) ADC(operand uint8) uint8 {
	overflow := true
	if ((o.A ^ operand) & 0x80) > 0 {
		overflow = false
	}

	var carry uint8 = 0
	if BitTest(Carry, o.Status) {
		carry = 1
	}

	sum := int16(o.A) + int16(operand) + int16(carry)
	o.Log("\n  ADC: %02x + %02x + %02x = %02x (%v)\n", o.A, operand, carry, sum, sum)

	a := uint8(sum & 0x00ff)

	o.SetStatus(Negative, IsNegative(a))

	if overflow && sum < 0x80 {
		overflow = false
	}
	o.SetStatus(Overflow, overflow)
	o.SetStatus(Zero, a == 0)
	o.SetStatus(Carry, sum >= 0x100)

	o.A = a
	return a
}

func (o *CPU) SBC(operand uint8) uint8 {
	return o.ADC(^operand)
}
