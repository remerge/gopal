package gopal

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

type Pal struct {
	fields  map[string]int
	idx     []byte
	data    []byte
	offsets Offsets
	mmap    []byte
}

func MMapPal(filename string) (*Pal, error) {
	p := &Pal{}

	fi, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, errors.New(fmt.Sprintf("%s is a directory, file needed", filename))
	}

	file, err := os.OpenFile(filename, os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	mmap, err := syscall.Mmap(int(file.Fd()), 0, int(fi.Size()), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	err = p.From(mmap)
	if err != nil {
		syscall.Munmap(mmap)
		return nil, err
	}
	p.mmap = mmap
	runtime.SetFinalizer(p, func(p *Pal) {
		p.Free()
	})
	return p, nil
}

func (p *Pal) Free() {
	if p == nil || p.mmap == nil {
		return
	}
	syscall.Munmap(p.mmap)
	p.data = nil
	p.mmap = nil
}

func (p *Pal) From(b []byte) error {
	var err error

	if len(b) < int(unsafe.Sizeof(PalHeader{})) {
		return fmt.Errorf("buffer seems to be too small len=%d", len(b))
	}

	h := &PalHeader{}
	h.Read(b)
	if err = h.Validate(); err != nil {
		return fmt.Errorf("header invalid (len=%d): %v", len(b), err)
	}

	// TODO  - checksum the file or at least check if the size matches

	hm := h.HeadSize + h.MapSize
	hmi := hm + h.IdxSize

	if hm > hmi || hmi > uint64(len(b)-1) || h.HeadSize > hm {
		return fmt.Errorf("header invalid h=%#v len(b)=%v", h, len(b))
	}

	p.idx = b[hm:hmi]
	p.data = b[hmi : len(b)-1]

	dec := gob.NewDecoder(bytes.NewBuffer(b[h.HeadSize:hm]))
	err = dec.Decode(&p.fields)
	if err != nil {
		return err
	}

	if p.offsets, err = GetOffsets(h.Magic); err != nil {
		return err
	}
	if _, err = p.offsets.Read(p.idx); err != nil {
		return err
	}
	return nil
}

type Row struct {
	// offset int
	data   []byte
	fields map[string]int
}

func (p *Pal) Fields() []string {
	return p.GetRandom().Fields()
}

func (r *Row) Fields() []string {
	keys := make([]string, 0, len(r.fields))
	for k := range r.fields {
		keys = append(keys, k)
	}
	return keys
}

func (r *Row) Get(field string) string {
	if fieldNum, ok := r.fields[field]; ok {
		headerSize := len(r.fields) * 4
		start := headerSize
		if fieldNum > 0 {
			start = headerSize + int(binary.LittleEndian.Uint32(r.data[(fieldNum-1)*4:4+(fieldNum-1)*4]))
		}
		end := headerSize + int(binary.LittleEndian.Uint32(r.data[fieldNum*4:4+fieldNum*4]))
		return string(r.data[start:end])
	}
	return ""
}

func (r *Row) String() string {
	var s []string
	for field := range r.fields {
		s = append(s, fmt.Sprintf("%s=%s", field, r.Get(field)))
	}
	return fmt.Sprintf("Row(%s)", strings.Join(s, ","))
}

func (p *Pal) Get(id string) *Row {
	b := p.offsets.Get([]byte(id))
	if b == nil {
		return nil
	}
	offset := int(binary.LittleEndian.Uint64(b))
	return &Row{data: p.data[offset:], fields: p.fields}
}

func (p *Pal) GetRandom() *Row {
	b := p.offsets.GetRandomValue()
	if b == nil {
		return nil
	}
	offset := int(binary.LittleEndian.Uint64(b))
	return &Row{data: p.data[offset:], fields: p.fields}
}
