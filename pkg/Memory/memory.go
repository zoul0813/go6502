package Memory

import "fmt"

type Memory struct {
	Bytes    []byte
	Offset   uint16
	ReadOnly bool
	Next     uint16
}

// New allocates addresses offset through offset+size, inclusive.
// Size is the last relative address, as required by IO.Memory.
func New(size uint32, offset uint16, readOnly bool) *Memory {
	return &Memory{Bytes: make([]byte, int(size)+1), Offset: offset, ReadOnly: readOnly, Next: offset}
}

func (o *Memory) Size() uint16 { return uint16(len(o.Bytes) - 1) }

func (o *Memory) index(addr uint16, width int) (int, error) {
	i := int(addr) - int(o.Offset)
	if i < 0 || i+width > len(o.Bytes) || int(addr)+width > 65536 {
		return 0, fmt.Errorf("%04x: %d-byte access outside memory at %04x (%d bytes)", addr, width, o.Offset, len(o.Bytes))
	}
	return i, nil
}
func (o *Memory) Get(addr uint16) (byte, error) {
	i, err := o.index(addr, 1)
	if err != nil {
		return 0, err
	}
	return o.Bytes[i], nil
}
func (o *Memory) Set(addr uint16, value byte) error {
	if o.ReadOnly {
		return fmt.Errorf("attempt to write to ROM at %04x", addr)
	}
	i, err := o.index(addr, 1)
	if err != nil {
		return err
	}
	o.Bytes[i] = value
	return nil
}
func (o *Memory) GetWord(addr uint16) (uint16, error) {
	i, err := o.index(addr, 2)
	if err != nil {
		return 0, err
	}
	return uint16(o.Bytes[i]) | uint16(o.Bytes[i+1])<<8, nil
}
func (o *Memory) SetWord(addr uint16, value uint16) error {
	if o.ReadOnly {
		return fmt.Errorf("attempt to write to ROM at %04x", addr)
	}
	i, err := o.index(addr, 2)
	if err != nil {
		return err
	}
	o.Bytes[i] = byte(value)
	o.Bytes[i+1] = byte(value >> 8)
	return nil
}

// Load initializes memory (including ROM) without changing its write protection.
// Validate first so an oversized image cannot leave a partially loaded device.
func (o *Memory) Load(data []byte) (uint16, error) {
	if _, err := o.index(o.Offset, len(data)); err != nil {
		return 0, err
	}
	copy(o.Bytes, data)
	return uint16(len(data)), nil
}
func (o *Memory) Goto(addr uint16) error {
	if _, err := o.index(addr, 1); err != nil {
		return err
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
func (o *Memory) next(offset byte) error {
	next := int(o.Next) + int(offset)
	if next > 65535 {
		return fmt.Errorf("memory cursor exceeds address space")
	}
	if _, err := o.index(uint16(next), 1); err != nil {
		return err
	}
	o.Next = uint16(next)
	return nil
}
