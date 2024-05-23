package main

import (
	"fmt"
	"time"

	"github.com/HandyGold75/GOLib/gonos"
)

func main() {
	zp := gonos.NewZonePlayer("10.0.1.32")

	ti := zp.GetTrackInfo()

	fmt.Println("-----")
	fmt.Printf("QuePosition: %v\n", ti.QuePosition)
	fmt.Printf("Duration: %v\n", ti.Duration)
	fmt.Printf("URI: %v\n", ti.URI)
	fmt.Printf("Progress: %v\n", ti.Progress)
	fmt.Printf("AlbumArtURI: %v\n", ti.AlbumArtURI)
	fmt.Printf("Title: %v\n", ti.Title)
	fmt.Printf("Class: %v\n", ti.Class)
	fmt.Printf("Creator: %v\n", ti.Creator)
	fmt.Printf("Album: %v\n", ti.Album)

	time.Sleep(time.Second * time.Duration(5))

	qi := zp.GetQueInfo()

	fmt.Println("-----")
	fmt.Printf("Count: %v\n", qi.Count)
	fmt.Printf("TotalCount: %v\n", qi.TotalCount)

	for _, track := range qi.Tracks {
		fmt.Println("---")
		fmt.Printf("AlbumArtURI: %v\n", track.AlbumArtURI)
		fmt.Printf("Title: %v\n", track.Title)
		fmt.Printf("Class: %v\n", track.Class)
		fmt.Printf("Creator: %v\n", track.Creator)
		fmt.Printf("Album: %v\n", track.Album)
	}
}
