package main

import (
	"github.com/zoul0813/go6502/pkg/CPU"
	"github.com/zoul0813/go6502/pkg/IO"
	"github.com/zoul0813/go6502/pkg/Memory"
	"testing"
	"time"
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

func TestFrameBudgetAndResume(t *testing.T) {
	previousCPU, previousIO, previousClock := cpu, io, clockSpeed
	t.Cleanup(func() { cpu, io, clockSpeed = previousCPU, previousIO, previousClock })
	memory := Memory.New(0xffff, 0, false)
	for i := uint16(0); i < 16; i++ {
		memory.Set(0x200+i, byte(CPU.NOP))
	}
	io = IO.New([]*IO.Device{IO.NewDevice("RAM", memory, 0)})
	cpu = CPU.New(0x200, 0xff, 0, 0, 0, 0, true, false)
	clockSpeed = time.Millisecond
	g := &Game{}
	if err := g.advance(time.Second); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x200 {
		t.Fatal("paused CPU executed")
	}
	cpu.SingleStep = false
	for i := 0; i < 3; i++ {
		if err := g.advance(time.Millisecond / 2); err != nil {
			t.Fatal(err)
		}
	}
	if cpu.PC != 0x201 {
		t.Fatalf("PC=%04x, want 0201", cpu.PC)
	}
	if err := g.advance(time.Millisecond / 2); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x202 {
		t.Fatalf("PC=%04x, want 0202", cpu.PC)
	}
}

func TestClockInterval(t *testing.T) {
	for _, tt := range []struct {
		rate int
		hz   bool
		want time.Duration
	}{
		{1000, false, time.Microsecond}, {1, false, time.Millisecond}, {1, true, time.Second}, {1000, true, time.Millisecond},
	} {
		got, err := clockInterval(tt.rate, tt.hz)
		if err != nil || got != tt.want {
			t.Fatalf("rate=%d hz=%v got=%v err=%v", tt.rate, tt.hz, got, err)
		}
	}
	for _, rate := range []int{-1, 0, 1_000_001} {
		if _, err := clockInterval(rate, false); err == nil {
			t.Fatalf("invalid rate %d accepted", rate)
		}
	}
}
