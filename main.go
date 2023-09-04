package main

import (
	"fmt"
	"go6502/CPU"
	"go6502/Memory"
	"os"
	"strconv"
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
			word, _ := rom.GetWord(cpu.PC)
			fmt.Printf("%04x (ABS)", word)
			cpu.PC = word
		case CPU.JMP_IN:
			// jump indirect
			fmt.Printf("I: JMP ")
			zp, _ := rom.Get(cpu.PC)
			fmt.Printf("%04x (Indirect)", zp)
			addr, _ := rom.GetWord(uint16(zp))
			cpu.PC = addr
		case CPU.JSR_A:
			// jump to subroute, absolute
			fmt.Printf("I: JSR ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS)", addr)

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
			neg := CPU.BitSet(CPU.Negative, uint16(cpu.Status))
			if !neg {
				var offset int8 = int8(rel)
				cpu.PC += uint16(offset)
			}
		case CPU.BMI:
			// branch on minus
			fmt.Printf("I: BMI ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			neg := CPU.BitSet(CPU.Negative, uint16(cpu.Status))
			if neg {
				var offset int8 = int8(rel)
				cpu.PC += uint16(offset)
			}
		case CPU.BVC:
			// branch on overflow clear
			fmt.Printf("I: BVC ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			overflow := CPU.BitSet(CPU.Overflow, uint16(cpu.Status))
			if !overflow {
				var offset int8 = int8(rel)
				cpu.PC += uint16(offset)
			}
		case CPU.BVS:
			// branch on overflow set
			fmt.Printf("I: BVS ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			overflow := CPU.BitSet(CPU.Overflow, uint16(cpu.Status))
			if overflow {
				var offset int8 = int8(rel)
				cpu.PC += uint16(offset)
			}
		case CPU.BCC:
			// branch on carry clear
			fmt.Printf("I: BCC ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			if !carry {
				var offset int8 = int8(rel)
				cpu.PC += uint16(offset)
			}
		case CPU.BCS:
			// branch on carry set
			fmt.Printf("I: BCS ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			if carry {
				var offset int8 = int8(rel)
				cpu.PC += uint16(offset)
			}
		case CPU.BNE:
			// branch on not equal
			fmt.Printf("I: BNE ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			zero := CPU.BitSet(CPU.Zero, uint16(cpu.Status))
			if !zero {
				var offset int8 = int8(rel)
				cpu.PC += uint16(offset)
			}
		case CPU.BEQ:
			// branch on equal
			fmt.Printf("I: BEQ ")
			rel, _ := rom.Get(cpu.PC)
			cpu.PC++
			zero := CPU.BitSet(CPU.Zero, uint16(cpu.Status))
			if zero {
				var offset int8 = int8(rel)
				cpu.PC += uint16(offset)
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
			v, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate)", v)
			a := cpu.A
			cpu.A += v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.ADC_ZP:
			// add with carry, zero page
			fmt.Printf("I: ADC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			v, _ := rom.Get(uint16(zp))
			a := cpu.A
			cpu.A += v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.ADC_ZPX:
			// add with carry, zero page, x
			fmt.Printf("I: ADC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, X)", zp)
			v, _ := rom.Get(uint16(zp + cpu.X))
			a := cpu.A
			cpu.A += v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.ADC_A:
			// add with carry, absolute
			fmt.Printf("I: ADC ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS)", addr)
			v, _ := rom.Get(addr)
			a := cpu.A
			cpu.A += v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.ADC_AX:
			// add with carry, absolute, x
			fmt.Printf("I: ADC ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS, X)", addr)
			v, _ := rom.Get(addr + uint16(cpu.X))
			a := cpu.A
			cpu.A += v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.ADC_AY:
			// add with carry, absolute, y
			fmt.Printf("I: ADC ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS, Y)", addr)
			v, _ := rom.Get(addr + uint16(cpu.Y))
			a := cpu.A
			cpu.A += v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.ADC_INX:
			// add with carry, indirect, x
			fmt.Printf("I: ADC ")
			zp, _ := rom.Get(cpu.PC)
			fmt.Printf("%02x (Indirect, X)", zp)
			cpu.PC++
			addr, _ := rom.GetWord(uint16(zp))
			v, _ := rom.Get(addr + uint16(cpu.X))
			a := cpu.A
			cpu.A += v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.ADC_INY:
			// add with carry, indirect, y
			fmt.Printf("I: ADC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, Y)", zp)
			addr, _ := rom.GetWord(uint16(zp))
			v, _ := rom.Get(addr + uint16(cpu.Y))
			a := cpu.A
			cpu.A += v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))

		// Subtract (SBC)
		case CPU.SBC_I:
			// subtract with carry, immediate
			fmt.Printf("I: SBC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			v, _ := rom.Get(uint16(zp))
			a := cpu.A
			cpu.A -= v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.SBC_ZP:
			// subtract with carry, zero page
			fmt.Printf("I: SBC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			v, _ := rom.Get(uint16(zp))
			a := cpu.A
			cpu.A -= v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.SBC_ZPX:
			// subtract with carry, zero page, x
			fmt.Printf("I: SBC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, X)", zp)
			v, _ := rom.Get(uint16(zp + cpu.X))
			a := cpu.A
			cpu.A -= v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.SBC_A:
			// subtract with carry, absolute
			fmt.Printf("I: SBC ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS)", addr)
			v, _ := rom.Get(addr)
			a := cpu.A
			cpu.A -= v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.SBC_AX:
			// subtract with carry, absolute, x
			fmt.Printf("I: SBC ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS, X)", addr)
			v, _ := rom.Get(addr + uint16(cpu.X))
			a := cpu.A
			cpu.A -= v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.SBC_AY:
			// subtract with carry, absolute, y
			fmt.Printf("I: SBC ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS, Y)", addr)
			v, _ := rom.Get(addr + uint16(cpu.Y))
			a := cpu.A
			cpu.A -= v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.SBC_INX:
			// subtract with carry, indirect, x
			fmt.Printf("I: SBC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, X)", zp)
			addr, _ := rom.GetWord(uint16(zp + cpu.X))
			v, _ := rom.Get(addr + uint16(cpu.Y))
			a := cpu.A
			cpu.A -= v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))
		case CPU.SBC_INY:
			// subtract with carry, indirect, y
			fmt.Printf("I: SBC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, Y)", zp)
			addr, _ := rom.GetWord(uint16(zp + cpu.Y))
			v, _ := rom.Get(addr + uint16(cpu.Y))
			a := cpu.A
			cpu.A -= v

			// N V Z C
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(a, cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, CPU.BitSet(CPU.Bit7, uint16(cpu.A)))

		// Compare (A, X, Y)
		case CPU.CMP_I:
			// compare accumulator
			fmt.Printf("I: CMP ")
			v, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate)", v)
			a := cpu.A
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == v)
			cpu.SetStatus(CPU.Carry, cpu.A >= v)
		case CPU.CMP_ZP:
			// compare accumulator, zero page
			fmt.Printf("I: CMP ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			v, _ := rom.Get(uint16(zp))
			a := cpu.A
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == v)
			cpu.SetStatus(CPU.Carry, cpu.A >= v)
		case CPU.CMP_ZPX:
			// compare accumulator, zero page, x
			fmt.Printf("I: CMP ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, X)", zp)
			v, _ := rom.Get(uint16(zp + cpu.X))
			a := cpu.A
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == v)
			cpu.SetStatus(CPU.Carry, cpu.A >= v)
		case CPU.CMP_A:
			// compare accumulator, absolute
			fmt.Printf("I: CMP ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS)", addr)
			v, _ := rom.Get(addr)
			a := cpu.A
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == v)
			cpu.SetStatus(CPU.Carry, cpu.A >= v)
		case CPU.CMP_AX:
			// compare accumulator, absolute, x
			fmt.Printf("I: CMP ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS, X)", addr)
			v, _ := rom.Get(addr + uint16(cpu.X))
			a := cpu.A
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == v)
			cpu.SetStatus(CPU.Carry, cpu.A >= v)
		case CPU.CMP_AY:
			// compare accumulator, absolute, y
			fmt.Printf("I: CMP ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS, Y)", addr)
			v, _ := rom.Get(addr + uint16(cpu.Y))
			a := cpu.A
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == v)
			cpu.SetStatus(CPU.Carry, cpu.A >= v)
		case CPU.CMP_INX:
			// compare accumulator, indirect, x
			fmt.Printf("I: CMP ")
			zp, _ := rom.GetWord(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, X)", zp)
			addr, _ := rom.GetWord(zp)
			v, _ := rom.Get(addr + uint16(cpu.X))
			a := cpu.A
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == v)
			cpu.SetStatus(CPU.Carry, cpu.A >= v)
		case CPU.CMP_INY:
			// compare accumulator, indirect y
			fmt.Printf("I: CMP ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, Y)", zp)
			addr, _ := rom.GetWord(uint16(zp))
			v, _ := rom.Get(addr + uint16(cpu.Y))
			a := cpu.A
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.A == v)
			cpu.SetStatus(CPU.Carry, cpu.A >= v)
		case CPU.CPX:
			// compare x
			fmt.Printf("I: CPX ")
			v, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate)", v)
			a := cpu.X
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.X == v)
			cpu.SetStatus(CPU.Carry, cpu.X >= v)
		case CPU.CPX_ZP:
			// compare x, zero page
			fmt.Printf("I: CPX ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			v, _ := rom.Get(uint16(zp))
			a := cpu.X
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.X == v)
			cpu.SetStatus(CPU.Carry, cpu.X >= v)
		case CPU.CPX_A:
			// compare x, absolute
			fmt.Printf("I: CPX ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS)", addr)
			v, _ := rom.Get(addr)
			a := cpu.X
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.X == v)
			cpu.SetStatus(CPU.Carry, cpu.X >= v)
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
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS)", addr)
			v, _ := rom.Get(addr)
			a := cpu.Y
			r := a - v // actually do the math, so we can determine if it's negative
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(r))
			cpu.SetStatus(CPU.Zero, cpu.Y == v)
			cpu.SetStatus(CPU.Carry, cpu.Y >= v)

		// AND
		case CPU.AND_I:
			// and with a
			fmt.Printf("I: AND ")
			v, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate)", v)
			cpu.A &= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_ZP:
			// and with a, zero page
			fmt.Printf("I: AND ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			v, _ := rom.Get(uint16(zp))
			cpu.A &= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_ZPX:
			// and with a, zero page, x
			fmt.Printf("I: AND ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, X)", zp)
			v, _ := rom.Get(uint16(zp + cpu.X))
			cpu.A &= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_A:
			// and with a, absolute
			fmt.Printf("I: AND ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ABS)", addr)
			cpu.PC += 2
			v, _ := rom.Get(addr)
			cpu.A &= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_AX:
			// and with a, absolute, x
			fmt.Printf("I: AND ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS, X)", addr)
			cpu.PC += 2
			v, _ := rom.Get(addr + uint16(cpu.X))
			cpu.A &= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_AY:
			// and with a, absolute, y
			fmt.Printf("I: AND ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%02x (ABS, Y)", addr)
			cpu.PC += 2
			v, _ := rom.Get(addr + uint16(cpu.Y))
			cpu.A &= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_INX:
			// and with a, indirect, x
			fmt.Printf("I: AND ")
			zp, _ := rom.GetWord(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, X)", zp)
			addr, _ := rom.GetWord(zp)
			cpu.PC += 2
			v, _ := rom.Get(addr + uint16(cpu.X))
			cpu.A &= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.AND_INY:
			// and with a, indirect, y
			fmt.Printf("I: AND ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, Y)", zp)
			addr, _ := rom.GetWord(uint16(zp))
			cpu.PC += 2
			v, _ := rom.Get(addr + uint16(cpu.Y))
			cpu.A &= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)

		// EOR
		case CPU.EOR_I:
			// exclusive or, immediate
			fmt.Printf("I: EOR ")
			v, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate)", v)
			cpu.A ^= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_ZP:
			// exclusive or, zero page
			fmt.Printf("I: EOR ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			v, _ := rom.Get(uint16(zp))
			fmt.Printf("%02x (ZP)", v)
			cpu.A ^= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_ZPX:
			// exclusive or, zeor page, x
			fmt.Printf("I: EOR ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			v, _ := rom.Get(uint16(zp + cpu.X))
			fmt.Printf("%02x (ZP, X)", v)
			cpu.A ^= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_A:
			// exclusive or, absolute
			fmt.Printf("I: EOR ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			v, _ := rom.Get(addr)
			fmt.Printf("%02x (ABS)", v)
			cpu.A ^= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_AX:
			// exclusive or, absolute, x
			fmt.Printf("I: EOR ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			v, _ := rom.Get(addr + uint16(cpu.X))
			fmt.Printf("%02x (ABS, X)", v)
			cpu.A ^= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_AY:
			// exclusive or, absolute, y
			fmt.Printf("I: EOR ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			v, _ := rom.Get(addr + uint16(cpu.Y))
			fmt.Printf("%02x (ABS, Y)", v)
			cpu.A ^= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_INX:
			// exclusive or, indirect, x
			fmt.Printf("I: EOR ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, X)", zp)
			addr, _ := rom.GetWord(uint16(zp))
			v, _ := rom.Get(addr + uint16(cpu.X))
			cpu.A ^= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.EOR_INY:
			// exclusive or, indirect, y
			fmt.Printf("I: EOR ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, Y)", zp)
			addr, _ := rom.GetWord(uint16(zp))
			v, _ := rom.Get(addr + uint16(cpu.Y))
			cpu.A ^= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)

			// ORA
		case CPU.ORA_I:
			// exclusive or, immediate
			fmt.Printf("I: ORA ")
			v, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate)", v)
			cpu.A |= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_ZP:
			// exclusive or, zero page
			fmt.Printf("I: ORA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			v, _ := rom.Get(uint16(zp))
			fmt.Printf("%02x (ZP)", v)
			cpu.A |= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_ZPX:
			// exclusive or, zeor page, x
			fmt.Printf("I: ORA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			v, _ := rom.Get(uint16(zp + cpu.X))
			fmt.Printf("%02x (ZP, X)", v)
			cpu.A |= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_A:
			// exclusive or, absolute
			fmt.Printf("I: ORA ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			v, _ := rom.Get(addr)
			fmt.Printf("%02x (ABS)", v)
			cpu.A |= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_AX:
			// exclusive or, absolute, x
			fmt.Printf("I: ORA ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			v, _ := rom.Get(addr + uint16(cpu.X))
			fmt.Printf("%02x (ABS, X)", v)
			cpu.A |= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_AY:
			// exclusive or, absolute, y
			fmt.Printf("I: ORA ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			v, _ := rom.Get(addr + uint16(cpu.Y))
			fmt.Printf("%02x (ABS, Y)", v)
			cpu.A |= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_INX:
			// exclusive or, indirect, x
			fmt.Printf("I: ORA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, X)", zp)
			addr, _ := rom.GetWord(uint16(zp))
			v, _ := rom.Get(addr + uint16(cpu.X))
			cpu.A |= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
		case CPU.ORA_INY:
			// exclusive or, indirect, y
			fmt.Printf("I: ORA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, Y)", zp)
			addr, _ := rom.GetWord(uint16(zp))
			v, _ := rom.Get(addr + uint16(cpu.Y))
			cpu.A |= v
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)

		// Store Instructions (STA, STX, STY)
		case CPU.STA_ZP:
			// store a, zero page
			fmt.Printf("I: STA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			rom.Set(uint16(zp), cpu.A)
		case CPU.STA_ZPX:
			// store a, zero page, x
			fmt.Printf("I: STA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, X)", zp)
			addr := uint16(zp + cpu.X)
			fmt.Printf(" %02x %02x %04x", zp, cpu.X, addr)
			rom.Set(addr, cpu.A)
		case CPU.STA_A:
			// store a, zero page, x
			fmt.Printf("I: STA ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS)", addr)
			rom.Set(addr, cpu.A)
		case CPU.STA_AX:
			// store a, zero page, x
			fmt.Printf("I: STA ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			addr += uint16(cpu.X)
			fmt.Printf("%04x (ABS, X)", addr)
			rom.Set(addr, cpu.A)
		case CPU.STA_AY:
			// store a, zero page, x
			fmt.Printf("I: STA ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			addr += uint16(cpu.Y)
			fmt.Printf("%04x (ABS, Y)", addr)
			rom.Set(addr, cpu.A)
		case CPU.STA_INX: // store a, indirect, x
			fmt.Printf("I: STA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%04x (Indirect, X)", zp)
			addr := uint16(zp + cpu.X)
			fmt.Printf(" %02x %02x %04x", zp, cpu.X, addr)
			rom.Set(addr, cpu.A)
		case CPU.STA_INY: // store a, indirect, y
			fmt.Printf("I: STA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Indirect, Y)", zp)
			addr := uint16(zp + cpu.Y)
			rom.Set(addr, cpu.A)
		case CPU.STX_ZP:
			// store a, zero page
			fmt.Printf("I: STX ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			rom.Set(uint16(zp), cpu.X)
		case CPU.STX_ZPY:
			// store a, zero page, x
			fmt.Printf("I: STX ")
			zp, _ := rom.Get(cpu.PC)
			addr := uint16(zp + cpu.X)
			cpu.PC++
			fmt.Printf("%04x (ZP, Y)", zp)
			rom.Set(addr, cpu.X)
		case CPU.STX_A:
			// store a, zero page, x
			fmt.Printf("I: STX ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS)", addr)
			rom.Set(addr, cpu.X)
		case CPU.STY_ZP:
			// store a, zero page
			fmt.Printf("I: STY ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			rom.Set(uint16(zp), cpu.Y)
		case CPU.STY_ZPX:
			// store a, zero page, x
			fmt.Printf("I: STY ")
			zp, _ := rom.Get(cpu.PC)
			addr := uint16(zp + cpu.X)
			cpu.PC++
			fmt.Printf("%04x (ZP, X)", zp)
			rom.Set(addr, cpu.Y)
		case CPU.STY_A:
			// store a, zero page, x
			fmt.Printf("I: STY ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS)", addr)
			rom.Set(addr, cpu.Y)

		// INC/DEC Instructions
		case CPU.INC_ZP:
			// increment zero page
			fmt.Printf("I: INC ")
			zp, _ := rom.Get(cpu.PC)
			addr := uint16(zp)
			cpu.PC++
			b, _ = rom.Get(addr)
			fmt.Printf("%02x (ZP)", b)
			b++
			rom.Set(addr, b)
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.INC_ZPX:
			// increment zero page, x
			fmt.Printf("I: INC ")
			zp, _ := rom.Get(cpu.PC)
			addr := uint16(zp + cpu.X)
			cpu.PC++
			b, _ = rom.Get(addr)
			fmt.Printf("%02x (ZP, X)", b)
			b++
			rom.Set(addr, b)
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.INC_A:
			// increment absolute
			fmt.Printf("I: INC ")
			word, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			b, _ = rom.Get(word)
			fmt.Printf("%02x (ABS)", b)
			b++
			rom.Set(word, b)
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.INC_AX:
			// increment absolute, x
			fmt.Printf("I: INC ")
			word, _ := rom.GetWord(cpu.PC)
			addr := word + uint16(cpu.X)
			cpu.PC += 2
			b, _ = rom.Get(addr)
			fmt.Printf("%02x (ABS, X)", b)
			b++
			rom.Set(addr, b)
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.DEC_ZP:
			// increment zero page
			fmt.Printf("I: INC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			b, _ = rom.Get(uint16(zp))
			fmt.Printf("%02x (ZP)", b)
			b--
			rom.Set(uint16(zp), b)
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.DEC_ZPX:
			// increment zero page, x
			fmt.Printf("I: INC ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			b, _ = rom.Get(uint16(zp + cpu.X))
			fmt.Printf("%02x (ZP, X)", b)
			b--
			rom.Set(uint16(zp), b)
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.DEC_A:
			// increment absolute
			fmt.Printf("I: INC ")
			word, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			b, _ = rom.Get(word)
			fmt.Printf("%02x (ABS)", b)
			b--
			rom.Set(word, b)
			cpu.SetStatus(CPU.Zero, b == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(b))
		case CPU.DEC_AX:
			// increment absolute, x
			fmt.Printf("I: INC ")
			word, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			b, _ = rom.Get(word + uint16(cpu.X))
			fmt.Printf("%02x (ABS, X)", b)
			b--
			rom.Set(word, b)
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
			b, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate)", b)
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_ZP:
			// load A zero page
			fmt.Printf("I: LDA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			b, _ = rom.Get(uint16(zp))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_ZPX:
			// load A zero page, x index
			fmt.Printf("I: LDA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, X)", zp)
			b, _ = rom.Get(uint16(zp + cpu.X))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_A:
			// load A absolute
			fmt.Printf("I: LDA ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS)", a)
			b, _ := rom.Get(a)
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_AX:
			// load A absolute, x
			fmt.Printf("I: LDA ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.X))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_AY:
			// load A absolute, y
			fmt.Printf("I: LDA ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.Y))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_IX:
			// load A x index, indirect
			fmt.Printf("I: LDA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%04x (Indirect, X)", zp)
			addr := uint16(zp + cpu.X)
			fmt.Printf(" %02x %02x %04x", zp, cpu.X, addr)
			b, _ := rom.Get(addr)
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_INY:
			// load A indirect, y index
			fmt.Printf("I: LDA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%04x (Indirect, Y)", zp)
			addr := uint16(zp + cpu.Y)
			fmt.Printf(" %02x %02x %04x", zp, cpu.Y, addr)
			b, _ := rom.Get(addr)
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDX_I:
			// load X immediate
			fmt.Printf("I: LDX ")
			b, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate))", b)
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDX_ZP:
			// load X zero page
			fmt.Printf("I: LDX ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			b, _ = rom.Get(uint16(zp))
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDX_ZPY:
			// load X zero page, y index
			fmt.Printf("I: LDX ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, Y)", zp)
			b, _ = rom.Get(uint16(zp + cpu.Y))
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDX_A:
			// load X absolute
			fmt.Printf("I: LDX ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS)", a)
			b, _ := rom.Get(a)
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDX_AY:
			// load X absolute, y index
			fmt.Printf("I: LDX ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, Y)", a)
			b, _ := rom.Get(a + uint16(cpu.Y))
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.X))
		case CPU.LDY_I:
			// load Y immediate
			fmt.Printf("I: LDY ")
			b, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate))", b)
			cpu.Y = b
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.LDY_ZP:
			// load Y zero page
			fmt.Printf("I: LDY ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			b, _ = rom.Get(uint16(zp))
			cpu.Y = b
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.LDY_ZPX:
			// load Y zero page, x index
			fmt.Printf("I: LDY ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, Y)", zp)
			b, _ = rom.Get(uint16(zp + cpu.Y))
			cpu.Y = b
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.LDY_A:
			// load Y absolute
			fmt.Printf("I: LDY ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS)", a)
			b, _ := rom.Get(a)
			cpu.Y = b
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.Y))
		case CPU.LDY_AX:
			// load Y absolute, x index
			fmt.Printf("I: LDY ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.Y))
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
			panic("Not yet implemented ... ")

		// Bit Shift Instructions
		case CPU.ROL:
			// rotate left, a
			fmt.Printf("I: ROL ")
			fmt.Printf("%02x (A)", cpu.A)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b7 := CPU.BitSet(cpu.A, CPU.Bit7)
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
			zp, _ := rom.Get(cpu.PC)
			fmt.Printf("%02x (ZP)", zp)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b7 := CPU.BitSet(zp, CPU.Bit7)
			v := (zp << 1)
			if carry {
				v++
			}
			rom.Set(cpu.PC, v)
			cpu.PC++ // increment it after we've set the value

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ROL_ZPX:
			// rotate left, zerp page, x
			fmt.Printf("I: ROL ")
			zp, _ := rom.Get(cpu.PC + uint16(cpu.X))
			fmt.Printf("%02x (ZP, X)", zp)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b7 := CPU.BitSet(zp, CPU.Bit7)
			v := (zp << 1)
			if carry {
				v++
			}
			rom.Set(cpu.PC, v)
			cpu.PC++ // increment it after we've set the value

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ROL_A:
			// rotate left, absolute
			fmt.Printf("I: ROL ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			abs, _ := rom.Get(addr)
			fmt.Printf("%02x (ABS)", abs)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b7 := CPU.BitSet(abs, CPU.Bit7)
			v := (abs << 1)
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
			addr, _ := rom.GetWord(cpu.PC + uint16(cpu.X))
			cpu.PC += 2
			abs, _ := rom.Get(addr)
			fmt.Printf("%02x (ABS, X)", abs)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b7 := CPU.BitSet(abs, CPU.Bit7)
			v := (abs << 1)
			if carry {
				v++
			}
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ROR:
			// rotate left, a
			fmt.Printf("I: ROL ")
			fmt.Printf("%02x (A)", cpu.A)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b0 := CPU.BitSet(cpu.A, CPU.Bit0)
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
			fmt.Printf("I: ROL ")
			zp, _ := rom.Get(cpu.PC)
			fmt.Printf("%02x (ZP)", zp)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b0 := CPU.BitSet(zp, CPU.Bit0)
			v := (zp >> 1)
			if carry {
				v |= 0b10000000 // bitset
			}
			rom.Set(cpu.PC, v)
			cpu.PC++ // increment it after we've set the value

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b0)
		case CPU.ROR_ZPX:
			// rotate left, zerp page, x
			fmt.Printf("I: ROL ")
			zp, _ := rom.Get(cpu.PC + uint16(cpu.X))
			fmt.Printf("%02x (ZP, X)", zp)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b0 := CPU.BitSet(zp, CPU.Bit0)
			v := (zp >> 1)
			if carry {
				v |= 0b10000000 // bitset
			}
			rom.Set(cpu.PC, v)
			cpu.PC++ // increment it after we've set the value

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b0)
		case CPU.ROR_A:
			// rotate left, absolute
			fmt.Printf("I: ROL ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			abs, _ := rom.Get(addr)
			fmt.Printf("%02x (ABS)", abs)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b0 := CPU.BitSet(abs, CPU.Bit0)
			v := (abs >> 1)
			if carry {
				v |= 0b10000000 // bitset
			}
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b0)
		case CPU.ROR_AX:
			// rotate left, absolute, x
			fmt.Printf("I: ROL ")
			addr, _ := rom.GetWord(cpu.PC + uint16(cpu.X))
			cpu.PC += 2
			abs, _ := rom.Get(addr)
			fmt.Printf("%02x (ABS, X)", abs)
			carry := CPU.BitSet(CPU.Carry, uint16(cpu.Status))
			b0 := CPU.BitSet(abs, CPU.Bit0)
			v := (abs >> 1)
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
			fmt.Printf("%02x (A)", cpu.A)
			b7 := CPU.BitSet(cpu.A, CPU.Bit7)
			a := (cpu.A << 1)
			cpu.A = a

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ASL_ZP:
			// arithmetic shift left, zero page
			fmt.Printf("I: ROL ")
			zp, _ := rom.Get(cpu.PC)
			fmt.Printf("%02x (ZP)", zp)
			b7 := CPU.BitSet(zp, CPU.Bit7)
			v := (zp << 1)
			rom.Set(cpu.PC, v)
			cpu.PC++ // increment it after we've set the value

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ASL_ZPX:
			// arithmetic shift left, zero page, x
			fmt.Printf("I: ROL ")
			zp, _ := rom.Get(cpu.PC + uint16(cpu.X))
			fmt.Printf("%02x (ZP, X)", zp)
			b7 := CPU.BitSet(zp, CPU.Bit7)
			v := (zp << 1)
			rom.Set(cpu.PC, v)
			cpu.PC++ // increment it after we've set the value

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ASL_A:
			// arithmetic shift left, absolute
			fmt.Printf("I: ROL ")
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			abs, _ := rom.Get(addr)
			fmt.Printf("%02x (ABS)", abs)
			b7 := CPU.BitSet(abs, CPU.Bit7)
			v := (abs << 1)
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.ASL_AX:
			// arithmetic shift left, absolute, x
			fmt.Printf("I: ROL ")
			addr, _ := rom.GetWord(cpu.PC + uint16(cpu.X))
			cpu.PC += 2
			abs, _ := rom.Get(addr)
			fmt.Printf("%02x (ABS, X)", abs)
			b7 := CPU.BitSet(abs, CPU.Bit7)
			v := (abs << 1)
			rom.Set(addr, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, b7)
		case CPU.LSR:
			// logical shift right, a
			fmt.Printf("I: LSR ")
			carry := CPU.BitSet(CPU.Bit1, uint16(cpu.A))
			cpu.A = cpu.A >> 1

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Carry, carry)
		case CPU.LSR_ZP:
			// logical shift right, zero page
			zp, _ := rom.Get(cpu.PC)
			carry := CPU.BitSet(CPU.Bit1, uint16(zp))
			rom.Set(cpu.PC, zp)
			cpu.PC++

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(zp))
			cpu.SetStatus(CPU.Zero, zp == 0)
			cpu.SetStatus(CPU.Carry, carry)
		case CPU.LSR_A:
			// logical shift right, absolute
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			v, _ := rom.Get(addr)
			carry := CPU.BitSet(CPU.Bit1, uint16(v))
			rom.Set(cpu.PC, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(v))
			cpu.SetStatus(CPU.Zero, v == 0)
			cpu.SetStatus(CPU.Carry, carry)
		case CPU.LSR_AX:
			// logical shift right, absolute, x
			addr, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			v, _ := rom.Get(addr + uint16(cpu.X))
			carry := CPU.BitSet(CPU.Bit1, uint16(v))
			rom.Set(cpu.PC, v)

			cpu.SetStatus(CPU.Negative, CPU.IsNegative(v))
			cpu.SetStatus(CPU.Zero, v == 0)
			cpu.SetStatus(CPU.Carry, carry)
		case CPU.BIT_ZP:
			// bit zero page
			fmt.Printf("I: BIT ")
			zp, _ := rom.Get(cpu.PC)
			fmt.Printf("%02x (ZP)", zp)
			cpu.PC++
			a := cpu.A
			v := a & zp

			cpu.SetStatus(CPU.Overflow, CPU.IsOverflow(zp, v))
			cpu.SetStatus(CPU.Zero, v == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(zp))
		case CPU.BIT_A:
			// bit absolute
			fmt.Printf("I: BIT ")
			addr, _ := rom.GetWord(cpu.PC)
			b, _ := rom.Get(addr)
			cpu.PC += 2
			fmt.Printf("%02x (ABS)", addr)
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
			fmt.Print("I: DEBUG (halted) \n")
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
			b := 0xFC
			var c int8 = int8(b - 127)
			var d int8 = int8(b)
			fmt.Printf("%02x %02x %08b %v %v\n", b, c, c, c, d)
		}

		// fmt.Printf("\n")
	}
exitDebugConsole:
	fmt.Printf("\n")
}
