package main

import (
	"fmt"
	"go6502/CPU"
	"go6502/Memory"
	"os"
)

const ROM_HEAD = 0x8000
const ZP_HEAD = 0x000

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
		0x0000,     //SP
		0x00,       //A
		0xf0,       //X
		0xFE,       //Y
		0b00000000, //Status
	)

	word, _ := rom.GetWord(cpu.PC)
	cpu.PC = word

	cpu.Debug()

	for {
		fmt.Printf("%%: ")
		var input string
		fmt.Scanln(&input)
		switch input {
		case "zp":
			rom.Dump(0x0000, 0xff)
			continue
			// rom.Dump(0xfff0, 0x0f)
			// rom.Dump(0x8000, 0xff)
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
		case CPU.BRK:
			// break
			fmt.Printf("I: BRK ")
			// TODO: non-maskable interrupt
			cpu.PC++
		case CPU.ADC_I: // add with carry, immediate
		case CPU.ADC_ZP: // add with carry, zero page
		case CPU.ADC_ZPX: // add with carry, zero page, x
		case CPU.ADC_A: // add with carry, absolute
		case CPU.ADC_AX: // add with carry, absolute, x
		case CPU.ADC_AY: // add with carry, absolute, y
		case CPU.ADC_INX: // add with carry, indirect, x
		case CPU.ADC_INY: // add with carry, indirect, y
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
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(b))
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
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(b))
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
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(b))
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
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(b))
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
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(b))
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
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(b))
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
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(b))
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
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(b))
		case CPU.INX:
			// increment x
			fmt.Printf("I: INX ")
			cpu.X++
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			// cpu.SetStatus(CPU.Overflow, cpu.X < x)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.X))
		case CPU.DEX:
			// increment x
			fmt.Printf("I: DEX ")
			cpu.X--
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.X))
		case CPU.INY:
			// increment y
			fmt.Printf("I: INY ")
			cpu.Y++
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.Y))
		case CPU.DEY:
			// increment y
			fmt.Printf("I: INY ")
			cpu.Y--
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.Y))
		case CPU.LDA_I:
			// load A immediate
			fmt.Printf("I: LDA ")
			b, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate)", b)
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		case CPU.LDA_ZP:
			// load A zero page
			fmt.Printf("I: LDA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			b, _ = rom.Get(uint16(zp))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		case CPU.LDA_ZPX:
			// load A zero page, x index
			fmt.Printf("I: LDA ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, X)", zp)
			b, _ = rom.Get(uint16(zp + cpu.X))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		case CPU.LDA_A:
			// load A absolute
			fmt.Printf("I: LDA ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS)", a)
			b, _ := rom.Get(a)
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		case CPU.LDA_AX:
			// load A absolute, x
			fmt.Printf("I: LDA ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.X))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		case CPU.LDA_AY:
			// load A absolute, y
			fmt.Printf("I: LDA ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.Y))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		case CPU.LDA_IX:
			// load A x index, indirect
			fmt.Printf("I: LDA ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.Y))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		case CPU.LDA_INY:
			// load A indirect, y index
			fmt.Printf("I: LDA ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.X))
			cpu.A = b
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		case CPU.LDX_I:
			// load X immediate
			fmt.Printf("I: LDX ")
			b, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate))", b)
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.X))
		case CPU.LDX_ZP:
			// load X zero page
			fmt.Printf("I: LDX ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			b, _ = rom.Get(uint16(zp))
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.X))
		case CPU.LDX_ZPY:
			// load X zero page, y index
			fmt.Printf("I: LDX ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, Y)", zp)
			b, _ = rom.Get(uint16(zp + cpu.Y))
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.X))
		case CPU.LDX_A:
			// load X absolute
			fmt.Printf("I: LDX ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS)", a)
			b, _ := rom.Get(a)
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.X))
		case CPU.LDX_AY:
			// load X absolute, y index
			fmt.Printf("I: LDX ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, Y)", a)
			b, _ := rom.Get(a + uint16(cpu.Y))
			cpu.X = b
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.X))
		case CPU.LDY_I:
			// load Y immediate
			fmt.Printf("I: LDY ")
			b, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate))", b)
			cpu.Y = b
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.Y))
		case CPU.LDY_ZP:
			// load Y zero page
			fmt.Printf("I: LDY ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP)", zp)
			b, _ = rom.Get(uint16(zp))
			cpu.Y = b
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.Y))
		case CPU.LDY_ZPX:
			// load Y zero page, x index
			fmt.Printf("I: LDY ")
			zp, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (ZP, Y)", zp)
			b, _ = rom.Get(uint16(zp + cpu.Y))
			cpu.Y = b
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.Y))
		case CPU.LDY_A:
			// load Y absolute
			fmt.Printf("I: LDY ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS)", a)
			b, _ := rom.Get(a)
			cpu.Y = b
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.Y))
		case CPU.LDY_AX:
			// load Y absolute, x index
			fmt.Printf("I: LDY ")
			a, _ := rom.GetWord(cpu.PC)
			cpu.PC += 2
			fmt.Printf("%04x (ABS, X)", a)
			b, _ := rom.Get(a + uint16(cpu.Y))
			cpu.Y = b
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.Y))
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
		case CPU.TAX:
			// transfer a to x
			fmt.Printf("I: TAX ")
			cpu.X = cpu.A
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.X))
		case CPU.TXA:
			// transfer x to a
			fmt.Printf("I: TXA ")
			cpu.A = cpu.X
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		case CPU.TAY:
			// transfer a to y
			fmt.Printf("I: TAY ")
			cpu.Y = cpu.A
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.Y))
		case CPU.TYA:
			// transfer y to a
			fmt.Printf("I: TYA ")
			cpu.A = cpu.Y
			cpu.SetStatus(CPU.Zero, cpu.A == 0)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.A))
		}

		fmt.Print("\n") // always end the instructions debug lines
		cpu.Debug()
	}
}
