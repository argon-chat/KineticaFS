package guid

import (
	"testing"
)

func TestGenerateRandomEntropy(t *testing.T) {
	entropy1, err := GenerateRandomEntropy()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entropy2, err := GenerateRandomEntropy()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entropy1 == entropy2 {
		t.Errorf("expected different entropy values, got identical: %d", entropy1)
	}
}

func TestGuid_Pack(t *testing.T) {
	g := NewGuid(
		0x12345678,
		0x12,
		0xFFFF,
		0x0123456789ABCDEF,
		0x0A,
	)
	bytes := g.Pack()

	if len(bytes) != 16 {
		t.Fatalf("expected 16 bytes, got %d", len(bytes))
	}
	if bytes[0] != 0x12 || bytes[1] != 0x34 || bytes[2] != 0x56 || bytes[3] != 0x78 {
		t.Errorf("epochTs bytes incorrect: %x", bytes[0:4])
	}
	if bytes[4] != 0x12 {
		t.Errorf("regionId byte incorrect: %x", bytes[4])
	}
	if bytes[5] != 0xFF || bytes[6] != 0xFF {
		t.Errorf("bucketCode bytes incorrect: %x", bytes[5:7])
	}
	expectedEntropy := []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF}
	for i := 0; i < 8; i++ {
		if bytes[7+i] != expectedEntropy[i] {
			t.Errorf("randomEntropy byte %d incorrect: %x", i, bytes[7:15])
			break
		}
	}
	if bytes[15]&0x0F != 0x0A {
		t.Errorf("reservedFlags incorrect in last byte: %x", bytes[15])
	}
	var checksum byte
	for i := 0; i < 15; i++ {
		checksum ^= bytes[i]
	}
	checksum &= 0x0F
	if (bytes[15]>>4)&0x0F != checksum {
		t.Errorf("checksum incorrect in last byte: %x, expected: %x", (bytes[15]>>4)&0x0F, checksum)
	}
}
