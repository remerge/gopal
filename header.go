package gopal

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

func (h *PalHeader) Validate() (err error) {
	switch h.Magic {
	case V1Magic, V2Magic:
	default:
		err = fmt.Errorf(`invalid magic %x (valid: %x, %x)`, h.Magic, V1Magic, V2Magic)
	}
	return
}

