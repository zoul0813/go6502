package main

import (
	"fmt"
	"os"

	"github.com/zoul0813/go6502/pkg/Display"
	"github.com/zoul0813/go6502/pkg/IO"
	"github.com/zoul0813/go6502/pkg/Keyboard"
	"github.com/zoul0813/go6502/pkg/Memory"
)

func main() {

	// RAM
	ram := Memory.New(0x8000, 0x0000, false)

	// Display
	display := Display.New(0xD012, 40, 25)

	// Keyboard
	keyboard := Keyboard.New(0xD010)

	// ROM
	rom := Memory.New(0x1000, 0xF000, true)

	devices := []*IO.Device{
		IO.NewDevice("RAM", ram, 0x0000),
		IO.NewDevice("Keyboard", keyboard, 0xD010),
		IO.NewDevice("Display", display, 0xD012),
		IO.NewDevice("ROM", rom, 0xF000),
	}
	io := IO.New(devices)

	io.List()

	f, err := os.ReadFile("../rom/rom.bin")
	if err != nil {
		fmt.Printf("Can't read file rom/rom.bin\n\n")
		panic(err)
	}
	fmt.Printf("Loading rom $%04x (%v) bytes at %v ($%04x)\n", len(f), len(f), 0xF000, 0xF000)
	io.LoadRom(f, 0xF000)
	io.Dump(0xF000, 0xFF)

	f, err = os.ReadFile("../rom/rom.cfg")
	if err != nil {
		fmt.Printf("Can't read file rom/rom.cfg\n\n")
		panic(err)
	}
	fmt.Printf("Loading ram $%04x (%v) bytes at %v ($%04x)\n", len(f), len(f), 0x0000, 0x0000)
	io.LoadRom(f, 0x0000)
	io.Dump(0x0000, 0xFF)
}
