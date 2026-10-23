package CPU

import "testing"

func TestBinaryArithmeticExhaustive(t *testing.T) {
	for _, subtract := range []bool{false, true} {
		for a := 0; a < 256; a++ {
			for b := 0; b < 256; b++ {
				for carry := 0; carry < 2; carry++ {
					c := CPU{A: byte(a), Status: Reserved | I | byte(carry)}
					unsigned, signed := a+b+carry, int(int8(a))+int(int8(b))+carry
					wantCarry := unsigned > 255
					if subtract {
						unsigned = a - b - (1 - carry)
						signed = int(int8(a)) - int(int8(b)) - (1 - carry)
						wantCarry = unsigned >= 0
						c.SBC(byte(b))
					} else {
						c.ADC(byte(b))
					}
					result := byte(unsigned)
					want := byte(Reserved|I) | result&N
					if result == 0 {
						want |= Z
					}
					if wantCarry {
						want |= C
					}
					if signed < -128 || signed > 127 {
						want |= V
					}
					if c.A != result || c.Status != want {
						t.Fatalf("subtract=%v a=%02x b=%02x carry=%d got A=%02x P=%02x want A=%02x P=%02x", subtract, a, b, carry, c.A, c.Status, result, want)
					}
				}
			}
		}
	}
}
