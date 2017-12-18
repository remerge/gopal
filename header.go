package gopal

import (
	"bytes"
	"encoding/binary"
	"io"
)

type PalHeader struct {
	// TODO add a checksum
	Magic    uint64
	HeadSize uint64
	MapSize  uint64
	IdxSize  uint64
}

func (h *PalHeader) Read(b []byte) error {
	return binary.Read(bytes.NewBuffer(b), binary.LittleEndian, h)
}

func (h *PalHeader) WriteTo(w io.Writer) (int64, error) {
	return 0, binary.Write(w, binary.LittleEndian, h)
}

func (h *PalHeader) Valid() bool {
	return h.Magic == 0x19820304
}

