package guid

import (
	"crypto/rand"
	"encoding/binary"
)

type Guid struct {
	epochTs       uint32
	regionId      byte
	bucketCode    uint16
	randomEntropy uint64
	reservedFlags byte
}

func NewGuid(epochTs uint32, regionId byte, bucketCode uint16, randomEntropy uint64, reservedFlags byte) *Guid {
	return &Guid{
		epochTs:       epochTs,
		regionId:      regionId,
		bucketCode:    bucketCode,
		randomEntropy: randomEntropy,
		reservedFlags: reservedFlags,
	}
}

func GenerateRandomEntropy() (uint64, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(b[:]), nil
}

func (g *Guid) Pack() (bytes [16]byte) {
	binary.BigEndian.PutUint32(bytes[0:4], g.epochTs)
	bytes[4] = g.regionId
	binary.BigEndian.PutUint16(bytes[5:7], g.bucketCode)
	binary.BigEndian.PutUint64(bytes[7:15], g.randomEntropy)

	lastByte := g.reservedFlags & 0x0F

	var checksum byte
	for i := 0; i < 15; i++ {
		checksum ^= bytes[i]
	}
	checksum &= 0x0F

	bytes[15] = (checksum << 4) | lastByte

	return bytes
}
