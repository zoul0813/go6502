package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/zoul0813/go6502/pkg/CPU"
	"strings"
	"testing"
)

func TestTextInputOncePerFrame(t *testing.T) {
	testMachine(t)
	g := &Game{}
	if err := g.handleInput([]ebiten.Key{ebiten.KeyShiftLeft, ebiten.KeyA, ebiten.KeyB}, []rune("abλ😀")); err != nil {
		t.Fatal(err)
	}
	if err := g.handleInput(nil, []rune("c")); err != nil {
		t.Fatal(err)
	} // Held-key repeat.
	for _, want := range []byte{'A', 'B', 'C'} {
		if got, _ := io.Get(0xd010); got != want|0x80 {
			t.Fatalf("got=%02x want=%02x", got, want|0x80)
		}
	}
	if ready, _ := io.Get(0xd011); ready&0x80 != 0 {
		t.Fatal("text was duplicated")
	}
}
func TestMonitorBackspaceAndHome(t *testing.T) {
	testMachine(t)
	monitorIdle(t)
	g := &Game{}
	if err := g.handleInput(nil, []rune("030X")); err != nil {
		t.Fatal(err)
	}
	monitorIdle(t)
	if err := g.handleInput([]ebiten.Key{ebiten.KeyBackspace}, nil); err != nil {
		t.Fatal(err)
	}
	monitorIdle(t)
	if err := g.handleInput(nil, []rune("0: 42")); err != nil {
		t.Fatal(err)
	}
	if err := g.handleInput([]ebiten.Key{ebiten.KeyEnter}, nil); err != nil {
		t.Fatal(err)
	}
	monitorIdle(t)
	if got, _ := io.Get(0x300); got != 0x42 {
		t.Fatalf("backspace command stored %02x", got)
	}
	if err := g.handleInput([]ebiten.Key{ebiten.KeyHome}, nil); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0xff00 {
		t.Fatalf("Home PC=%04x", cpu.PC)
	}
	monitorIdle(t)
	monitorCommand(t, "0300")
	if !strings.Contains(display.All(false), "0300: 42") {
		t.Fatal("Home broke monitor or cleared RAM")
	}
}
func TestFunctionKeyPauseAndStep(t *testing.T) {
	testMachine(t)
	g := &Game{}
	cpu.PC = 0x300
	ram.Set(0x300, byte(CPU.NOP))
	if err := g.handleInput([]ebiten.Key{ebiten.KeyF7}, nil); err != nil {
		t.Fatal(err)
	}
	if !cpu.SingleStep {
		t.Fatal("F7 did not pause")
	}
	if err := g.handleInput([]ebiten.Key{ebiten.KeyF8}, nil); err != nil {
		t.Fatal(err)
	}
	if cpu.PC != 0x301 {
		t.Fatalf("F8 PC=%04x", cpu.PC)
	}
	if err := g.handleInput([]ebiten.Key{ebiten.KeyF7}, nil); err != nil {
		t.Fatal(err)
	}
	if cpu.SingleStep {
		t.Fatal("F7 did not resume")
	}
}
