package Display

import (
	"fmt"
)

type Display struct {
	buffer []byte
	size   int
	offset uint16
	cols   int
	rows   int
	mode   uint8
	col    int
	row    int
	blink  int
}

func New(offset uint16, cols int, rows int) *Display {
	size := cols * rows
	return &Display{
		offset: offset,
		cols:   cols,
		rows:   rows,
		size:   size,
		buffer: make([]byte, size),
		mode:   0x00,
		blink:  0,
	}
}

func (d *Display) All(blink bool) string {
	t := ""

	for r := 0; r < d.rows; r++ {
		for c := 0; c < d.cols; c++ {
			b := d.buffer[r*d.cols+c]
			if b == 13 {
				b = 0x20 // replace CR with space
			}
			if r == d.row && c == d.col && blink {
				t += "@"
			} else {
				t += string(b)
			}
		}
		t += "\n"
	}

	// for i := 0; i < d.size; i++ {
	// 	b := d.buffer[i]
	// 	// convert CR to LF
	// 	if b == 13 {
	// 		b = 10
	// 	}
	// 	t += string(b)
	// }

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
		d.col = 0
		d.row = 0
		for i, _ := range d.buffer {
			d.buffer[i] = 0x20 // Space
		}
		return nil
	}

	if d.mode == 0x00 {
		fmt.Printf("display in configuration mode, ignoring: $%04x $%02x\n", addr, value)
		return nil
	}

	// strip bit 7
	c := value & 0b01111111 // $7F

	fmt.Printf("display: set: $%04x $%02x %v '%v' @ (%v, %v)\n", addr, value, c, string(c), d.col, d.row)
	d.buffer[d.row*d.cols+d.col] = c
	d.col++
	fmt.Printf("display: col: %v\n", d.col)
	if d.col >= d.cols {
		d.row++
		d.col = 0
		fmt.Printf("display: row: %v\n", d.row)
		if d.row >= d.rows {
			d.row--
			d.buffer = d.buffer[d.cols:]
			fmt.Printf("display: col: %v, row: %v\n", d.col, d.row)
		}
	}

	if c == 13 {
		d.row++
		d.col = 0
		if d.row >= d.rows {
			d.row--
			d.buffer = d.buffer[d.cols:d.size]
			for i := 0; i < d.cols; i++ {
				// d.buffer[d.row*d.cols+i] = 0x20 // Space
				d.buffer = append(d.buffer, 0x20)
			}
			fmt.Printf("display: col: %v, row: %v\n", d.col, d.row)
		}
		fmt.Printf("display: col: %v, row: %v\n", d.col, d.row)
	}

	// d.buffer[d.index] = c
	// d.index++
	// if d.index >= d.size {
	// 	d.index = 0
	// }
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
