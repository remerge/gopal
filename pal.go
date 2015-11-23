package gopal

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/alecthomas/mph"
)

type Pal struct {
	fields map[string]int
	idx    []byte
	data   []byte
	chd    *mph.CHD
	mmaped bool
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
	mmap, err := syscall.Mmap(int(file.Fd()), 0, int(fi.Size()), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	err = p.From(mmap)
	if err != nil {
		syscall.Munmap(mmap)
		return nil, err
	}
	p.mmaped = true
	runtime.SetFinalizer(p, func(p *Pal) {
		p.Free()
	})
	return p, nil
}

func (p *Pal) Free() {
	if !p.mmaped {
		return
	}
	if p.data == nil {
		return
	}
	syscall.Munmap(p.data)
	p.data = nil
	p.mmaped = false
}

func (p *Pal) From(b []byte) error {
	if len(b) < int(unsafe.Sizeof(PalHeader{})) {
		return errors.New("can't load pal, seems to be too small")
	}
	h := &PalHeader{}
	h.Read(b)
	if !h.Valid() {
		return errors.New("can't load pal, header invalid")
	}

	// TODO  - checksum the file or at least check if the size matches

	hm := h.HeadSize + h.MapSize
	hmi := hm + h.IdxSize
	p.idx = b[hm:hmi]
	p.data = b[hmi : len(b)-1]
	dec := gob.NewDecoder(bytes.NewBuffer(b[h.HeadSize:hm]))
	err := dec.Decode(&p.fields)
	if err != nil {
		return err
	}
	chd, err := mph.Mmap(p.idx)
	if err != nil {
		return err
	}
	p.chd = chd
	return nil
}

type Row struct {
	// offset int
	data   []byte
	fields map[string]int
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

func (p *Pal) Get(id string) *Row {
	b := p.chd.Get([]byte(id))
	if b == nil {
		return nil
	}
	offset := int(binary.LittleEndian.Uint64(b))
	return &Row{data: p.data[offset:], fields: p.fields}
}
