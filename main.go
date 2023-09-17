package main

import (
	"fmt"
	"os"

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
		0xFF,       // SP
		0x00,       // A
		0xf0,       // X
		0xFE,       // Y
		0b00110000, // Status
		false,      // Single Step
		true,       // DebugMode
	)

	word, _ := rom.GetWord(cpu.PC)
	cpu.PC = word

	cpu.Debug()

	for {
		if cpu.SingleStep {
			DebugConsole(cpu, rom)
		}

		halted, _ := cpu.Step(rom)

		if halted {
			DebugConsole(cpu, rom)
		}

		fmt.Print("\n") // always end the instructions debug lines
		cpu.Debug()
	}
}
