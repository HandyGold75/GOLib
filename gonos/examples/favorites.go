package main

import (
	"fmt"
	"time"

	"github.com/HandyGold75/gonos"
)

func main() {
	zp := gonos.NewZonePlayer("10.0.1.2")

	fi := zp.GetFavoritesInfo()

	fmt.Println("-----")
	fmt.Printf("Count: %v\n", fi.Count)
	fmt.Printf("TotalCount: %v\n", fi.TotalCount)

	for _, favorites := range fi.Favorites {
		fmt.Println("---")
		fmt.Printf("AlbumArtURI: %v\n", favorites.AlbumArtURI)
		fmt.Printf("Title: %v\n", favorites.Title)
		fmt.Printf("Description: %v\n", favorites.Description)
		fmt.Printf("Class: %v\n", favorites.Class)
		fmt.Printf("Type: %v\n", favorites.Type)
	}

	time.Sleep(time.Second * time.Duration(5))

	fi = zp.GetFavoritesRadioStationsInfo()

	fmt.Println("-----")
	fmt.Printf("Count: %v\n", fi.Count)
	fmt.Printf("TotalCount: %v\n", fi.TotalCount)

	for _, favorites := range fi.Favorites {
		fmt.Println("---")
		fmt.Printf("AlbumArtURI: %v\n", favorites.AlbumArtURI)
		fmt.Printf("Title: %v\n", favorites.Title)
		fmt.Printf("Description: %v\n", favorites.Description)
		fmt.Printf("Class: %v\n", favorites.Class)
		fmt.Printf("Type: %v\n", favorites.Type)
	}

	time.Sleep(time.Second * time.Duration(5))

	fi = zp.GetFavoritesRadioShowsInfo()

	fmt.Println("-----")
	fmt.Printf("Count: %v\n", fi.Count)
	fmt.Printf("TotalCount: %v\n", fi.TotalCount)

	for _, favorites := range fi.Favorites {
		fmt.Println("---")
		fmt.Printf("AlbumArtURI: %v\n", favorites.AlbumArtURI)
		fmt.Printf("Title: %v\n", favorites.Title)
		fmt.Printf("Description: %v\n", favorites.Description)
		fmt.Printf("Class: %v\n", favorites.Class)
		fmt.Printf("Type: %v\n", favorites.Type)
	}
}
