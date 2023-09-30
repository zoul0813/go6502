package Memory

import (
	"fmt"
	"log"
)

type Memory struct {
	Bytes    []byte
	Offset   uint16
	ReadOnly bool
	Next     uint16
}

func New(size uint32, offset uint16, readOnly bool) *Memory {
	fmt.Printf("Creating %v ($%04x) bytes at %v ($%04x) ", size, size, offset, offset)
	o := &Memory{
		Bytes:    make([]byte, size+1), // we add 1, cause 0x0000:0xFFFF is (0:65536)
		Offset:   offset,
		ReadOnly: readOnly,
		Next:     0,
	}
	fmt.Printf(" -> %v (%04x) bytes\n", len(o.Bytes), len(o.Bytes))
	return o
}

// IO.Memory Interface
func (o *Memory) Size() uint16 {
	return uint16(len(o.Bytes) - 1)
}
func (o *Memory) Get(addr uint16) (byte, error) {
	a := addr - o.Offset
	if int(a) <= len(o.Bytes) {
		return o.Bytes[a], nil
	}
	return 0x00, fmt.Errorf("%04x is out of range of %04x", addr, len(o.Bytes))
}

func (o *Memory) Set(addr uint16, value byte) error {
	// o.mutex.Lock()
	// defer o.mutex.Unlock()

	if o.ReadOnly {
		fmt.Printf("attempt to write to ROM at %04x\n", addr)
		return fmt.Errorf("attempt to write to ROM at %04x", addr)
	}

	a := addr - o.Offset
	// fmt.Printf("SET: %04x [%04x]: %02x (%v, %04x)\n", addr, a, value, len(o.Bytes), len(o.Bytes))
	if int(a) >= len(o.Bytes) {
		return fmt.Errorf("%04x is out of range of %04x, trying to set %02x", addr, len(o.Bytes), value)
	}
	o.Bytes[a] = value
	return nil
}

func (o *Memory) GetWord(addr uint16) (uint16, error) {
	a := addr - o.Offset
	if int(a+1) <= len(o.Bytes) {
		lo := o.Bytes[a]
		hi := o.Bytes[a+1]
		var word uint16 = (uint16(hi) << 8) + uint16(lo)
		return word, nil
	}
	return 0x00, fmt.Errorf("%04x is out of range of %04x", addr, len(o.Bytes))
}

func (o *Memory) SetWord(addr uint16, value uint16) error {
	// o.mutex.Lock()
	// defer o.mutex.Unlock()

	if o.ReadOnly {
		return fmt.Errorf("attempt to write to ROM at %04x", addr)
	}

	a := addr - o.Offset
	if int(a)+1 > len(o.Bytes) {
		return fmt.Errorf("%04x is out of range of %04x", addr, len(o.Bytes))
	}

	hi := byte(value >> 8)
	lo := uint8(value)

	// fmt.Printf("SET: %04x: %04x (%04x) -> %02x = %02x", addr, value, o.Offset, a, lo)
	o.Bytes[a] = lo
	// fmt.Printf(" | %02x = %02x\n", a+1, hi)
	o.Bytes[a+1] = hi
	return nil
}

func (o *Memory) Load(bytes []byte) (uint16, error) {
	ro := o.ReadOnly
	o.ReadOnly = false
	for i, b := range bytes {
		// fmt.Printf("%v: %04x: %02x (%v, %04x)\n", i, o.Offset+uint16(i), b, len(o.Bytes), len(bytes))
		err := o.Set((o.Offset + uint16(i)), b)
		if err != nil {
			fmt.Printf("Load Error: %v\n", err)
			log.Fatal(err)
		}
	}
	o.ReadOnly = ro
	fmt.Printf("Loaded %v bytes into %v\n", len(bytes), "memory")
	return uint16(len(bytes)), nil
}

// Memory.Memory public

func (o *Memory) Goto(addr uint16) error {
	a := addr - o.Offset
	if int(a) > len(o.Bytes) {
		return fmt.Errorf("%04x is out of range of %04x", addr, len(o.Bytes))
	}
	o.Next = addr

	return nil
}

func (o *Memory) Write(value byte) error {
	if err := o.Set(o.Next, value); err != nil {
		return err
	}
	return o.next(1)
}

func (o *Memory) WriteWord(value uint16) error {
	if err := o.SetWord(o.Next, value); err != nil {
		return err
	}
	return o.next(2)
}

// Memory.Memory private

func (o *Memory) next(offset byte) error {
	next := int(o.Next) + int(offset)
	if next < len(o.Bytes) {
		o.Next = uint16(next)
	} else {
		return fmt.Errorf("%04x is out of range of %04x", next, len(o.Bytes))
	}
	return nil
}
