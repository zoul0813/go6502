package main

import (
	"github.com/zoul0813/go6502/pkg/CPU"
	"github.com/zoul0813/go6502/pkg/IO"
	"github.com/zoul0813/go6502/pkg/Memory"
	"testing"
)

func TestResumeFromSingleStep(t *testing.T) {
	previousCPU, previousIO := cpu, io
	t.Cleanup(func() { cpu, io = previousCPU, previousIO })
	memory := Memory.New(0xffff, 0, false)
	io = IO.New([]*IO.Device{IO.NewDevice("RAM", memory, 0)})
	memory.Set(0x200, byte(CPU.NOP))
	cpu = CPU.New(0x200, 0xff, 0, 0, 0, 0, true, false)
	if err := automaticStep(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x200 {
		t.Fatal("paused CPU executed")
	}
	cpu.SingleStep = false
	if err := automaticStep(); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x201 {
		t.Fatalf("resumed PC=%04x", cpu.PC)
	}
}
