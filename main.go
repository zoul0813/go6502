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

	cpu.Debug()

	word, _ := rom.GetWord(cpu.PC)
	cpu.PC = word

	cpu.Debug()

	for {

		b, _ := rom.Get(cpu.PC)
		var instr CPU.OpCode = CPU.OpCode(b)
		cpu.PC++ // increment the stack pointer
		fmt.Printf("Instruction: %02x\n", instr)
		switch instr {
		case CPU.JMP_A:
			// jump absolute
			fmt.Printf("I: JMP ")
			word, _ := rom.GetWord(cpu.PC)
			fmt.Printf("%04x", word)
			cpu.PC = word
		case CPU.BRK:
			// break
			fmt.Printf("I: BRK ")
			// TODO: non-maskable interrupt
			cpu.PC++
		case CPU.INX:
			// increment x
			fmt.Printf("I: INX ")
			cpu.X++
			cpu.SetStatus(CPU.Zero, cpu.X == 0)
			// cpu.SetStatus(CPU.Overflow, cpu.X < x)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.X))
		case CPU.INY:
			// increment y
			fmt.Printf("I: INY ")
			cpu.Y++
			cpu.SetStatus(CPU.Zero, cpu.Y == 0)
			// cpu.SetStatus(CPU.Overflow, cpu.Y < y)
			cpu.SetStatus(CPU.Negative, cpu.IsNegative(cpu.Y))
		case CPU.LDA_I:
			// load A immediate
			fmt.Printf("I: LDA ")
			b, _ := rom.Get(cpu.PC)
			cpu.PC++
			fmt.Printf("%02x (Immediate))", b)
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
		case CPU.LDA_AY:
			// load A absolute, y
		case CPU.LDA_IX:
			// load A x index, indirect
		case CPU.LDA_INY:
			// load A indirect, y index
		case CPU.CLC:
			// clear carry
			fmt.Printf("I: CLC ")
			cpu.SetStatus(CPU.Carry, false)
			// cpu.Status &= 0b11111110
		case CPU.SEC:
			// set carry
			fmt.Printf("I: SEC ")
			cpu.SetStatus(CPU.Carry, true)
			// cpu.Status |= 0b00000001
		}

		fmt.Print("\n") // always end the instructions debug lines

		cpu.Debug()
		fmt.Printf("[Enter] for next Clock Tick")
		var input string
		fmt.Scanln(&input)
	}
}
