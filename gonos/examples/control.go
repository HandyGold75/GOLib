package main

import (
	"fmt"
	"time"

	"github.com/HandyGold75/gonos"
)

func main() {
	zp := gonos.NewZonePlayer("10.0.1.2")

	zp.Play()
	playing, _ := zp.GetState()
	fmt.Println(playing)

	time.Sleep(time.Second * time.Duration(5))

	zp.Pause()
	playing, _ = zp.GetState()
	fmt.Println(playing)

	time.Sleep(time.Second * time.Duration(5))

	zp.Next()
	ti, _ := zp.GetTrackInfo()
	fmt.Println(ti.TItle)

	time.Sleep(time.Second * time.Duration(5))

	zp.Stop()
	playing, _ = zp.GetState()
	fmt.Println(playing)

	time.Sleep(time.Second * time.Duration(5))

	zp.Previous()
	ti, _ = zp.GetTrackInfo()
	fmt.Println(ti.TItle)

	time.Sleep(time.Second * time.Duration(5))

	zp.Seek("00:00:30")
	ti, _ = zp.GetTrackInfo()
	fmt.Println(ti.Progress)
}
