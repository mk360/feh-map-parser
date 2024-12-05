package main

import (
	"feh-map-editor/decoder"
	"feh-map-editor/updater"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "update" {
		updater.Update()
	} else {
		decoder.Decode("S8084C.bin")
	}
}
