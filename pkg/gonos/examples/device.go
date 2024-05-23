package main

import (
	"fmt"
	"time"

	"github.com/HandyGold75/gonos"
)

func main() {
	zp := gonos.NewZonePlayer("10.0.1.2")

	zp.SetBass(4)
	bass, _ := zp.GetBass()
	fmt.Println(bass)

	time.Sleep(time.Second * time.Duration(5))

	zp.SetTreble(2)
	treble, _ := zp.GetTreble()
	fmt.Println(treble)

	time.Sleep(time.Second * time.Duration(5))

	zp.SetLoudness(true)
	loudness, _ := zp.GetLoudness()
	fmt.Println(loudness)

	time.Sleep(time.Second * time.Duration(5))

	zp.SetLEDState(true)
	ledState, _ := zp.GetLedState()
	fmt.Println(ledState)

	time.Sleep(time.Second * time.Duration(5))

	zp.SetPlayerName("gonos")
	name, _ := zp.GetPlayerName()
	fmt.Println(name)
}
