package Keyboard

import "fmt"

type Keyboard struct {
	buffer []byte
	offset uint16
	mode   uint8
}

func New(offset uint16) *Keyboard {
	return &Keyboard{
		buffer: make([]byte, 0),
		offset: offset,
		mode:   0x00,
	}
}

func (k *Keyboard) AppendKeys(keys []rune) {
	for _, key := range keys {
		k.buffer = append(k.buffer, byte(0x80|key))
	}
	// fmt.Printf("Keyboard: %v\n", k.buffer)
}

// IO.Memory Interface
func (k *Keyboard) Size() uint16 {
	return 1
}

func (k *Keyboard) Set(addr uint16, value byte) error {
	if addr == k.offset+1 {
		fmt.Printf("keyboard setting control registers: $%04x $%02x\n", addr, value)
		k.mode = value
		return nil
	}

	if k.mode == 0x00 {
		fmt.Printf("keyboard in configuration mode, ignoring: $%04x $%02x\n", addr, value)
		return nil
	}

	return fmt.Errorf("Keyboard: not yet implemented")
}
func (k *Keyboard) SetWord(addr uint16, value uint16) error {
	return fmt.Errorf("Keyboard: not yet implemented")
}

func (k *Keyboard) Get(addr uint16) (byte, error) {
	// always return ready?
	if addr == k.offset+1 {
		return 0x80, nil
	}

	if len(k.buffer) > 1 {
		key := k.buffer[0]
		k.buffer = k.buffer[1:]
		fmt.Printf("Keyboard: %v\n", k.buffer)
		// set bit 7 to 1
		return 0x80 | key, nil
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
