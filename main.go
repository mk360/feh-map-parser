package main

import (
	"feh-map-editor/updater"
	"os"
)

var XOR_ID []byte = []byte{0x81, 0x00, 0x80, 0xA4, 0x5A, 0x16, 0x6F, 0x78, 0x57, 0x81, 0x2D, 0xF7, 0xFC, 0x66, 0x0F, 0x27, 0x75, 0x35, 0xB4, 0x34, 0x10, 0xEE, 0xA2, 0xDB, 0xCC, 0xE3, 0x35, 0x99, 0x43, 0x48, 0xD2, 0xBB, 0x93, 0xC1}

func main() {
	updater.Update()
	byteArray, _ := os.ReadFile("S8084C.bin")

}

func readBytes(byteArray *[]byte, index int) (int, []byte) {
	var subArray []byte = []byte{} // optimization, sort of a guess
	var lastByte byte = 0xff
	var curIndex = index

	for {
		currentByte := (*byteArray)[curIndex]
		if (lastByte == 0 && currentByte == 0) || (lastByte != 0 && currentByte == 0) {
			return curIndex, subArray
		}
		subArray = append(subArray, currentByte)
		lastByte = currentByte
		curIndex++
	}
}

func skipNullBytes(byteArray *[]byte, currentIndex int) int {
	var i = 0
	for {
		if (*byteArray)[currentIndex+i] != 0 {
			return currentIndex + i
		}

		i++
	}
}

func readRawBytes(byteArray *[]byte, curIndex int, length int) []byte {
	var subarray []byte = []byte{}
	for i := 0; i < length; i++ {
		subarray = append(subarray, (*byteArray)[curIndex+i])
	}

	return subarray
}

func encodeOrDecodeString(encoded []byte, key []byte) []byte {
	var decryptedArray []byte = make([]byte, len(encoded))
	for i, curByte := range encoded {
		var keyByte = key[i]
		if keyByte == curByte {
			decryptedArray[i] = keyByte
		} else {
			var decrypted = keyByte ^ curByte
			decryptedArray[i] = decrypted
		}
	}

	return decryptedArray
}

// func parseBinaryFile(fileData []byte) {
// 	var curIndex = 0x340

// }
