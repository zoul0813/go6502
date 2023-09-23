package IO

import (
	"fmt"
	"sync"

	"github.com/zoul0813/go6502/pkg/Debug"
)

type Memory interface {
	Set(addr uint16, value byte) error
	SetWord(addr uint16, value uint16) error
	Get(addr uint16) (byte, error)
	GetWord(addr uint16) (uint16, error)
	Load(bytes []byte) (uint16, error)
}

type Device struct {
	Name   string
	Chip   Memory
	Size   uint16
	Offset uint16
}

type IO struct {
	Devices []*Device
	Size    uint16
	mutex   sync.Mutex
}

func New(devices []*Device, size uint16) *IO {
	return &IO{
		Devices: devices,
		Size:    size,
	}
}

func NewDevice(name string, chip Memory, size uint16, offset uint16) *Device {
	return &Device{
		Name:   name,
		Chip:   chip,
		Size:   size,
		Offset: offset,
	}
}

func (io *IO) getDevice(addr uint16) (*Device, error) {
	var device *Device
	for _, d := range io.Devices {
		// fmt.Printf("\tTesting: %v 0x%04x 0x%04x 0x%04x\n", d.Name, addr, d.Offset, d.Offset+d.Size)
		if addr >= d.Offset && addr <= d.Offset+d.Size {
			// fmt.Printf("\tFound: %v\n", d.Name)
			device = d
			break
		}
	}

	var err error = nil
	if device == nil {
		// fmt.Printf("\tDevice not found for address 0x%04x\n", addr)
		err = fmt.Errorf("Device not found for address 0x%04x", addr)
	} else {
		// fmt.Printf("\n\tDevice: %v\n", device.Name)
	}
	return device, err
}

func (io *IO) Set(addr uint16, value byte) error {
	var err error
	device, err := io.getDevice(addr)
	if err != nil {
		// fmt.Printf("Error: Device not found: %v\n", err)
		return err
	}
	// fmt.Printf("\tFound: %v, locking...", device.Name)
	chip := device.Chip
	io.mutex.Lock()
	defer io.mutex.Unlock()
	// fmt.Printf(" %02x -> %04x...", value, addr)
	err = chip.Set(addr, value)
	// device.Chip = &chip
	// fmt.Print(" Done\n")
	return err
}

func (io *IO) SetWord(addr uint16, value uint16) error {
	var err error
	device, err := io.getDevice(addr)
	if err != nil {
		return err
	}
	chip := device.Chip
	io.mutex.Lock()
	defer io.mutex.Unlock()
	err = chip.SetWord(addr, value)
	// device.Chip = &chip
	return err
}

func (io *IO) Get(addr uint16) (byte, error) {
	var err error
	device, err := io.getDevice(addr)
	if err != nil {
		return 0xEA, err
	}
	chip := device.Chip
	return chip.Get(addr)
}

func (io *IO) GetWord(addr uint16) (uint16, error) {
	var err error
	device, err := io.getDevice(addr)
	if err != nil {
		return 0xEAEA, err
	}
	chip := device.Chip
	return chip.GetWord(addr)
}

func (io *IO) Load(bytes []byte) (uint16, error) {
	return 0, fmt.Errorf("not implemented: %v", len(bytes))

	// if len(bytes) > int(io.Size-offset)+1 {
	// 	fmt.Printf("%04x is too large for ROM with %04x\n", len(bytes), io.Size)
	// 	return 0, fmt.Errorf("%04x is too large for ROM with %04x", len(bytes), io.Size)
	// }

	// device, err := io.getDevice(offset)
	// if err != nil {
	// 	fmt.Printf("Load Error: Device Not Found: %v\n", err)
	// 	return 0, err
	// }

	// size, err := device.Chip.Load(bytes)

	// for i, b := range bytes {
	// 	fmt.Printf("\t%04x: %02x\n", offset+uint16(i), b)
	// 	err := io.Set(offset+uint16(i), b)
	// 	if err != nil {
	// 		fmt.Printf("Load Error: %v\n", err)
	// 		log.Fatal(err)
	// 	}
	// }

	// return size, err
}

func (io *IO) Dump(addr uint16, size uint16) {
	a := addr
	if int(a+size) > int(io.Size) {
		fmt.Printf("%04x is out of range of %04x", addr, io.Size)
		return
	}

	fmt.Printf("Memory Dump (%04x:%04x)\n", addr, addr+size)
	fmt.Print(Debug.Colorize(Debug.DebugColor, "%s", "---- 0001 0203 0405 0607 0809 0A0B 0C0D 0E0F\n"))
	fmt.Print(Debug.Colorize(Debug.DebugColor, "%s", "---- ---- ---- ---- ---- ---- ---- ---- ----\n"))
	// fmt.Printf("0000000 2aa5 3818 0000 0000 0000 0000 0000 0000")
	var lcv int = 0
	var i uint16 = addr
	var end = addr + size
	for i < end {
		fmt.Print(Debug.Colorize(Debug.AddrColor, "%04x ", i))
		for w := 0; w < 8; w++ {
			// fmt.Printf("%02x %02x ", i, i+1)
			w1, _ := io.Get(i)
			i++
			w2, _ := io.Get(i)
			if i < end {
				i++
			}

			// fmt.Printf("%02x%02x ", w1, w2)
			w1c := Debug.HiColor
			if w1 == 0 {
				w1c = Debug.EmptyColor
			}
			w2c := Debug.LoColor
			if w2 == 0 {
				w2c = Debug.EmptyColor
			}
			fmt.Printf("%s%s ", Debug.Colorize(w1c, "%02x", w1), Debug.Colorize(w2c, "%02x", w2))
			lcv++
		}

		fmt.Printf("\n")
	}
	fmt.Print(Debug.Colorize(Debug.DebugColor, "%s", "------- ---- ---- ---- ---- ---- ---- ---- ----\n"))
	fmt.Printf("%v bytes, %v loops, (%07x:%07x):%07x\n\n", size+1, lcv, addr, end, i)
}
