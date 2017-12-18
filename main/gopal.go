package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	pal "github.com/remerge/gopal"
)

func csv2pal(csvname, palname, _sep string) (err error) {
	rr := strings.NewReader(_sep)
	sep, _, _ := rr.ReadRune()
	f, err := os.Open(csvname)
	defer f.Close()
	if err != nil {
		return
	}
	r := csv.NewReader(f)

	r.Comma = sep
	headerRead := false
	var pb *pal.Builder
	for {
		var record []string
		record, err = r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("err", err)
			continue
		}
		if !headerRead {
			headerRead = true
			pb = pal.NewBuilder(record)
		} else {
			pb.AddRow(record)
		}
	}
	out, err := os.Create(palname)
	if err != nil {
		return
	}
	defer out.Close()
	err = pb.BuildTo(out)
	return
}

func main() {
	csvfile := os.Args[1]
	palfile := os.Args[2]

	err := csv2pal(csvfile, palfile, ",")
	if err != nil {
		fmt.Println(err)
	}

}
