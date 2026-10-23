package CPU

import "testing"

func TestBRKAndRTI(t *testing.T) {
	for _, sp := range []byte{0xff, 0} {
		c, m := testCPU(byte(BRK), 0xea, byte(NOP))
		c.SP = sp
		c.Status = Reserved | C | D | V
		m.SetWord(0xfffe, 0x3456)
		m.Set(0x3456, byte(RTI))
		step(t, c, m)
		if c.PC != 0x3456 || c.SP != sp-3 || c.Status != (Reserved|C|D|V|I) {
			t.Fatalf("BRK CPU=%+v", c)
		}
		for i, want := range []byte{0x02, 0x02, Reserved | B | C | D | V} {
			addr := 0x100 + uint16(sp-byte(i))
			if got := m.Bytes[addr]; got != want {
				t.Fatalf("stack[%04x]=%02x want %02x", addr, got, want)
			}
		}
		step(t, c, m)
		if c.PC != 0x202 || c.SP != sp || c.Status != (Reserved|C|D|V) {
			t.Fatalf("RTI CPU=%+v", c)
		}
	}
}
func TestBRKWrapsPC(t *testing.T) {
	c, m := testCPU()
	c.PC = 0xffff
	m.Bytes[0xffff] = byte(BRK) // also vector high byte
	m.Bytes[0xfffe] = 0x80
	step(t, c, m)
	if c.PC != 0x80 || m.Bytes[0x1ff] != 0 || m.Bytes[0x1fe] != 1 {
		t.Fatalf("BRK wrap CPU=%+v", c)
	}
}
func TestStackStatusBits(t *testing.T) {
	c, m := testCPU(byte(PHP), byte(PLP))
	c.Status = C
	step(t, c, m)
	if m.Bytes[0x1ff] != C|B|Reserved {
		t.Fatalf("PHP pushed %02x", m.Bytes[0x1ff])
	}
	step(t, c, m)
	if c.Status != C|Reserved {
		t.Fatalf("PLP status=%02x", c.Status)
	}
}
