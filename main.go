package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zoul0813/go6502/pkg/CPU"
	"github.com/zoul0813/go6502/pkg/Memory"
)

const ROM_HEAD = 0x8000
const ZP_HEAD = 0x000
const STACK_HEAD = 0x100

func main() {
	fmt.Printf("Go 6502... \n")

	// ram := Memory.New(0x8000, 0x0000, false)
	rom := Memory.New(0xffff, 0x0000, false)
	f, err := os.ReadFile("rom/rom.bin")
	if err != nil {
		fmt.Printf("Can't read file rom/rom.bin\n\n")
		panic(err)
	}

	rom.Load(f)

	// rom.Dump(0x0000, 0xff) // Zero Page
	rom.Dump(0xfff0, 0x0f) // Reset Vectors
	rom.Dump(0x8000, 0xff) // Start of Program?

	// should just loop infinitely now ...

	cpu := CPU.New(
		0xfffc,     // PC
		0x00,       // SP
		0x00,       // A
		0xf0,       // X
		0xFE,       // Y
		0b00000000, // Status
		false,      // Single Step
	)

	word, _ := rom.GetWord(cpu.PC)
	cpu.PC = word

	cpu.Debug()

	for {
		if cpu.SingleStep {
			DebugConsole(cpu, rom)
		}

		b, _ := rom.Get(cpu.PC)
		var instr CPU.OpCode = CPU.OpCode(b)
		cpu.PC++ // increment the stack pointer
		fmt.Printf("Instruction: %02x\n", instr)
		switch instr {
		// Jump/Branch Instructions
		case CPU.JMP_A:
			// jump absolute
			fmt.Printf("I: JMP ")
			addr, _ := cpu.Absolute(rom)
			fmt.Printf("%04x (ABS)", addr)
			cpu.PC = addr
		case CPU.JMP_IN:
			// jump indirect
			fmt.Printf("I: JMP ")
			// from, _ := rom.GetWord(cpu.PC)
			// addr, _ := rom.GetWord(from)
			addr, _ := cpu.Indirect(rom)
			fmt.Printf("%04x (Indirect)", addr)
			cpu.PC = addr
		case CPU.JSR_A:
			// jump to subroute, absolute
			fmt.Printf("I: JSR ")
			// addr, _ := rom.GetWord(cpu.PC) // jump to here
			// cpu.PC += 2
			addr, _ := cpu.Absolute(rom)
			fmt.Printf("%04x (ABS)", addr)

			pc := cpu.PC - 1
			var lo uint8 = uint8(pc & 0b0000000011111111)
			var hi uint8 = uint8(pc >> 8)

			cpu.SP--
			rom.Set(STACK_HEAD+uint16(cpu.SP), hi)
			cpu.SP--
			rom.Set(STACK_HEAD+uint16(cpu.SP), lo)

			cpu.PC = addr
		case CPU.RTS:
			fmt.Printf("I: RTS ")
			fmt.Printf("(Implied)")
			lo, _ := rom.Get(STACK_HEAD + uint16(cpu.SP))
			cpu.SP++
			hi, _ := rom.Get(STACK_HEAD + uint16(cpu.SP))
			cpu.SP++
			var addr uint16 = (uint16(hi) << 8) | uint16(lo)
			cpu.PC = addr + 1
		case CPU.RTI:
			fmt.Printf("I: RTI ")
			fmt.Printf("(Implied)")
			status, _ := rom.Get(STACK_HEAD + uint16(cpu.SP))
			cpu.SP++
			cpu.Status = status
			lo, _ := rom.Get(STACK_HEAD + uint16(cpu.SP))
			cpu.SP++
			hi, _ := rom.Get(STACK_HEAD + uint16(cpu.SP))
			cpu.SP++
			var addr uint16 = (uint16(hi) << 8) | uint16(lo)
			cpu.PC = addr
		case CPU.BPL:
			// branch on plus
			fmt.Printf("I: BPL ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Rel)", rel)
			neg := CPU.BitTest(CPU.Negative, cpu.Status)
			if !neg {
				cpu.PC += uint16(rel)
			}
		case CPU.BMI:
			// branch on minus
			fmt.Printf("I: BMI ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Rel)", rel)
			neg := CPU.BitTest(CPU.Negative, cpu.Status)
			if neg {
				cpu.PC += uint16(rel)
			}
		case CPU.BVC:
			// branch on overflow clear
			fmt.Printf("I: BVC ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Rel)", rel)
			overflow := CPU.BitTest(CPU.Overflow, cpu.Status)
			if !overflow {
				cpu.PC += uint16(rel)
			}
		case CPU.BVS:
			// branch on overflow set
			fmt.Printf("I: BVS ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Rel)", rel)
			overflow := CPU.BitTest(CPU.Overflow, cpu.Status)
			if overflow {
				cpu.PC += uint16(rel)
			}
		case CPU.BCC:
			// branch on carry clear
			fmt.Printf("I: BCC ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Rel)", rel)
			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			if !carry {
				cpu.PC += uint16(rel)
			}
		case CPU.BCS:
			// branch on carry set
			fmt.Printf("I: BCS ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Rel)", rel)
			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			if carry {
				cpu.PC += uint16(rel)
			}
		case CPU.BNE:
			// branch on not equal
			fmt.Printf("I: BNE ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Rel)", rel)
			zero := CPU.BitTest(CPU.Zero, cpu.Status)
			if !zero {
				cpu.PC += uint16(rel)
			}
		case CPU.BEQ:
			// branch on equal
			fmt.Printf("I: BEQ ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Rel)", rel)
			zero := CPU.BitTest(CPU.Zero, cpu.Status)
			if zero {
				cpu.PC += uint16(rel)
			}

		// Misc
		case CPU.BRK: // TODO: NMI
			// break
			fmt.Printf("I: BRK ")
			// TODO: non-maskable interrupt
			cpu.PC += 2
		case CPU.NOP:
			fmt.Printf("I: NOP")

		// Add (ADC)
		case CPU.ADC_I:
			// add with carry, immediate
			fmt.Printf("I: ADC ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate)", b)

			cpu.ADC(b)
			// N V Z C
		case CPU.ADC_ZP:
			// add with carry, zero page
			fmt.Printf("I: ADC ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			cpu.ADC(b)
			// N V Z C
		case CPU.ADC_ZPX:
			// add with carry, zero page, x
			fmt.Printf("I: ADC ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			cpu.ADC(b)
			// N V Z C
		case CPU.ADC_A:
			// add with carry, absolute
			fmt.Printf("I: ADC ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			cpu.ADC(b)
			// N V Z C
		case CPU.ADC_AX:
			// add with carry, absolute, x
			fmt.Printf("I: ADC ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			cpu.ADC(b)
			// N V Z C
		case CPU.ADC_AY:
			// add with carry, absolute, y
			fmt.Printf("I: ADC ")
			addr, _ := cpu.AbsoluteY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, Y)", addr)

			cpu.ADC(b)
			// N V Z C
		case CPU.ADC_INX:
			// add with carry, indirect, x
			fmt.Printf("I: ADC ")
			addr, _ := cpu.IndirectX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, X)", addr)

			cpu.ADC(b)
			// N V Z C
		case CPU.ADC_INY:
			// add with carry, indirect, y
			fmt.Printf("I: ADC ")
			addr, _ := cpu.IndirectY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, Y)", addr)

			cpu.ADC(b)
			// N V Z C

		// Subtract (SBC)
		case CPU.SBC_I:
			// subtract with carry, immediate
			fmt.Printf("I: SBC ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate)", b)

			cpu.SBC(b)
			// N V Z C
		case CPU.SBC_ZP:
			// subtract with carry, zero page
			fmt.Printf("I: SBC ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			cpu.SBC(b)
			// N V Z C
		case CPU.SBC_ZPX:
			// subtract with carry, zero page, x
			fmt.Printf("I: SBC ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			cpu.SBC(b)
			// N V Z C
		case CPU.SBC_A:
			// subtract with carry, absolute
			fmt.Printf("I: SBC ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			cpu.SBC(b)
			// N V Z C
		case CPU.SBC_AX:
			// subtract with carry, absolute, x
			fmt.Printf("I: SBC ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			cpu.SBC(b)
			// N V Z C
		case CPU.SBC_AY:
			// subtract with carry, absolute, y
			fmt.Printf("I: SBC ")
			addr, _ := cpu.AbsoluteY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, Y)", addr)

			cpu.SBC(b)
			// N V Z C
		case CPU.SBC_INX:
			// subtract with carry, indirect, x
			fmt.Printf("I: SBC ")
			addr, _ := cpu.IndirectX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, X)", addr)

			cpu.SBC(b)
			// N V Z C
		case CPU.SBC_INY:
			// subtract with carry, indirect, y
			fmt.Printf("I: SBC ")
			addr, _ := cpu.IndirectY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, Y)", addr)

			cpu.SBC(b)
			// N V Z C

		// Compare (A, X, Y)
		case CPU.CMP_I:
			// compare accumulator
			fmt.Printf("I: CMP ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate)", b)

			a := cpu.A
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == b)
			cpu.SetStatus(CPU.Carry, cpu.A >= b)
		case CPU.CMP_ZP:
			// compare accumulator, zero page
			fmt.Printf("I: CMP ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			a := cpu.A
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == b)
			cpu.SetStatus(CPU.Carry, cpu.A >= b)
		case CPU.CMP_ZPX:
			// compare accumulator, zero page, x
			fmt.Printf("I: CMP ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			a := cpu.A
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == b)
			cpu.SetStatus(CPU.Carry, cpu.A >= b)
		case CPU.CMP_A:
			// compare accumulator, absolute
			fmt.Printf("I: CMP ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			a := cpu.A
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == b)
			cpu.SetStatus(CPU.Carry, cpu.A >= b)
		case CPU.CMP_AX:
			// compare accumulator, absolute, x
			fmt.Printf("I: CMP ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			a := cpu.A
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == b)
			cpu.SetStatus(CPU.Carry, cpu.A >= b)
		case CPU.CMP_AY:
			// compare accumulator, absolute, y
			fmt.Printf("I: CMP ")
			addr, _ := cpu.AbsoluteY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, Y)", addr)

			a := cpu.A
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == b)
			cpu.SetStatus(CPU.Carry, cpu.A >= b)
		case CPU.CMP_INX:
			// compare accumulator, indirect, x
			fmt.Printf("I: CMP ")
			addr, _ := cpu.IndirectX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, X)", addr)

			a := cpu.A
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == b)
			cpu.SetStatus(CPU.Carry, cpu.A >= b)
		case CPU.CMP_INY:
			// compare accumulator, indirect y
			fmt.Printf("I: CMP ")
			addr, _ := cpu.IndirectY(rom)
			b, _ := rom.Get(addr)

			fmt.Printf("%04x (Indirect, Y)", addr)

			a := cpu.A
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == b)
			cpu.SetStatus(CPU.Carry, cpu.A >= b)
		case CPU.CPX:
			// compare x
			fmt.Printf("I: CPX ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate)", b)

			a := cpu.X
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.X == b)
			cpu.SetStatus(CPU.Carry, cpu.X >= b)
		case CPU.CPX_ZP:
			// compare x, zero page
			fmt.Printf("I: CPX ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			a := cpu.X
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.X == b)
			cpu.SetStatus(CPU.Carry, cpu.X >= b)
		case CPU.CPX_A:
			// compare x, absolute
			fmt.Printf("I: CPX ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			a := cpu.X
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.X == b)
			cpu.SetStatus(CPU.Carry, cpu.X >= b)
		case CPU.CPY:
			// compare y
			fmt.Printf("I: CPY ")
			v, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate)", v)
			a := cpu.Y
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.Y == v)
			cpu.SetStatus(CPU.Carry, cpu.Y >= v)
		case CPU.CPY_ZP:
			// compare x, zero page
			fmt.Printf("I: CPY ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			v, _ := rom.Get(uint16(zp))
			a := cpu.Y
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.Y == v)
			cpu.SetStatus(CPU.Carry, cpu.Y >= v)
		case CPU.CPY_A:
			// compare x, absolute
			fmt.Printf("I: CPY ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			a := cpu.Y
			r := a - b // actually do the math, so we can determine if it's negative

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.Y == b)
			cpu.SetStatus(CPU.Carry, cpu.Y >= b)

		// AND
		case CPU.AND_I:
			// and with a
			fmt.Printf("I: AND ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate)", b)

			cpu.A &= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_ZP:
			// and with a, zero page
			fmt.Printf("I: AND ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			cpu.A &= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_ZPX:
			// and with a, zero page, x
			fmt.Printf("I: AND ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			cpu.A &= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_A:
			// and with a, absolute
			fmt.Printf("I: AND ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			cpu.A &= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_AX:
			// and with a, absolute, x
			fmt.Printf("I: AND ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			cpu.A &= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_AY:
			// and with a, absolute, y
			fmt.Printf("I: AND ")
			addr, _ := cpu.AbsoluteY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, Y)", addr)

			cpu.A &= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_INX:
			// and with a, indirect, x
			fmt.Printf("I: AND ")
			addr, _ := cpu.IndirectX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, X)", addr)

			cpu.A &= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_INY:
			// and with a, indirect, y
			fmt.Printf("I: AND ")
			addr, _ := cpu.IndirectY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, Y)", addr)

			cpu.A &= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)

		// EOR
		case CPU.EOR_I:
			// exclusive or, immediate
			fmt.Printf("I: EOR ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate)", b)

			cpu.A ^= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_ZP:
			// exclusive or, zero page
			fmt.Printf("I: EOR ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			cpu.A ^= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_ZPX:
			// exclusive or, zeor page, x
			fmt.Printf("I: EOR ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			cpu.A ^= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_A:
			// exclusive or, absolute
			fmt.Printf("I: EOR ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			cpu.A ^= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_AX:
			// exclusive or, absolute, x
			fmt.Printf("I: EOR ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			cpu.A ^= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_AY:
			// exclusive or, absolute, y
			fmt.Printf("I: EOR ")
			addr, _ := cpu.AbsoluteY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, Y)", addr)

			cpu.A ^= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_INX:
			// exclusive or, indirect, x
			fmt.Printf("I: EOR ")
			addr, _ := cpu.IndirectX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, X)", addr)
			cpu.A ^= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_INY:
			// exclusive or, indirect, y
			fmt.Printf("I: EOR ")
			addr, _ := cpu.IndirectY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, Y)", addr)
			cpu.A ^= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)

			// ORA
		case CPU.ORA_I:
			// exclusive or, immediate
			fmt.Printf("I: ORA ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate)", b)

			cpu.A |= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_ZP:
			// exclusive or, zero page
			fmt.Printf("I: ORA ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			cpu.A |= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_ZPX:
			// exclusive or, zeor page, x
			fmt.Printf("I: ORA ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			cpu.A |= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_A:
			// exclusive or, absolute
			fmt.Printf("I: ORA ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			cpu.A |= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_AX:
			// exclusive or, absolute, x
			fmt.Printf("I: ORA ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			cpu.A |= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_AY:
			// exclusive or, absolute, y
			fmt.Printf("I: ORA ")
			addr, _ := cpu.AbsoluteY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, Y)", addr)

			cpu.A |= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_INX:
			// exclusive or, indirect, x
			fmt.Printf("I: ORA ")
			addr, _ := cpu.IndirectX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, X)", addr)
			cpu.A |= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_INY:
			// exclusive or, indirect, y
			fmt.Printf("I: ORA ")
			addr, _ := cpu.IndirectY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, Y)", addr)

			cpu.A |= b

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)

		// Store Instructions (STA, STX, STY)
		case CPU.STA_ZP:
			// store a, zero page
			fmt.Printf("I: STA ")
			addr, _ := cpu.ZeroPage(rom)
			fmt.Printf("%02x (ZP)", addr)

			rom.Set(addr, cpu.A)
		case CPU.STA_ZPX:
			// store a, zero page, x
			fmt.Printf("I: STA ")
			addr, _ := cpu.ZeroPageX(rom)
			fmt.Printf("%02x (ZP, X)", addr)

			rom.Set(addr, cpu.A)
		case CPU.STA_A:
			// store a, zero page, x
			fmt.Printf("I: STA ")
			addr, _ := cpu.Absolute(rom)
			fmt.Printf("%04x (ABS)", addr)

			rom.Set(addr, cpu.A)
		case CPU.STA_AX:
			// store a, zero page, x
			fmt.Printf("I: STA ")
			addr, _ := cpu.AbsoluteX(rom)
			fmt.Printf("%04x (ABS, X)", addr)

			rom.Set(addr, cpu.A)
		case CPU.STA_AY:
			// store a, zero page, x
			fmt.Printf("I: STA ")
			addr, _ := cpu.AbsoluteY(rom)
			fmt.Printf("%04x (ABS, Y)", addr)

			rom.Set(addr, cpu.A)
		case CPU.STA_INX: // store a, indirect, x
			fmt.Printf("I: STA ")
			addr, _ := cpu.IndirectX(rom)
			fmt.Printf("%04x (Indirect, X)", addr)

			rom.Set(addr, cpu.A)
		case CPU.STA_INY: // store a, indirect, y
			fmt.Printf("I: STA ")
			addr, _ := cpu.IndirectY(rom)
			fmt.Printf("%04x (Indirect, Y)", addr)

			rom.Set(addr, cpu.A)
		case CPU.STX_ZP:
			// store a, zero page
			fmt.Printf("I: STX ")
			addr, _ := cpu.ZeroPage(rom)
			fmt.Printf("%02x (ZP)", addr)

			rom.Set(addr, cpu.X)
		case CPU.STX_ZPY:
			// store a, zero page, y
			fmt.Printf("I: STY ")
			addr, _ := cpu.ZeroPageY(rom)
			fmt.Printf("%02x (ZP, Y)", addr)

			rom.Set(addr, cpu.X)
		case CPU.STX_A:
			// store a, zero page, x
			fmt.Printf("I: STX ")
			addr, _ := cpu.Absolute(rom)
			fmt.Printf("%04x (ABS)", addr)

			rom.Set(addr, cpu.X)
		case CPU.STY_ZP:
			// store a, zero page
			fmt.Printf("I: STY ")
			addr, _ := cpu.ZeroPage(rom)
			fmt.Printf("%02x (ZP)", addr)

			rom.Set(addr, cpu.Y)
		case CPU.STY_ZPX:
			// store a, zero page, x
			fmt.Printf("I: STY ")
			addr, _ := cpu.ZeroPageX(rom)
			fmt.Printf("%02x (ZP, X)", addr)

			rom.Set(addr, cpu.Y)
		case CPU.STY_A:
			// store a, zero page, x
			fmt.Printf("I: STY ")
			addr, _ := cpu.Absolute(rom)
			fmt.Printf("%04x (ABS)", addr)

			rom.Set(addr, cpu.Y)

		// INC/DEC Instructions
		case CPU.INC_ZP:
			// increment zero page
			fmt.Printf("I: INC ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			b++
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.INC_ZPX:
			// increment zero page, x
			fmt.Printf("I: INC ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			b++
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.INC_A:
			// increment absolute
			fmt.Printf("I: INC ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			b++
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.INC_AX:
			// increment absolute, x
			fmt.Printf("I: INC ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			b++
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.DEC_ZP:
			// increment zero page
			fmt.Printf("I: DEC ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			b--
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.DEC_ZPX:
			// increment zero page, x
			fmt.Printf("I: DEC ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			b--
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.DEC_A:
			// increment absolute
			fmt.Printf("I: DEC ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			b--
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.DEC_AX:
			// increment absolute, x
			fmt.Printf("I: DEC ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			b--
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.INX:
			// increment x
			fmt.Printf("I: INX ")
			cpu.X++
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			// cpu.SetStatus(CPU.Overflow, cpu.X < x)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.DEX:
			// increment x
			fmt.Printf("I: DEX ")
			cpu.X--
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.INY:
			// increment y
			fmt.Printf("I: INY ")
			cpu.Y++
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.DEY:
			// increment y
			fmt.Printf("I: INY ")
			cpu.Y--
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))

		// Load Register Instructions
		case CPU.LDA_I:
			// load A immediate
			fmt.Printf("I: LDA ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate)", b)

			cpu.A = b

			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_ZP:
			// load A zero page
			fmt.Printf("I: LDA ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			cpu.A = b

			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_ZPX:
			// load A zero page, x index
			fmt.Printf("I: LDA ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			cpu.A = b

			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_A:
			// load A absolute
			fmt.Printf("I: LDA ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			cpu.A = b

			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_AX:
			// load A absolute, x
			fmt.Printf("I: LDA ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			cpu.A = b

			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_AY:
			// load A absolute, y
			fmt.Printf("I: LDA ")
			addr, _ := cpu.AbsoluteY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, Y)", addr)

			cpu.A = b

			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_INX:
			// load A x index, indirect
			fmt.Printf("I: LDA ")
			addr, _ := cpu.IndirectX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, X)", addr)

			cpu.A = b

			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_INY:
			// load A indirect, y index
			fmt.Printf("I: LDA ")
			addr, _ := cpu.IndirectY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (Indirect, Y)", addr)

			cpu.A = b

			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDX_I:
			// load X immediate
			fmt.Printf("I: LDX ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate))", b)

			cpu.X = b

			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDX_ZP:
			// load X zero page
			fmt.Printf("I: LDX ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			cpu.X = b

			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDX_ZPY:
			// load X zero page, y index
			fmt.Printf("I: LDX ")
			addr, _ := cpu.ZeroPageY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, Y)", addr)

			cpu.X = b

			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDX_A:
			// load X absolute
			fmt.Printf("I: LDX ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			cpu.X = b

			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDX_AY:
			// load X absolute, y index
			fmt.Printf("I: LDX ")
			addr, _ := cpu.AbsoluteY(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, Y)", addr)

			cpu.X = b

			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDY_I:
			// load Y immediate
			fmt.Printf("I: LDY ")
			addr, _ := cpu.Immediate(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (Immediate))", b)

			cpu.Y = b

			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.LDY_ZP:
			// load Y zero page
			fmt.Printf("I: LDY ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			cpu.Y = b

			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.LDY_ZPX:
			// load Y zero page, x index
			fmt.Printf("I: LDY ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, Y)", addr)

			cpu.Y = b

			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.LDY_A:
			// load Y absolute
			fmt.Printf("I: LDY ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			cpu.Y = b

			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.LDY_AX:
			// load Y absolute, x index
			fmt.Printf("I: LDY ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			cpu.Y = b

			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))

		// Register Transfer Instructions
		case CPU.TAX:
			// transfer a to x
			fmt.Printf("I: TAX ")
			cpu.X = cpu.A
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.TXA:
			// transfer x to a
			fmt.Printf("I: TXA ")
			cpu.A = cpu.X
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.TAY:
			// transfer a to y
			fmt.Printf("I: TAY ")
			cpu.Y = cpu.A
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.TYA:
			// transfer y to a
			fmt.Printf("I: TYA ")
			cpu.A = cpu.Y
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))

		// Status Instructions
		case CPU.CLC:
			// clear carry
			fmt.Printf("I: CLC ")
			cpu.SetStatus(CPU.Carry, false)
		case CPU.SEC:
			// set carry
			fmt.Printf("I: SEC ")
			cpu.SetStatus(CPU.Carry, true)
		case CPU.CLI:
			// clear interrupt
			fmt.Printf("I: CLI ")
			cpu.SetStatus(CPU.Interrupt, false)
		case CPU.SEI:
			// set interrupt
			fmt.Printf("I: SEI ")
			cpu.SetStatus(CPU.Interrupt, true)
		case CPU.CLV:
			// clear overflow
			fmt.Printf("I CLV ")
			cpu.SetStatus(CPU.Overflow, false)
		case CPU.CLD:
			// clear decimal
			fmt.Printf("I CLD ")
			cpu.SetStatus(CPU.Decimal, false)
		case CPU.SED:
			// set decimal
			fmt.Printf("I SED ")
			cpu.SetStatus(CPU.Decimal, true)
			fmt.Printf("Not yet implemented ... ")
			cpu.SetStatus(CPU.Decimal, true)

		// Bit Shift Instructions
		case CPU.ROL:
			// rotate left, a
			fmt.Printf("I: ROL ")
			fmt.Printf("%02x (A)", cpu.A)
			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b7 := CPU.BitTest(cpu.A, CPU.Bit7)
			a := (cpu.A << 1)
			if carry {
				a++
			}
			cpu.A = a

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ROL_ZP:
			// rotate left, zerp page
			fmt.Printf("I: ROL ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b7 := CPU.BitTest(b, CPU.Bit7)
			v := (b << 1)
			if carry {
				v++
			}
			rom.Set(uint16(b), v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ROL_ZPX:
			// rotate left, zerp page, x
			fmt.Printf("I: ROL ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b7 := CPU.BitTest(b, CPU.Bit7)
			v := (b << 1)
			if carry {
				v++
			}
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ROL_A:
			// rotate left, absolute
			fmt.Printf("I: ROL ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b7 := CPU.BitTest(b, CPU.Bit7)
			v := (b << 1)
			if carry {
				v++
			}
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ROL_AX:
			// rotate left, absolute, x
			fmt.Printf("I: ROL ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b7 := CPU.BitTest(b, CPU.Bit7)
			v := (b << 1)
			if carry {
				v++
			}
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ROR:
			// rotate left, a
			fmt.Printf("I: ROR ")
			fmt.Printf("%02x (A)", cpu.A)
			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b0 := CPU.BitTest(cpu.A, CPU.Bit0)
			v := (cpu.A >> 1)
			if carry {
				v |= 0b10000000 // bitset
			}
			cpu.A = v
			fmt.Printf(" %08b", v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b0)
		case CPU.ROR_ZP:
			// rotate left, zerp page
			fmt.Printf("I: ROR ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b0 := CPU.BitTest(b, CPU.Bit0)
			v := (b >> 1)
			if carry {
				v |= 0b10000000 // bitset
			}
			rom.Set(uint16(b), v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b0)
		case CPU.ROR_ZPX:
			// rotate left, zerp page, x
			fmt.Printf("I: ROR ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b0 := CPU.BitTest(b, CPU.Bit0)
			v := (b >> 1)
			if carry {
				v |= 0b10000000 // bitset
			}
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b0)
		case CPU.ROR_A:
			// rotate left, absolute
			fmt.Printf("I: ROR ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)
			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b0 := CPU.BitTest(b, CPU.Bit0)
			v := b >> 1
			if carry {
				v |= 0b10000000 // bitset
			}
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b0)
		case CPU.ROR_AX:
			// rotate left, absolute, x
			fmt.Printf("I: ROR ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			carry := CPU.BitTest(CPU.Carry, cpu.Status)
			b0 := CPU.BitTest(b, CPU.Bit0)
			v := b >> 1
			if carry {
				v |= 0b10000000 // bitset
			}
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b0)
		case CPU.ASL:
			// arithmetic shift left, a
			fmt.Printf("I: ASL ")
			fmt.Printf("%02x", cpu.A)

			b7 := CPU.BitTest(cpu.A, CPU.Bit7)
			a := cpu.A << 1
			cpu.A = a

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ASL_ZP:
			// arithmetic shift left, zero page
			fmt.Printf("I: ASL ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			b7 := CPU.BitTest(b, CPU.Bit7)
			b = b << 1
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ASL_ZPX:
			// arithmetic shift left, zero page, x
			fmt.Printf("I: ASL ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			b7 := CPU.BitTest(b, CPU.Bit7)
			v := b << 1
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ASL_A:
			// arithmetic shift left, absolute
			fmt.Printf("I: ROL ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			b7 := CPU.BitTest(b, CPU.Bit7)
			v := b << 1
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ASL_AX:
			// arithmetic shift left, absolute, x
			fmt.Printf("I: ASL ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			b7 := CPU.BitTest(b, CPU.Bit7)
			v := b << 1
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.LSR:
			// logical shift right, a
			fmt.Printf("I: LSR ")
			carry := CPU.BitTest(CPU.Bit0, cpu.A)
			cpu.A = cpu.A >> 1

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, carry)
		case CPU.LSR_ZP:
			// logical shift right, zero page
			fmt.Printf("I: LSR ")
			// zp, _ := rom.Get(cpu.PC)
			// cpu.PC++
			// b, _ := rom.Get(uint16(zp))
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			carry := CPU.BitTest(CPU.Bit1, b)
			b = b >> 1
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Carry, carry)
		case CPU.LSR_ZPX:
			// logical shift right, zero page
			fmt.Printf("I: LSR ")
			addr, _ := cpu.ZeroPageX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", addr)

			carry := CPU.BitTest(CPU.Bit1, b)
			b = b >> 1
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Carry, carry)
		case CPU.LSR_A:
			// logical shift right, absolute
			fmt.Printf("I: LSR ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			carry := CPU.BitTest(CPU.Bit1, b)
			b = b >> 1
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Carry, carry)
		case CPU.LSR_AX:
			// logical shift right, absolute, x
			fmt.Printf("I: LSR ")
			addr, _ := cpu.AbsoluteX(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS, X)", addr)

			carry := CPU.BitTest(CPU.Bit1, b)
			b = b >> 1
			rom.Set(addr, b)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Carry, carry)
		case CPU.BIT_ZP:
			// bit zero page
			fmt.Printf("I: BIT ")
			addr, _ := cpu.ZeroPage(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%02x (ZP)", addr)

			a := cpu.A
			v := a & b

			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(b, v))
			cpu.SetStatus(CPU.Zero, v == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.BIT_A:
			// bit absolute
			fmt.Printf("I: BIT ")
			addr, _ := cpu.Absolute(rom)
			b, _ := rom.Get(addr)
			fmt.Printf("%04x (ABS)", addr)

			a := cpu.A
			v := a & b

			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(b, v))
			cpu.SetStatus(CPU.Zero, v == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))

		// Stack Instructions
		case CPU.TXS:
			// transfer x to stack
			fmt.Printf("I: TXS ")
			cpu.SP--
			rom.Set(STACK_HEAD+uint16(cpu.SP), cpu.X)
		case CPU.TSX:
			// transfer stack to x
			fmt.Printf("I: TSX ")
			x, _ := rom.Get(STACK_HEAD + uint16(cpu.SP))
			cpu.X = x
			cpu.SP++
		case CPU.PHA:
			// push accumulater
			fmt.Printf("I: PHA ")
			cpu.SP--
			rom.Set(STACK_HEAD+uint16(cpu.SP), cpu.A)
		case CPU.PLA:
			// pull accumulater
			fmt.Printf("I: PLA ")
			a, _ := rom.Get(STACK_HEAD + uint16(cpu.SP))
			cpu.A = a
			cpu.SP++
		case CPU.PHP:
			// push status to stack
			fmt.Printf("I: PHP ")
			cpu.SP--
			rom.Set(STACK_HEAD+uint16(cpu.SP), cpu.Status)
		case CPU.PLP:
			// pull status from stack
			fmt.Printf("I: PLP ")
			status, _ := rom.Get(STACK_HEAD + uint16(cpu.SP))
			cpu.Status = status
			cpu.SP++

		case CPU.DEBUG:
			fmt.Print("I: DEBUG ")
			bp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (halted)", bp)
			fmt.Print("\n")
			DebugConsole(cpu, rom)
		}

		fmt.Print("\n") // always end the instructions debug lines
		cpu.Debug()
	}
}

func DebugConsole(cpu *CPU.CPU, rom *Memory.Memory) {
	// DEBUG CONSOLE
	for {
		var pre string
		if cpu.SingleStep {
			pre = fmt.Sprintf("S ")
		}
		fmt.Printf("%v%%: ", pre)
		var cmd string
		var arg1 string
		var arg2 string
		fmt.Scanln(&cmd, &arg1, &arg2)
		switch cmd {
		case "":
			if cpu.SingleStep {
				fmt.Printf("Single Stepping\n")
				goto exitDebugConsole
			}
		case "c":
			fallthrough
		case "continue":
			goto exitDebugConsole
		case "zp":
			fallthrough
		case "zeropage":
			rom.Dump(0x0000, 0xff)
		case "s":
			fallthrough
		case "stack":
			rom.Dump(0x0100, 0xff)
		case "m":
			fallthrough
		case "mem":
			var start uint16 = 0x00
			var end uint16 = 0xFF
			if len(arg1) > 0 {
				s, _ := strconv.ParseInt(arg1, 16, 16)
				start = uint16(s)
			}
			if len(arg2) > 0 {
				e, _ := strconv.ParseInt(arg2, 16, 16)
				end = uint16(e)
			}
			rom.Dump(start, end)
		case "d":
			fallthrough
		case "debug":
			cpu.Debug()
		case "db":
			fallthrough
		case "debug:bit":
			cpu.Debug()
			cpu.DebugBits()
		case "ss":
			fallthrough
		case "singlestep":
			cpu.SingleStep = !cpu.SingleStep
			fmt.Printf("SingleStep = %v\n", cpu.SingleStep)
		case "h":
			fallthrough
		case "help":
			fmt.Printf("Press enter to cycle the next clock tick\n")
			fmt.Printf("\n")
			fmt.Printf("c|continue            continue execution\n")
			fmt.Printf("zp|zeropage           mem dump of zero page\n")
			fmt.Printf("s|stack               show stack ($0100:$1FF)\n")
			fmt.Printf("m|mem [start, len]    show memory ($start..$len)\n")
			fmt.Printf("d|debug               print registers\n")
			fmt.Printf("db|debug:bit          print registers as bits\n")
			fmt.Printf("ss|singlestep         toggle single step\n")
			fmt.Printf("h|help                this helpful message\n")
			fmt.Printf("\n")
		case "q":
			fallthrough
		case "quit":
			os.Exit(0)
		case "test":
			// special command for just doing quick tests
			// code is volatile
			// b := 0xFC
			var a int8 = 85
			fmt.Printf("A: %08b\n", a)
			a = a & 0x68
			fmt.Printf("&: %08b\n", a)
			a = a | 56
			fmt.Printf("|: %08b\n", a)
			a = a ^ 17
			fmt.Printf("^: %08b\n", a)
			fmt.Printf("A: %08b  %02x %0x\n", a, a, a)
		}

		// fmt.Printf("\n")
	}
exitDebugConsole:
	fmt.Printf("\n")
}
