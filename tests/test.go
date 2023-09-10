package main

import "fmt"

func main() {
	fmt.Printf("Testing...\n")

	var a uint8 = 0x2e
	fmt.Printf("A: %08b  %02x %0x\n", a, a, a)
	a = a | 0x41
	fmt.Printf("+: %08b  %02x %0x\n", a, a, a)
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
	fmt.Printf("A: %08b  %02x %0x\n", a, a, a)
}
