package monitor

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestOriginalWozmon(t *testing.T) {
	image := Wozmon()
	if len(image) != 256 {
		t.Fatalf("ROM length=%d", len(image))
	}
	listing, err := os.ReadFile("../../docs/roms/wozmon.txt")
	if err != nil {
		t.Fatal(err)
	}
	var expected []byte
	for _, line := range strings.Split(strings.TrimSpace(string(listing)), "\n") {
		fields := strings.SplitN(line, ":", 2)
		if len(fields) != 2 {
			t.Fatalf("invalid listing line %q", line)
		}
		var address uint16
		if _, err := fmt.Sscanf(fields[0], "%x", &address); err != nil {
			t.Fatal(err)
		}
		if int(address) != int(Address)+len(expected) {
			t.Fatalf("listing address=%04x", address)
		}
		data, err := hex.DecodeString(strings.Join(strings.Fields(fields[1]), ""))
		if err != nil {
			t.Fatal(err)
		}
		expected = append(expected, data...)
	}
	if !bytes.Equal(image, expected) {
		t.Fatal("ROM differs from original hex listing")
	}
	if got := fmt.Sprintf("%x", sha256.Sum256(image)); got != "e5af0d1c4057bd8e0ef5cb069c208ff7cc0984a7dff53b12c5cf119de8cb5c25" {
		t.Fatalf("ROM checksum=%s", got)
	}
	image[0] = 0
	if Wozmon()[0] != 0xd8 {
		t.Fatal("caller changed embedded ROM")
	}
	if image[0xfc] != 0 || image[0xfd] != 0xff {
		t.Fatal("reset vector is not FF00")
	}
}
