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

	. "github.com/smartystreets/goconvey/convey"
)

func TestPal_GetRandom(t *testing.T) {

}

func TestPal(t *testing.T) {



	Convey("pal", t, func() {

		data := [][]string{
			{"123", "test", "somevalue"},
			{"xfz56", "test2", "somevalue2"},
			{"00000", "test3", "somevalue3"},
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

			// get random row
			row = p.GetRandom()
			So(row, ShouldNotBeNil)
			So(row.Get("id"), ShouldBeIn, []string{"123", "xfz56", "00000"})

			// Get non-existent
			So(p.Get("aaaaaa"), ShouldBeNil)

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
	Convey(`build`, t, func() {
		for j := 0; j < 1024; j++ {
			Convey(fmt.Sprintf("%d", j), func() {
				b1 := NewBuilder([]string{"id", "value"})
				for i := 0; i < j; i++ {
					b1.Add(strconv.Itoa(i), []string{
						strconv.Itoa(i),
						strconv.Itoa(i),
					})
				}
				var buf bytes.Buffer

				err := b1.BuildTo(&buf)
				So(err, ShouldBeNil)
			})
		}
	})
	Convey("mmap", t, func() {
		Convey("v1", func() {
			p, err := MMapPal("testdata/v1.pal")
			So(err, ShouldBeNil)
			row := p.Get("2")
			So(row.Get("val2"), ShouldEqual, "2-2")
			So(p.Get("aaaaaa"), ShouldBeNil)
			p.Free()
		})
		Convey("v2", func() {
			p, err := MMapPal("testdata/v2.pal")
			So(err, ShouldBeNil)
			row := p.Get("2")
			So(row.Get("val2"), ShouldEqual, "2-2")
			So(p.Get("aaaaaa"), ShouldBeNil)
			t.Log(p.GetRandom())
			p.Free()
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
