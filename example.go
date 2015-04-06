package main

import (
	"fmt"
	"github.com/MMore/go-openal/openal"
)

func main() {
	data := MustAsset("resources/de/f/1s.wav")
	// data2 := MustAsset("resources/de/f/2s.wav")
	go PlaySoundFromMemory(data)
	// go PlaySoundFromMemory(data2)
	// PlaySound("ahooga")
	// PlaySound("welcome")

	var s string
	fmt.Scanln(&s)
}
