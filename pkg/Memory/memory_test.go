package Memory

import (
	"bytes"
	"testing"
)

func TestMemoryBounds(t *testing.T) {
	for _, base := range []uint16{0, 0xf000} {
		m := New(0xff, base, false)
		if err := m.SetWord(base+0xfe, 0x1234); err != nil {
			t.Fatal(err)
		}
		if got, err := m.GetWord(base + 0xfe); got != 0x1234 || err != nil {
			t.Fatalf("word=%04x err=%v", got, err)
		}
		if _, err := m.Get(base + 0x100); err == nil {
			t.Fatal("out-of-range Get succeeded")
		}
		if _, err := m.GetWord(base + 0xff); err == nil {
			t.Fatal("out-of-range GetWord succeeded")
		}
		before := append([]byte(nil), m.Bytes...)
		if err := m.SetWord(base+0xff, 0xffff); err == nil {
			t.Fatal("out-of-range SetWord succeeded")
		}
		if !bytes.Equal(before, m.Bytes) {
			t.Fatal("failed write changed memory")
		}
		if err := m.Goto(base + 0x100); err == nil {
			t.Fatal("invalid cursor accepted")
		}
		if base > 0 {
			if _, err := m.Get(base - 1); err == nil {
				t.Fatal("address below offset accepted")
			}
		}
	}
	m := New(0xffff, 0, false)
	if _, err := m.GetWord(0xffff); err == nil {
		t.Fatal("word past address space accepted")
	}
	if err := m.SetWord(0xffff, 0x1234); err == nil {
		t.Fatal("word write past address space accepted")
	}
}
func TestLoadAndCursor(t *testing.T) {
	m := New(3, 0xf000, true)
	if _, err := m.Load([]byte{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}
	if !m.ReadOnly || m.Set(0xf000, 0) == nil {
		t.Fatal("ROM protection lost")
	}
	if _, err := m.Load([]byte{9, 9, 9, 9, 9}); err == nil {
		t.Fatal("oversized load accepted")
	}
	if !bytes.Equal(m.Bytes, []byte{1, 2, 3, 4}) || !m.ReadOnly {
		t.Fatal("failed load mutated ROM")
	}
	m = New(3, 0x2000, false)
	if err := m.WriteWord(0x1234); err != nil {
		t.Fatal(err)
	}
	if m.Next != 0x2002 {
		t.Fatalf("cursor=%04x", m.Next)
	}
}
