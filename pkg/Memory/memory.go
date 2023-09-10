package Memory

import (
	"errors"
	"fmt"
)

// type Addressable interface {
// 	Load() (uint16, error)
// 	Get(addr uint16) (uint16, error)
// 	GetWord(addr uint16) (uint16, error)
// 	Set(addr uint16, value byte) error
// 	SetWord(addr uint16, value uint16) error
// }

type Memory struct {
	Bytes    []byte
	Offset   uint16
	ReadOnly bool
	Next     uint16
}

func New(size uint32, offset uint16, readOnly bool) *Memory {
	fmt.Printf("Memory: %v at %04x\n", size, offset)
	return &Memory{
		Bytes:    make([]byte, size+1), // we add 1, cause 0x0000:0xFFFF is (0:65536)
		Offset:   offset,
		ReadOnly: readOnly,
		Next:     0,
	}
}

func (o *Memory) Load(bytes []byte) (uint16, error) {
	if len(bytes) > len(o.Bytes) {
		return 0, errors.New(fmt.Sprintf("%04x is too large for ROM with %04x", len(bytes), len(o.Bytes)))
	}

	size := copy(o.Bytes, bytes)

	return uint16(size), nil
}

func (o *Memory) Get(addr uint16) (byte, error) {
	a := addr - o.Offset
	if int(a) <= len(o.Bytes) {
		return o.Bytes[a], nil
	}
	return 0x00, errors.New(fmt.Sprintf("%04x is out of range of %04x", addr, len(o.Bytes)))
}

func (o *Memory) GetWord(addr uint16) (uint16, error) {
	a := addr - o.Offset
	if int(a+1) <= len(o.Bytes) {
		lo := o.Bytes[a]
		hi := o.Bytes[a+1]
		var word uint16 = (uint16(hi) << 8) + uint16(lo)
		return word, nil
	}
	return 0x00, errors.New(fmt.Sprintf("%04x is out of range of %04x", addr, len(o.Bytes)))
}

func (o *Memory) Set(addr uint16, value byte) error {
	if o.ReadOnly {
		return errors.New(fmt.Sprintf("Attempt to write to ROM at %04x", addr))
	}

	a := addr - o.Offset
	if int(a) > len(o.Bytes) {
		return errors.New(fmt.Sprintf("%04x is out of range of %04x", addr, len(o.Bytes)))
	}
	o.Bytes[a] = value
	return nil
}

func (o *Memory) SetWord(addr uint16, value uint16) error {
	if o.ReadOnly {
		return errors.New(fmt.Sprintf("Attempt to write to ROM at %04x", addr))
	}

	a := addr - o.Offset
	if int(a)+1 > len(o.Bytes) {
		return errors.New(fmt.Sprintf("%04x is out of range of %04x", addr, len(o.Bytes)))
	}

	hi := byte(value >> 8)
	lo := uint8(value)

	// fmt.Printf("SET: %04x: %04x (%04x) -> %02x = %02x", addr, value, o.Offset, a, lo)
	o.Bytes[a] = lo
	// fmt.Printf(" | %02x = %02x\n", a+1, hi)
	o.Bytes[a+1] = hi
	return nil
}

func (o *Memory) Goto(addr uint16) error {
	a := addr - o.Offset
	if int(a) > len(o.Bytes) {
		return errors.New(fmt.Sprintf("%04x is out of range of %04x", addr, len(o.Bytes)))
	}
	o.Next = addr

	return nil
}

func (o *Memory) next(offset byte) error {
	next := int(o.Next) + int(offset)
	if next < len(o.Bytes) {
		o.Next = uint16(next)
	} else {
		return errors.New(fmt.Sprintf("%04x is out of range of %04x", next, len(o.Bytes)))
	}
	return nil
}

func (o *Memory) Write(value byte) error {
	err := o.Set(o.Next, value)
	o.next(1)
	return err
}

func (o *Memory) WriteWord(value uint16) error {
	err := o.SetWord(o.Next, value)
	o.next(2)
	return err
}

func (o *Memory) Dump(addr uint16, size uint16) {
	a := addr - o.Offset
	if int(a+size) > len(o.Bytes) {
		fmt.Printf("%04x is out of range of %04x", addr, len(o.Bytes))
		return
	}

	fmt.Printf("Memory Dump (%07x:%07x)\n", addr, addr+size)
	fmt.Printf(Colorize(DebugColor, "%s", "------- 0001 0203 0405 0607 0809 0A0B 0C0D 0E0F\n"))
	fmt.Printf(Colorize(DebugColor, "%s", "------- ---- ---- ---- ---- ---- ---- ---- ----\n"))
	// fmt.Printf("0000000 2aa5 3818 0000 0000 0000 0000 0000 0000")
	var lcv int = 0
	var i uint16 = addr
	var end = addr + size
	for i < end {
		fmt.Printf(Colorize(AddrColor, "%07x ", i))
		for w := 0; w < 8; w++ {
			// fmt.Print("ww")
			w1, _ := o.Get(i)
			i++
			w2, _ := o.Get(i)
			if i < end {
				i++
			}

			// fmt.Printf("%02x%02x ", w1, w2)
			w1c := HiColor
			if w1 == 0 {
				w1c = EmptyColor
			}
			w2c := LoColor
			if w2 == 0 {
				w2c = EmptyColor
			}
			fmt.Printf("%s%s ", Colorize(w1c, "%02x", w1), Colorize(w2c, "%02x", w2))
			lcv++
		}

		fmt.Printf("\n")
	}
	fmt.Printf(Colorize(DebugColor, "%s", "------- ---- ---- ---- ---- ---- ---- ---- ----\n"))
	fmt.Printf("%v bytes, %v loops, (%07x:%07x):%07x\n\n", size+1, lcv, addr, end, i)
}

const (
	EmptyColor = "\033[1;90m%s\033[0m"
	HiColor    = "\033[1;31m%s\033[0m"
	LoColor    = "\033[1;91m%s\033[0m"
	AddrColor  = "\033[1;32m%s\033[0m"
	DebugColor = "\033[0;37m%s\033[0m"
)

func Colorize(color string, format string, v any) string {
	f := fmt.Sprintf(format, v)
	return fmt.Sprintf(color, f)
}
