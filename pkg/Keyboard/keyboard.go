package Keyboard

import "fmt"

type Keyboard struct {
	buffer []byte
}

func New(addr uint16) *Keyboard {
	return &Keyboard{
		buffer: make([]byte, 0),
	}
}

func (k *Keyboard) AppendKeys(keys []rune) {
	for _, key := range keys {
		k.buffer = append(k.buffer, byte(key))
	}
	// fmt.Printf("Keyboard: %v\n", k.buffer)
}

// IO.Memory Interface
func (k *Keyboard) Size() uint16 {
	return 2
}

func (k *Keyboard) Set(addr uint16, value byte) error {
	return fmt.Errorf("Keyboard: not yet implemented")
}
func (k *Keyboard) SetWord(addr uint16, value uint16) error {
	return fmt.Errorf("Keyboard: not yet implemented")
}

func (k *Keyboard) Get(addr uint16) (byte, error) {
	if len(k.buffer) > 1 {
		key := k.buffer[0]
		k.buffer = k.buffer[1:]
		// fmt.Printf("Keyboard: %v\n", k.buffer)
		return key, nil
	}
	return 0x00, nil

	// return byte(key), fmt.Errorf("Keyboard: read-only, invalid memory access")
}
func (k *Keyboard) GetWord(addr uint16) (uint16, error) {
	return 0xEAEA, fmt.Errorf("Keyboard: read-only, invalid memory access")
}

func (k *Keyboard) Load(bytes []byte) (uint16, error) {
	return 0x00, fmt.Errorf("not implemented: %v", len(bytes))
}
