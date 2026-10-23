package CPU

import "testing"

func TestIndirectZeroPageWrap(t *testing.T) {
	for _, op := range []byte{LDA_INX, LDA_INY} {
		c, m := testCPU(op, 0xff)
		c.Y = 1
		if op == LDA_INX {
			m.Bytes[0x201] = 0xfe
			c.X = 1
			c.Y = 0
		}
		m.Bytes[0xff] = 0x34
		m.Bytes[0] = 0x12
		m.Bytes[0x100] = 0x56
		m.Bytes[0x1234+uint16(c.Y)] = 0xab
		step(t, c, m)
		if c.A != 0xab || c.PC != 0x202 {
			t.Fatalf("opcode=%02x CPU=%+v", op, c)
		}
	}
	c, m := testCPU(LDA_INY, 0xff)
	c.Y = 2
	m.Bytes[0xff] = 0xff
	m.Bytes[0] = 0xff
	m.Bytes[1] = 0x42
	step(t, c, m)
	if c.Address != 1 || c.A != 0x42 {
		t.Fatalf("indexed address overflow CPU=%+v", c)
	}
}
func TestJMPIndirectPageWrap(t *testing.T) {
	for _, pointer := range []uint16{0x12fe, 0x12ff, 0xffff} {
		c, m := testCPU(JMP_IN, byte(pointer), byte(pointer>>8))
		m.Bytes[pointer] = 0x34
		m.Bytes[(pointer&0xff00)|uint16(byte(pointer+1))] = 0x12
		step(t, c, m)
		if c.PC != 0x1234 {
			t.Fatalf("pointer=%04x PC=%04x", pointer, c.PC)
		}
	}
}
func TestPageWordErrors(t *testing.T) {
	_, m := testCPU()
	m.Bytes = m.Bytes[:1]
	if _, err := readPageWord(m, 0); err == nil {
		t.Fatal("missing high byte accepted")
	}
	if _, err := readPageWord(m, 1); err == nil {
		t.Fatal("missing low byte accepted")
	}
}
