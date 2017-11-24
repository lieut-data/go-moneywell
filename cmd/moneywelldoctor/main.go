package main

import (
	"flag"
	"fmt"

	"github.com/lieut-data/go-moneywell/internal/doctor"
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("required: path to moneywell document")
		return
	}

	err := doctor.Diagnose(flag.Arg(0))
	if err != nil {
		fmt.Printf("do failed: %v\n", err)
	}
}
