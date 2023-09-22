package IO

type Memory interface {
	Set(addr uint16, value byte) error
	SetWord(addr uint16, value uint16) error
	Get(addr uint16) (byte, error)
	GetWord(addr uint16) (uint16, error)
}

type IO struct {
	Devices []Memory
}

func New(devices []Memory) *IO {
	return &IO{
		Devices: devices,
	}
}

func (io *IO) Set(addr uint16, value byte) error {
	return io.Devices[0].Set(addr, value)
}

func (io *IO) SetWord(addr uint16, value uint16) error {
	return io.Devices[0].SetWord(addr, value)
}

func (io *IO) Get(addr uint16) (byte, error) {
	return io.Devices[0].Get(addr)
}

func (io *IO) GetWord(addr uint16) (uint16, error) {
	return io.Devices[0].GetWord(addr)
}
