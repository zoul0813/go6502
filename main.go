package main

import (
	"fmt"
	"go6502/CPU"
	"go6502/Memory"
	"os"
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

	rom.Dump(0x0000, 0xff)
	rom.Dump(0xfff0, 0x0f)
	rom.Dump(0x8000, 0xff)

	// should just loop infinitely now ...

	cpu := CPU.New(
		0xfffc,     //PC
		0x00,       //SP
		0x00,       //A
		0xf0,       //X
		0xFE,       //Y
		0b00000000, //Status
	)

	word, _ := rom.GetWord(cpu.PC)
	cpu.PC = word

	cpu.Debug()

	for {
		// DEBUG CONSOLE
		fmt.Printf("%%: ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "zp":
			rom.Dump(0x0000, 0xff)
			continue
			// rom.Dump(0xfff0, 0x0f)
			// rom.Dump(0x8000, 0xff)
		case "stack":
			rom.Dump(0x0100, 0xff)
			continue
		case "d":
			fallthrough
		case "debug":
			cpu.Debug()
			continue
		case "db":
			fallthrough
		case "debug:bit":
			cpu.Debug()
			cpu.DebugBits()
			continue
		case "help":
			fmt.Printf("Press enter to cycle the next clock tick\n")
			fmt.Printf("\n")
			fmt.Printf("zp = mem dump of zero page\n")
			continue
		case "q":
			fallthrough
		case "quit":
			os.Exit(0)
		}

		fmt.Printf("\n")

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
			word, _ := rom.GetWord(cpu.PC)
			fmt.Printf("%04x (Indirect)", word)
			addr, _ := rom.GetWord(word)
			cpu.PC = addr

		// Misc
		case CPU.BRK: // TODO: NMI
			// break
			fmt.Printf("I: BRK ")
			// TODO: non-maskable interrupt
			cpu.PC += 2
		case CPU.NOP:
			fmt.Printf("I: NOP")
		// Add (ADC)
		case CPU.ADC_I: // add with carry, immediate
		case CPU.ADC_ZP: // add with carry, zero page
		case CPU.ADC_ZPX: // add with carry, zero page, x
		case CPU.ADC_A: // add with carry, absolute
		case CPU.ADC_AX: // add with carry, absolute, x
		case CPU.ADC_AY: // add with carry, absolute, y
		case CPU.ADC_INX: // add with carry, indirect, x
		case CPU.ADC_INY: // add with carry, indirect, y

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
			// wtf is indirect?
		case CPU.AND_INY:
			// and with a, indirect, y
			// wtf is indirect?
		// Store Instructions (STA, STX, STY)
		case CPU.STA_ZP:
			// store a, zero page
			fmt.Printf("I: STA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%04x (ZP)", zp)
			rom.Set(uint16(zp), cpu.A)
		case CPU.STA_ZPX:
			// store a, zero page, x
			fmt.Printf("I: STA ")
			zp, _ := rom.Get(cpu.PC)
			addr := uint16(zp + cpu.X)
			cpu.PC++
			fmt.Printf("%04x (ZP, X)", zp)
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
			// wtf is indirect?
		case CPU.STA_INY: // store a, indirect, y
			// wtf is indirect?
		case CPU.STX_ZP:
			// store a, zero page
			fmt.Printf("I: STX ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%04x (ZP)", zp)
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
			fmt.Printf("%04x (ZP)", zp)
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
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.Y))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, CPU.IsNegative(cpu.A))
		case CPU.LDA_INY:
			// load A indirect, y index
			fmt.Printf("I: LDA ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.X))
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
		}

		fmt.Print("\n") // always end the instructions debug lines
		cpu.Debug()
	}
}
