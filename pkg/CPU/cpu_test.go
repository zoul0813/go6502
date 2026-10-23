package CPU

import (
	"github.com/zoul0813/go6502/pkg/Memory"
	"testing"
)

func testCPU(code ...byte) (*CPU, *Memory.Memory) {
	m := &Memory.Memory{Bytes: make([]byte, 65536)}
	copy(m.Bytes[0x200:], code)
	return New(0x200, 0xff, 0, 0, 0, Reserved, false, false), m
}
func step(t *testing.T, c *CPU, m *Memory.Memory) {
	t.Helper()
	if halted, err := c.Step(m); halted || err != nil {
		t.Fatalf("halted=%v err=%v", halted, err)
	}
}
func TestStackRegisterTransfers(t *testing.T) {
	for _, value := range []byte{0, 0x42, 0x80, 0xff} {
		c, m := testCPU(byte(TXS))
		c.X = value
		c.Status = 0xff
		m.Bytes[0x1ff] = 0x55
		step(t, c, m)
		if c.SP != value || c.Status != 0xff || m.Bytes[0x1ff] != 0x55 {
			t.Fatalf("TXS value=%02x CPU=%+v", value, c)
		}
		c, m = testCPU(byte(TSX))
		c.SP = value
		c.Status = Reserved | C | V | Z | N
		m.Bytes[0x100+int(value)] = 0x55
		step(t, c, m)
		want := byte(Reserved | C | V)
		if value == 0 {
			want |= Z
		}
		want |= value & N
		if c.X != value || c.SP != value || c.Status != want {
			t.Fatalf("TSX value=%02x CPU=%+v want status=%02x", value, c, want)
		}
	}
}
