package Display

import "fmt"

type Display struct {
	buffer []byte
	index  int
	size   int
	offset uint16
}

func New(offset uint16, cols int, rows int) *Display {
	return &Display{
		offset: offset,
		index:  0,
		size:   cols * rows,
		buffer: make([]byte, cols*rows),
	}
}

func (d *Display) All() string {
	t := string(d.buffer)
	// fmt.Printf("Display:All: %v\n%v\n", len(t), d.buffer)
	return t
}

// func (d *Display) Write(addr uint16, value byte) error {
// 	if addr >= uint16(d.size) {
// 		return fmt.Errorf("Display: Out of range %02x", addr)
// 	}

// 	d.buffer[addr] = value
// 	return nil
// }

// func (d *Display) Read(addr uint16) (byte, error) {
// 	if addr >= uint16(d.size) {
// 		return 0x00, fmt.Errorf("Display: Out of range %02x", addr)
// 	}
// 	return d.buffer[addr], nil
// }

// IO.Memory Interface

func (d *Display) Size() uint16 {
	return 2
}

func (d *Display) Set(addr uint16, value byte) error {
	fmt.Printf("Display: Set: $%04x $%02x\n", addr, value)
	d.buffer[d.index] = value
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
	return d.buffer[addr-d.offset], nil
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
