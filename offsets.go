package gopal

import (
	"fmt"
	"io"

	"github.com/akaspin/chd"
	"github.com/remerge/mph"
)

const (
	V1Magic = 0x19820304
	V2Magic = 0x19820305
)

// Generic offsets mapper
type Offsets interface {
	io.Reader
	io.WriterTo

	Get(key []byte) []byte
	GetRandomKey() []byte
	GetRandomValue() []byte
}

func GetOffsets(magic uint64) (o Offsets, err error) {
	switch magic {
	case V1Magic:
		o = &v1offsets{}
	case V2Magic:
		o = chd.NewMap()
	default:
		err = fmt.Errorf(`invalid magic`)
	}
	return
}

type v1offsets struct {
	*mph.CHD
}

func (o *v1offsets) Read(p []byte) (n int, err error) {
	o.CHD, err = mph.Mmap(p)
	return
}

func (o *v1offsets) WriteTo(w io.Writer) (n int64, err error) {
	err = o.CHD.Write(w)
	return
}

