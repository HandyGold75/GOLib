package main

import (
	"Gonos"
	"fmt"
)

func main() {
	zps, err := Gonos.ScanZonePlayer("10.69.3.0/24")
	if err != nil {
		panic(err)
	}
	zp := zps[0]

	out, err := zp.GetQue()
	if err != nil {
		fmt.Println(err)
		fmt.Println()
	}
	fmt.Println(out)
}
