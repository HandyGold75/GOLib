package main

import (
	"fmt"
	"time"

	"github.com/HandyGold75/gonos"
)

func main() {
	zp := gonos.NewZonePlayer("10.0.1.32")

	zp.SetVolume(10)
	volume, _ := zp.GetVolume()
	fmt.Println(volume)

	time.Sleep(time.Second * time.Duration(5))

	zp.SetMute(true)
	isMute, _ := zp.GetMute()
	fmt.Println(isMute)
}
