package main

import (
	"fmt"
	"time"

	"github.com/HandyGold75/GOLib/gonos"
)

func main() {
	zp := gonos.NewZonePlayer("10.0.1.2")

	zp.PlayFromQue(0)
	ti, _ := zp.GetTrackInfo()
	fmt.Println(ti.QuePosition)

	time.Sleep(time.Second * time.Duration(5))

	zp.RemoveFromQue(1)
	qi, _ := zp.GetQueInfo()
	fmt.Println(qi.Tracks)

	time.Sleep(time.Second * time.Duration(5))

	zp.ClearQue()
	qi, _ = zp.GetQueInfo()
	fmt.Println(qi.Tracks)
}
