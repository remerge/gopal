package gopal

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPal(t *testing.T) {
	Convey("pal", t, func() {

		data := [][]string{
			[]string{"123", "test", "somevalue"},
			[]string{"xfz56", "test2", "somevalue2"},
			[]string{"00000", "test3", "somevalue3"},
		}

		b := NewBuilder([]string{"id", "name", "value"})
		So(b, ShouldNotBeNil)
		for _, values := range data {
			b.Add(values[0], values)
		}

		Convey("build pal and load from bytes", func() {
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			b.BuildTo(w)
			w.Flush()

			p := &Pal{}
			p.From(buf.Bytes())

			row := p.Get("123")
			So(row, ShouldNotBeNil)

			So(row.Get("id"), ShouldEqual, "123")
			So(row.Get("name"), ShouldEqual, "test")
			So(row.Get("value"), ShouldEqual, "somevalue")
		})

		Convey("build pal and mmap", func() {
			fn := "mmap.pal"
			os.Remove(fn)
			var buf bytes.Buffer
			w := bufio.NewWriter(&buf)
			b.BuildTo(w)
			w.Flush()

			err := ioutil.WriteFile(fn, buf.Bytes(), os.ModePerm)
			So(err, ShouldBeNil)

			p, err := MMapPal(fn)
			So(err, ShouldBeNil)
			row := p.Get("123")
			So(row, ShouldNotBeNil)
			So(row.Get("id"), ShouldEqual, "123")
			So(row.Get("name"), ShouldEqual, "test")
			So(row.Get("value"), ShouldEqual, "somevalue")

			p.Free()
			os.Remove(fn)
		})
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
