package main

import (
	"strings"
	"testing"
)

func testMachine(t *testing.T) {
	t.Helper()
	oldCPU, oldIO, oldRAM, oldROM, oldKeyboard, oldDisplay := cpu, io, ram, rom, keyboard, display
	t.Cleanup(func() { cpu, io, ram, rom, keyboard, display = oldCPU, oldIO, oldRAM, oldROM, oldKeyboard, oldDisplay })
	if err := initializeMachine(false, false); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0xff00 {
		t.Fatalf("reset PC=%04x", cpu.PC)
	}
}

func monitorIdle(t *testing.T) {
	t.Helper()
	for i := 0; i < 200000; i++ {
		ready, err := io.Get(0xd011)
		if err != nil {
			t.Fatal(err)
		}
		if cpu.PC == 0xff29 && ready&0x80 == 0 {
			return
		}
		if halted, err := cpu.Step(io); halted || err != nil {
			t.Fatalf("PC=%04x halted=%v err=%v", cpu.PC, halted, err)
		}
	}
	t.Fatalf("monitor did not return to keyboard poll: PC=%04x display=%q", cpu.PC, display.All(false))
}
func monitorCommand(t *testing.T, command string) {
	t.Helper()
	for _, b := range []byte(command + "\r") {
		keyboard.AppendKey(b)
	}
	monitorIdle(t)
}
func TestOriginalWozmonBootStoreExamineAndRun(t *testing.T) {
	testMachine(t)
	monitorIdle(t)
	if got := display.All(false); !strings.HasPrefix(got, "\\") {
		t.Fatalf("missing startup prompt: %q", got)
	}
	monitorCommand(t, "03FE: DE AD BE EF")
	for i, want := range []byte{0xde, 0xad, 0xbe, 0xef} {
		if got, _ := io.Get(0x3fe + uint16(i)); got != want {
			t.Fatalf("store at %04x got=%02x want=%02x", 0x3fe+i, got, want)
		}
	}
	monitorCommand(t, "03FE.0401")
	output := display.All(false)
	for _, want := range []string{"03FE: DE AD", "0400: BE EF"} {
		if !strings.Contains(output, want) {
			t.Fatalf("missing %q in %q", want, output)
		}
	}
	// Enter and execute a program through the monitor. It calls the ROM ECHO
	// entry twice, then returns to the monitor keyboard poll.
	monitorCommand(t, "0300: A9 C1 20 EF FF A9 DA 20 EF FF 4C 29 FF")
	monitorCommand(t, "0300R")
	if got := display.All(false); !strings.Contains(got, "AZ") {
		t.Fatalf("program did not print AZ: %q", got)
	}
	if err := io.Set(0xff00, 0); err == nil {
		t.Fatal("monitor ROM is writable")
	}
	monitorCommand(t, "0000.00FF") // More output than a full 40x24 screen.
	if got := display.All(false); !strings.Contains(got, "00F8:") {
		t.Fatalf("scrolling lost final dump line: %q", got)
	}
}
func TestMonitorResetAndEscape(t *testing.T) {
	testMachine(t)
	monitorIdle(t)
	monitorCommand(t, "0300: 42")
	if err := resetCPU(); err != nil {
		t.Fatal(err)
	}
	monitorIdle(t)
	if value, _ := io.Get(0x300); value != 0x42 {
		t.Fatal("reset cleared RAM")
	}
	keyboard.AppendKey(0x1b)
	monitorIdle(t)
	monitorCommand(t, "0300")
	if got := display.All(false); !strings.Contains(got, "0300: 42") {
		t.Fatalf("monitor unusable after reset/escape: %q", got)
	}
}
