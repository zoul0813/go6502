package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/zoul0813/go6502/pkg/CPU"
	"github.com/zoul0813/go6502/pkg/IO"
)

func DebugConsole(cpu *CPU.CPU, io *IO.IO) {
	// DEBUG CONSOLE
	for {
		var pre string
		if cpu.SingleStep {
			pre = "S "
		}
		fmt.Printf("%v%%: ", pre)
		var cmd string
		var arg1 string
		var arg2 string
		fmt.Scanln(&cmd, &arg1, &arg2)
		switch cmd {
		case "":
			if cpu.SingleStep {
				fmt.Printf("Single Stepping\n")
				goto exitDebugConsole
			}
		case "c":
			fallthrough
		case "continue":
			goto exitDebugConsole
		case "zp":
			fallthrough
		case "zeropage":
			io.Dump(0x0000, 0xff)
		case "s":
			fallthrough
		case "stack":
			io.Dump(0x0100, 0xff)
		case "m":
			fallthrough
		case "mem":
			var start uint16 = 0x00
			var end uint16 = 0xFF
			if len(arg1) > 0 {
				s, _ := strconv.ParseInt(arg1, 16, 16)
				start = uint16(s)
			}
			if len(arg2) > 0 {
				e, _ := strconv.ParseInt(arg2, 16, 16)
				end = uint16(e)
			}
			io.Dump(start, end)
		case "d":
			fallthrough
		case "debug":
			cpu.Debug()
		case "db":
			fallthrough
		case "debug:bit":
			cpu.Debug()
			cpu.DebugBits()
		case "ss":
			fallthrough
		case "singlestep":
			cpu.SingleStep = !cpu.SingleStep
			fmt.Printf("SingleStep = %v\n", cpu.SingleStep)
		case "h":
			fallthrough
		case "help":
			fmt.Printf("Press enter to cycle the next clock tick\n")
			fmt.Printf("\n")
			fmt.Printf("c|continue            continue execution\n")
			fmt.Printf("zp|zeropage           mem dump of zero page\n")
			fmt.Printf("s|stack               show stack ($0100:$1FF)\n")
			fmt.Printf("m|mem [start, len]    show memory ($start..$len)\n")
			fmt.Printf("d|debug               print registers\n")
			fmt.Printf("db|debug:bit          print registers as bits\n")
			fmt.Printf("ss|singlestep         toggle single step\n")
			fmt.Printf("h|help                this helpful message\n")
			fmt.Printf("\n")
		case "q":
			fallthrough
		case "quit":
			os.Exit(0)
		case "test":
			// special command for just doing quick tests
			// code is volatile
			// b := 0xFC
			var a int8 = 85
			fmt.Printf("A: %08b\n", a)
			a = a & 0x68
			fmt.Printf("&: %08b\n", a)
			a = a | 56
			fmt.Printf("|: %08b\n", a)
			a = a ^ 17
			fmt.Printf("^: %08b\n", a)
			fmt.Printf("A: %08b  %02x %0x\n", a, a, a)
		}

		// fmt.Printf("\n")
	}
exitDebugConsole:
	fmt.Printf("\n")
}
