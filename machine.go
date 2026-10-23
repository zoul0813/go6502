package main

import (
	"github.com/zoul0813/go6502/pkg/CPU"
	"github.com/zoul0813/go6502/pkg/Display"
	"github.com/zoul0813/go6502/pkg/IO"
	"github.com/zoul0813/go6502/pkg/Keyboard"
	"github.com/zoul0813/go6502/pkg/Memory"
	monitor "github.com/zoul0813/go6502/rom/monitor"
)

func initializeMachine(singleStep, debugMode bool) error {
	ram = Memory.New(0x7fff, 0, false)
	keyboard = Keyboard.New(0xd010)
	display = Display.New(0xd012, cols, rows)
	rom = Memory.New(0xff, monitor.Address, true)
	if _, err := rom.Load(monitor.Wozmon()); err != nil {
		return err
	}
	io = IO.New([]*IO.Device{
		IO.NewDevice("RAM", ram, 0),
		IO.NewDevice("Keyboard", keyboard, 0xd010),
		IO.NewDevice("Display", display, 0xd012),
		IO.NewDevice("Wozmon", rom, monitor.Address),
	})
	cpu = CPU.New(0, 0xfd, 0, 0, 0, CPU.Reserved|CPU.I, singleStep, debugMode)
	return resetCPU()
}

// Reset the CPU through the hardware vector while preserving RAM contents.
func resetCPU() error {
	pc, err := io.GetWord(0xfffc)
	if err != nil {
		return err
	}
	cpu.PC = pc
	cpu.SP = 0xfd
	cpu.Status = (cpu.Status &^ CPU.B) | CPU.Reserved | CPU.I
	return nil
}
