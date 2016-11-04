package gopal

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"unsafe"

	"github.com/remerge/mph"
)

type Builder struct {
	cdh             *mph.CHDBuilder
	numFields       int
	emptyHeaderSize int
	buf             bytes.Buffer
	emptyHeader     []byte
	pos             int
	fields          map[string]int
}

func NewBuilder(fields []string) *Builder {
	numFields := len(fields)
	m := make(map[string]int)
	for idx, f := range fields {
		m[f] = idx
	}

	emptyHeaderSize := numFields * 4
	return &Builder{
		cdh:             mph.Builder(),
		numFields:       numFields,
		fields:          m,
		emptyHeaderSize: emptyHeaderSize,
		emptyHeader:     make([]byte, emptyHeaderSize),
	}
}

// does the id lookup for us
func (b *Builder) AddRow(values []string) {
	// TODO : error checks
	i := b.fields["id"]
	id := values[i]
	b.Add(id, values)
}

func (b *Builder) Add(id string, values []string) {
	pos := b.pos
	offsetBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(offsetBuf, uint64(pos))
	b.cdh.Add([]byte(id), offsetBuf)
	// write placeholder
	n, _ := b.buf.Write(b.emptyHeader)
	pos += n
	// reserve space for offsets
	offset := 0
	for idx, v := range values {
		// update the offset
		n, _ := b.buf.Write([]byte(v))
		pos += n
		offset += n
		// store the offset
		op := b.buf.Bytes()[b.pos+idx*4 : b.pos+4+idx*4]
		binary.LittleEndian.PutUint32(op, uint32(offset))
	}
	b.pos = pos
}

func (b *Builder) BuildTo(w io.Writer) {
	h := &PalHeader{Magic: 0x19820304, HeadSize: uint64(unsafe.Sizeof(PalHeader{}))}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(b.fields)
	h.MapSize = uint64(buf.Len())
	cdhb, err := b.cdh.Build()
	if err != nil {
		fmt.Println(err)
	}
	cdhb.Write(&buf)
	h.IdxSize = uint64(buf.Len()) - h.MapSize
	h.WriteTo(w)
	buf.WriteTo(w)
	b.buf.WriteTo(w)
}

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
