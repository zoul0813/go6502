package CPU

import "testing"

func TestDecimalArithmetic(t *testing.T) {
	// Independent base-ten oracle for every valid packed-BCD operand.
	for a := 0; a < 100; a++ {
		for b := 0; b < 100; b++ {
			for carry := 0; carry < 2; carry++ {
				pack := func(n int) byte { return byte(n/10*16 + n%10) }
				for _, sub := range []bool{false, true} {
					c := CPU{A: pack(a), Status: Reserved | I | D | byte(carry)}
					want := a + b + carry
					wantCarry := want > 99
					if sub {
						want = a - b - (1 - carry)
						wantCarry = want >= 0
						c.SBC(pack(b))
					} else {
						c.ADC(pack(b))
					}
					want = (want + 100) % 100
					if c.A != pack(want) || (c.Status&C != 0) != wantCarry || c.Status&(Reserved|I|D) != (Reserved|I|D) {
						t.Fatalf("sub=%v a=%d b=%d C=%d got A=%02x P=%02x want %d C=%v", sub, a, b, carry, c.A, c.Status, want, wantCarry)
					}
				}
			}
		}
	}
}
func TestNMOSDecimalFlags(t *testing.T) {
	for _, tt := range []struct {
		a, b, carry, result, flags byte
		sub                        bool
	}{
		{0x50, 0x50, 0, 0x00, N | V | C, false}, // Decimal zero, binary nonzero.
		{0x79, 0x00, 1, 0x80, N | V, false},     // Low-digit carry feeds N/V.
		{0x99, 0x01, 0, 0x00, N | C, false},
		{0x00, 0x00, 1, 0x00, Z | C, true},
		{0x00, 0x01, 1, 0x99, N, true},
		{0x80, 0x01, 1, 0x79, V | C, true},
		{0x00, 0x00, 0, 0x99, N, true},
	} {
		c := CPU{A: tt.a, Status: D | tt.carry}
		if tt.sub {
			c.SBC(tt.b)
		} else {
			c.ADC(tt.b)
		}
		if c.A != tt.result || c.Status&(N|V|Z|C) != tt.flags {
			t.Fatalf("%+v: A=%02x P=%02x", tt, c.A, c.Status)
		}
	}
}
