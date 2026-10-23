package IO

import (
	RAM "github.com/zoul0813/go6502/pkg/Memory"
	"testing"
)

func TestWordAcrossDevicesAndAddressWrap(t *testing.T) {
	a, b := RAM.New(0, 0xffff, false), RAM.New(0, 0, false)
	a.Set(0xffff, 0x34)
	b.Set(0, 0x12)
	bus := New([]*Device{NewDevice("last", a, 0xffff), NewDevice("first", b, 0)})
	if got, err := bus.GetWord(0xffff); got != 0x1234 || err != nil {
		t.Fatalf("got=%04x err=%v", got, err)
	}
	bus = New([]*Device{NewDevice("first", b, 0)})
	if _, err := bus.GetWord(0); err == nil {
		t.Fatal("missing high byte accepted")
	}
}
