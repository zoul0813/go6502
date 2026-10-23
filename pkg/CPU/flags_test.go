package CPU

import "testing"

func TestMemoryShiftFlags(t *testing.T) {
	for _, op := range []byte{ROL_ZP, ROL_ZPX, ROL_A, ROL_AX, ROR_ZP, ROR_ZPX, ROR_A, ROR_AX, ASL_ZP, ASL_ZPX, ASL_A, ASL_AX, LSR_ZP, LSR_ZPX, LSR_A, LSR_AX} {
		c, m := testCPU(op, 0x10, 0x00)
		for operand := 0; operand < 256; operand++ {
			for carry := 0; carry < 2; carry++ {
				c.PC = 0x200
				c.A = 0x55
				c.Status = Reserved | I | V | byte(carry)
				m.Bytes[0x10] = byte(operand)
				result := operand << 1
				wantCarry := operand&0x80 != 0
				switch op {
				case ROL_ZP, ROL_ZPX, ROL_A, ROL_AX:
					result |= carry
				case ROR_ZP, ROR_ZPX, ROR_A, ROR_AX:
					result = operand>>1 | carry<<7
					wantCarry = operand&1 != 0
				case LSR_ZP, LSR_ZPX, LSR_A, LSR_AX:
					result = operand >> 1
					wantCarry = operand&1 != 0
				}
				want := byte(Reserved|I|V) | byte(result)&N
				if byte(result) == 0 {
					want |= Z
				}
				if wantCarry {
					want |= C
				}
				step(t, c, m)
				if m.Bytes[0x10] != byte(result) || c.Status != want || c.A != 0x55 {
					t.Fatalf("op=%02x operand=%02x carry=%d got mem=%02x P=%02x want mem=%02x P=%02x", op, operand, carry, m.Bytes[0x10], c.Status, byte(result), want)
				}
			}
		}
	}
}
func TestBITFlags(t *testing.T) {
	for _, op := range []byte{BIT_ZP, BIT_A} {
		c, m := testCPU(op, 0x10, 0)
		for a := 0; a < 256; a++ {
			for value := 0; value < 256; value++ {
				c.PC = 0x200
				c.A = byte(a)
				c.Status = Reserved | C
				m.Bytes[0x10] = byte(value)
				step(t, c, m)
				want := byte(Reserved|C) | byte(value)&(N|V)
				if a&value == 0 {
					want |= Z
				}
				if c.Status != want {
					t.Fatalf("BIT A=%02x mem=%02x P=%02x want=%02x", a, value, c.Status, want)
				}
			}
		}
	}
}
func TestINCAndPLAFlags(t *testing.T) {
	c, m := testCPU(INC_ZP, 0x10)
	for value := 0; value < 256; value++ {
		c.PC = 0x200
		c.Status = Reserved | C | V
		m.Bytes[0x10] = byte(value)
		step(t, c, m)
		result := byte(value + 1)
		want := byte(Reserved|C|V) | result&N
		if result == 0 {
			want |= Z
		}
		if m.Bytes[0x10] != result || c.Status != want {
			t.Fatalf("INC operand=%02x mem=%02x P=%02x want=%02x", value, m.Bytes[0x10], c.Status, want)
		}
	}
	c, m = testCPU(PLA)
	for value := 0; value < 256; value++ {
		c.PC = 0x200
		c.SP = 0xfe
		c.Status = Reserved | C | V
		m.Bytes[0x1ff] = byte(value)
		step(t, c, m)
		want := byte(Reserved|C|V) | byte(value)&N
		if value == 0 {
			want |= Z
		}
		if c.A != byte(value) || c.SP != 0xff || c.Status != want {
			t.Fatalf("PLA value=%02x CPU=%+v want=%02x", value, c, want)
		}
	}
}
