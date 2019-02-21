package gopal

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPal(t *testing.T) {
	t.Run(`build`, func(t *testing.T) {
		for j := 0; j < 1024; j++ {
			t.Run(fmt.Sprintf("%d", j), func(t *testing.T) {
				b1 := NewBuilder([]string{"id", "value"})
				for i := 0; i < j; i++ {
					b1.Add(strconv.Itoa(i), []string{
						strconv.Itoa(i),
						strconv.Itoa(i),
					})
				}
				var buf bytes.Buffer

				err := b1.BuildTo(&buf)
				assert.Nil(t, err)
			})
		}
	})

	t.Run("mmap", func(t *testing.T) {
		t.Run("v1", func(t *testing.T) {
			p, err := MMapPal("testdata/v1.pal")
			assert.Nil(t, err)
			row := p.Get("2")
			assert.Equal(t, "2-2", row.Get("val2"))
			assert.Nil(t, p.Get("aaaaaa"))
			p.Free()
		})
		t.Run("v2", func(t *testing.T) {
			p, err := MMapPal("testdata/v2.pal")
			assert.Nil(t, err)
			row := p.Get("2")
			assert.Equal(t, "2-2", row.Get("val2"))
			assert.Nil(t, p.Get("aaaaaa"))
			t.Log(p.GetRandom())
			p.Free()
		})
	})

	data := [][]string{
		{"123", "test", "somevalue"},
		{"xfz56", "test2", "somevalue2"},
		{"00000", "test3", "somevalue3"},
	}

	t.Run("build pal and load from bytes", func(t *testing.T) {
		b := NewBuilder([]string{"id", "name", "value"})
		assert.NotNil(t, b)
		for _, values := range data {
			b.Add(values[0], values)
		}

		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		b.BuildTo(w)
		w.Flush()

		p := &Pal{}
		p.From(buf.Bytes())

		row := p.Get("123")
		assert.NotNil(t, row)

		assert.Equal(t, "123", row.Get("id"))
		assert.Equal(t, "test", row.Get("name"))
		assert.Equal(t, "somevalue", row.Get("value"))

		// get random row
		row = p.GetRandom()
		assert.NotNil(t, row)
		assert.Contains(t, []string{"123", "xfz56", "00000"}, row.Get("id"))

		// Get non-existent
		assert.Nil(t, p.Get("aaaaaa"))
	})

	t.Run("build pal and mmap", func(t *testing.T) {
		b := NewBuilder([]string{"id", "name", "value"})
		assert.NotNil(t, b)
		for _, values := range data {
			b.Add(values[0], values)
		}

		fn := "mmap.pal"
		os.Remove(fn)
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		b.BuildTo(w)
		w.Flush()

		err := ioutil.WriteFile(fn, buf.Bytes(), os.ModePerm)
		assert.Nil(t, err)

		p, err := MMapPal(fn)
		assert.Nil(t, err)
		row := p.Get("123")
		assert.NotNil(t, row)
		assert.Equal(t, "123", row.Get("id"))
		assert.Equal(t, "test", row.Get("name"))
		assert.Equal(t, "somevalue", row.Get("value"))

		p.Free()
		os.Remove(fn)
	})
}

func BenchmarkPalGen(b *testing.B) {
	for n := 0; n < b.N; n++ {
		builder := NewBuilder([]string{"id", "name", "value", "anotherfield"})
		values := []string{"id", "name", "value", "anotherfield"}

		// 1000000 entries
		for i := 0; i < 1000000; i++ {
			builder.Add(strconv.Itoa(i), values)
		}
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)
		builder.BuildTo(w)
		w.Flush()
	}
}
func BenchmarkPal(b *testing.B) {
	builder := NewBuilder([]string{"id", "name", "value", "anotherfield"})
	values := []string{"id", "name", "value", "anotherfield"}
	// 1000000 entries
	for i := 0; i < 1000000; i++ {
		builder.Add(strconv.Itoa(i), values)
	}
	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	builder.BuildTo(w)
	w.Flush()

	p := &Pal{}
	p.From(buf.Bytes())

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		p.Get(strconv.Itoa(rand.Intn(1000000))).Get("value")
	}
}
