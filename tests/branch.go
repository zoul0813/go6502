package main

import "fmt"

const (
	Bit0 = 0b00000001
	Bit1 = 0b00000010
	Bit2 = 0b00000100
	Bit3 = 0b00001000
	Bit4 = 0b00010000
	Bit5 = 0b00100000
	Bit6 = 0b01000000
	Bit7 = 0b10000000
)

func BitTest(bit uint8, value uint8) bool {
	b := value & bit
	return b != 0
}

func main() {
	fmt.Printf("Testing...\n")

	var pc uint16 = 0xFFF2
	pc++
	var rel uint8 = 0xFB

	j := ^rel + 1

	fmt.Printf(" PC: %04x, rel: %02x, j: %02x\n", pc, rel, j)
	fmt.Printf("Rel: %08b\n", rel)
	fmt.Printf("  J: %08b\n", j)

	if BitTest(Bit7, rel) {
		pc -= uint16(j)
	} else {
		pc += uint16(j)
	}
	pc++

	fmt.Printf("PC: %04x\n", pc)

}
