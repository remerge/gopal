package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	pal "github.com/remerge/gopal"
)

func csv2pal(csvname, palname, _sep string) error {
	rr := strings.NewReader(_sep)
	sep, _, _ := rr.ReadRune()
	f, err := os.Open(csvname)
	defer f.Close()
	if err != nil {
		return err
	}
	r := csv.NewReader(f)

	r.Comma = sep
	headerRead := false
	var pb *pal.Builder
	for {
		record, err := r.Read()
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
		return err
	}
	defer out.Close()
	pb.BuildTo(out)
	return nil
}

func main() {
	csvfile := os.Args[1]
	palfile := os.Args[2]

	err := csv2pal(csvfile, palfile, ",")
	if err != nil {
		fmt.Println(err)
	}

}
