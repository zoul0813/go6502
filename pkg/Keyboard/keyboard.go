package Keyboard

import "fmt"

type Keyboard struct {
	buffer []byte
}

const (
	bufferSize = 16
)

func New(addr uint16) *Keyboard {
	return &Keyboard{
		buffer: make([]byte, bufferSize),
	}
}

func (k *Keyboard) Set(addr uint16, value byte) error {
	k.buffer[0] = value
	return fmt.Errorf("Keyboard: not yet implemented")
}
func (k *Keyboard) SetWord(addr uint16, value uint16) error {
	k.buffer[0] = byte(value >> 8)
	k.buffer[1] = byte(value & 0xFF)
	return fmt.Errorf("Keyboard: not yet implemented")
}

func (k *Keyboard) Get(addr uint16) (byte, error) {
	return 0xEB, fmt.Errorf("Keyboard: read-only, invalid memory access")
}
func (k *Keyboard) GetWord(addr uint16) (uint16, error) {
	return 0xEAEA, fmt.Errorf("Keyboard: read-only, invalid memory access")
}
