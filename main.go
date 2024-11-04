package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var XOR_ID []byte = []byte{0x81, 0x00, 0x80, 0xA4, 0x5A, 0x16, 0x6F, 0x78, 0x57, 0x81, 0x2D, 0xF7, 0xFC, 0x66, 0x0F, 0x27, 0x75, 0x35, 0xB4, 0x34, 0x10, 0xEE, 0xA2, 0xDB, 0xCC, 0xE3, 0x35, 0x99, 0x43, 0x48, 0xD2, 0xBB, 0x93, 0xC1}

type MapData struct {
	LastEnemyTurn    bool
	Width            int32
	Height           int32
	BaseTerrain      int8
	TurnsToDefend    byte
	TurnsToWin       byte
	TotalEnemies     int32
	TotalPlayerUnits int32
	TileLayout       []byte
}

func main() {
	// updater.Update()
	var mapData MapData = MapData{}
	byteArray, _ := os.ReadFile("S8084C.bin")
	var index = 0x41

	var totalPlayerUnits = readRawBytes(&byteArray, index, 4)
	var totalPlayersXor, _ = rawXor(&totalPlayerUnits, []byte{0x9a, 0xc7, 0x63, 0x9d})
	mapData.TotalPlayerUnits = byteArrayToInt32(&totalPlayersXor)
	index += 4

	var totalEnemiesBuffer = readRawBytes(&byteArray, index, 4)
	fmt.Println(totalEnemiesBuffer)
	var totalEnemies, _ = rawXor(&totalEnemiesBuffer, []byte{0xee, 0x10, 0x67, 0xac})
	mapData.TotalEnemies = byteArrayToInt32(&totalEnemies)
	fmt.Println(totalEnemies)
	index += 4

	var turnsToWin = readRawBytes(&byteArray, index, 1)[0] ^ 0xFD
	mapData.TurnsToWin = turnsToWin
	index++

	var lastEnemyTurn = (readRawBytes(&byteArray, index, 1)[0] ^ 0xC7) != 0
	mapData.LastEnemyTurn = lastEnemyTurn
	index++

	var turnsToDefend = readRawBytes(&byteArray, index, 1)[0] ^ 0xEC
	mapData.TurnsToDefend = turnsToDefend

	index = 0x59
	var widthSlice = readRawBytes(&byteArray, index, 4)
	var mapWidth, err = rawXor(&widthSlice, []byte{0x5f, 0xd7, 0x7c, 0x6b})

	if err != nil {
		log.Fatalln(err)
	}

	index += 4
	var heightSlice = readRawBytes(&byteArray, index, 4)
	var mapHeight, _ = rawXor(&heightSlice, []byte{0xd5, 0x12, 0xaa, 0x2b})
	mapData.Width = byteArrayToInt32(&mapWidth)
	mapData.Height = byteArrayToInt32(&mapHeight)

	index += 4
	var terrain = readRawBytes(&byteArray, index, 1)
	mapData.BaseTerrain = int8(terrain[0] ^ 0x41)

	index = skipNullBytes(&byteArray, index)

	var tileBytes = readRawBytes(&byteArray, index, 48)
	var tilesXor [48]byte = [48]byte{0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1, 0xa1}
	var tiles, _ = rawXor(&tileBytes, tilesXor[:])
	mapData.TileLayout = tiles

	fmt.Println(mapData)
}

func rawXor(byteArray *[]byte, xorKey []byte) ([]byte, error) {
	if len(*byteArray) != len(xorKey) {
		var err = errors.ErrUnsupported
		return nil, err
	}
	var withXor []byte = []byte{}
	for i, xorByte := range xorKey {
		withXor = append(withXor, (*byteArray)[i]^xorByte)
	}

	return withXor, nil
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
			return currentIndex + i + 1
		}

		i++
	}
}

func readRawBytes(byteArray *[]byte, curIndex int, length int) []byte {
	var subarray []byte = []byte{}
	for i := 0; i < length; i++ {
		subarray = append(subarray, (*byteArray)[curIndex-1+i])
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

func byteArrayToInt32(byteArray *[]byte) int32 {
	var value int32
	value |= int32((*byteArray)[0])
	value |= int32((*byteArray)[1]) << 8
	value |= int32((*byteArray)[2]) << 16
	value |= int32((*byteArray)[3]) << 24

	return value
}
