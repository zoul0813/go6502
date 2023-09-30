package main

import "fmt"

func main() {
	fmt.Printf("Testing...\n")

	var buffer []byte
	buffer = []byte{1}
	b := buffer[0]
	buffer = buffer[1:]
	fmt.Printf("b: %v, buffer: %v\n", b, buffer)

	var a uint16 = 0x4001
	var m uint16 = 0xEFFF
	fmt.Printf("A: %008b\n", a)
	fmt.Printf("M: %008b\n", m)
	// a = (a << 1) + 1
	// fmt.Printf("+: %08b  %02x %0x\n", a, a, a)
	// a = a << 1
	// fmt.Printf("-: %08b  %02x %0x\n", a, a, a)

	// fmt.Printf("A: %08b  %02x %0x\n", a, a, a)
	// a = a & 0x37
	// fmt.Printf("&: %08b  %02x %0x\n", a, a, a)
	// a = a | 0x23
	// fmt.Printf("|: %08b  %02x %0x\n", a, a, a)
	// a = a ^ 0x9d
	// fmt.Printf("^: %08b  %02x %0x\n", a, a, a)

	// final
	// fmt.Printf("A: %08b  %02x %0x\n", a, a, a)
}
