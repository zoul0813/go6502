package Display

import (
	"fmt"
)

type Display struct {
	buffer []byte
	index  int
	size   int
	offset uint16
	cols   int
	rows   int
	mode   uint8
}

func New(offset uint16, cols int, rows int) *Display {
	return &Display{
		offset: offset,
		index:  0,
		cols:   cols,
		rows:   rows,
		size:   cols * rows,
		buffer: make([]byte, cols*rows),
		mode:   0x00,
	}
}

func (d *Display) All() string {
	t := ""

	for i := 0; i < d.size; i++ {
		b := d.buffer[i]
		// convert CR to LF
		if b == 13 {
			b = 10
		}
		t += string(b)
	}

	return t
}

// IO.Memory Interface

func (d *Display) Size() uint16 {
	return 1
}

func (d *Display) Set(addr uint16, value byte) error {
	if addr == d.offset+1 {
		fmt.Printf("display setting control registers: $%04x $%02x\n", addr, value)
		d.mode = value
		return nil
	}

	if d.mode == 0x00 {
		fmt.Printf("display in configuration mode, ignoring: $%04x $%02x\n", addr, value)
		return nil
	}

	c := value & 0b01111111

	fmt.Printf("display: set: $%04x $%02x %v '%v'\n", addr, value, c, string(c))
	// strip bit 7
	d.buffer[d.index] = c
	d.index++
	if d.index >= d.size {
		d.index = 0
	}
	return nil // never give up, never surrender
}

func (d *Display) SetWord(addr uint16, value uint16) error {
	a := addr - d.offset
	if a+1 >= uint16(d.size) {
		return fmt.Errorf("Display: invalid offset: %04x", addr)
	}
	d.buffer[a] = uint8(value >> 8)
	d.buffer[a+1] = uint8(value & 0xFF)
	return nil
}

func (d *Display) Get(addr uint16) (byte, error) {
	if addr == d.offset+1 {
		return d.mode, nil
	}

	return 0x00, nil
}

func (d *Display) GetWord(addr uint16) (uint16, error) {
	a := addr - d.offset
	if a+1 >= uint16(d.size) {
		return 0x00, fmt.Errorf("Display: invalid offset: %04x", addr)
	}
	hi := d.buffer[a]
	lo := d.buffer[a+1]
	b := (uint16(hi) << 8) + uint16(lo)
	return b, nil
}

func (d *Display) Load(bytes []byte) (uint16, error) {
	return 0x00, fmt.Errorf("not implemented: %v", len(bytes))
}
