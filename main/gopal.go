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

	d := p.Get("")
	fmt.Println(d)
	return nil
}

func main() {
	palFile := os.Args[1]

	err := probePal(palFile)
	if err != nil {
		fmt.Println(err)
	}

}
