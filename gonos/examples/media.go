package main

import (
	"fmt"
	"time"

	"github.com/HandyGold75/gonos"
)

func main() {
	zp := gonos.NewZonePlayer("10.0.1.2")

	zp.SetPlayMode(false, false, false)
	mode, _ := zp.GetPlayMode()
	fmt.Println(mode)

	time.Sleep(time.Second * time.Duration(5))

	zp.SetShuffle(true)
	shuffling, _ := zp.GetShuffle()
	fmt.Println(shuffling)

	time.Sleep(time.Second * time.Duration(5))

	zp.SetRepeat(true)
	repeating, _ := zp.GetRepeat()
	fmt.Println(repeating)

	time.Sleep(time.Second * time.Duration(5))

	zp.SetRepeatOne(true)
	repeatingOnce, _ := zp.GetRepeatOne()
	fmt.Println(repeatingOnce)
}
