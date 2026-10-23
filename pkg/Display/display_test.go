package Display

import "testing"

func TestScrolling(t *testing.T) {
	for _, tt := range []struct{ name, input, want string }{
		{"wrap", "ABCD", "CD\n  \n"},
		{"repeated wrap", "ABCDEFGHIJ", "IJ\n  \n"},
		{"carriage return", "A\rB\rC", "B \nC \n"},
		{"CR at last column", "A\rB\rC\r", "C \n  \n"},
		{"ignore control", "A\x01B", "AB\n  \n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := New(0xd012, 2, 2)
			d.Set(0xd013, 0xa7)
			for _, b := range []byte(tt.input) {
				if err := d.Set(0xd012, b|0x80); err != nil {
					t.Fatal(err)
				}
				_ = d.All(true)
			}
			if got := d.All(false); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
			if len(d.buffer) != 4 {
				t.Fatalf("buffer length %d", len(d.buffer))
			}
		})
	}
}
