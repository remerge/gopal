package main

import (
	"fmt"
	"os"

	pal "github.com/remerge/gopal"
)

func probePal(file string) error {
	p, err := pal.MMapPal(file)
	if err != nil {
		return err
	}
	fields := p.Fields()
	fmt.Printf("Pal %v opened \n", file)
	fmt.Printf("Fields are %v \n", fields)

	for i := 0; i < 100000; i++ {
		for _, f := range fields {
			field := p.GetRandom().Get(f)
			fmt.Println(field)
		}
		p.GetRandom().Get("xz")
	}
	return nil
}

func main() {
	palFile := os.Args[1]

	err := probePal(palFile)
	if err != nil {
		fmt.Println(err)
	}

}
