package main

import (
	"Gonos"
	"fmt"
	"time"
)

func main() {
	zps, err := Gonos.ScanZonePlayer("10.69.3.0/24", time.Second)
	if err != nil {
		panic(err)
	}
	zp := zps[0]

	err = zp.AVTransport.Pause()
	if err != nil {
		fmt.Println(err)
		return
	}
}
