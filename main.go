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
	PlayerPositions  []Coords
	TileLayout       []byte
	Units            []UnitData
}

type Coords struct {
	X int16
	Y int16
}

type Stats struct {
	HP   int16
	Atk  int16
	Spd  int16
	Def  int16
	Res  int16
	Unk1 int16
	Unk2 int16
	Unk3 int16
}

type UnitData struct {
	Id               string
	X                int16
	Y                int16
	Rarity           byte
	Level            byte
	TrueLevel        byte
	UnknownByte      byte
	SpecialCooldown  int8
	IsEnemy          bool
	MovementGroup    byte
	MovementDelay    int8
	StartTurn        int8
	GoBackToHomeTile bool
	BreakTerrain     bool
	Stats            Stats
	Skills           []string
}

func main() {
	// updater.Update()
	var mapData MapData = MapData{
		PlayerPositions: []Coords{},
	}
	byteArray, _ := os.ReadFile("S8084C.bin")
	var index = 0x41

	var totalPlayerUnits = readRawBytes(&byteArray, index, 4)
	var totalPlayersXor, _ = rawXor(&totalPlayerUnits, []byte{0x9a, 0xc7, 0x63, 0x9d})
	mapData.TotalPlayerUnits = byteArrayToInt32(&totalPlayersXor)
	index += 4

	var totalEnemiesBuffer = readRawBytes(&byteArray, index, 4)
	var totalEnemies, _ = rawXor(&totalEnemiesBuffer, []byte{0xee, 0x10, 0x67, 0xac})
	mapData.TotalEnemies = byteArrayToInt32(&totalEnemies)
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
	index += 48

	for i := 0; i < int(mapData.TotalPlayerUnits); i++ {
		var playerPosition Coords = Coords{}
		var rawXBytes = readRawBytes(&byteArray, index, 2)
		var unlockedXCoords, _ = rawXor(&rawXBytes, []byte{0x32, 0xb3})
		var int16_xCoord = byteArrayToInt16(&unlockedXCoords)
		playerPosition.X = int16_xCoord

		index += 2
		var rawYBytes = readRawBytes(&byteArray, index, 2)
		var unlockedYCoords, _ = rawXor(&rawYBytes, []byte{0xb2, 0x28})
		var int16_yCoord = byteArrayToInt16(&unlockedYCoords)
		playerPosition.Y = int16_yCoord
		mapData.PlayerPositions = append(mapData.PlayerPositions, playerPosition)
		index += 2
		index = skipNullBytes(&byteArray, index)
	}

	index = 0x189

	for i := 0; i < int(mapData.TotalEnemies); i++ {
		var unitStruct UnitData = UnitData{}

		var rawXCoordinates = readRawBytes(&byteArray, index, 2)
		var xCoord, _ = rawXor(&rawXCoordinates, []byte{0x32, 0xb3})
		var x = byteArrayToInt16(&xCoord)
		unitStruct.X = x
		index += 2

		var rawYCoordinates = readRawBytes(&byteArray, index, 2)
		var yCoord, _ = rawXor(&rawYCoordinates, []byte{0xb2, 0x28})
		var y = byteArrayToInt16(&yCoord)
		unitStruct.Y = y
		index += 2

		var rarityByte = readRawBytes(&byteArray, index, 1)
		var decryptedRarity = rarityByte[0] ^ 0x61
		unitStruct.Rarity = decryptedRarity
		index++

		var levelByte = readRawBytes(&byteArray, index, 1)
		var decryptedDisplayedLevel = levelByte[0] ^ 0x2A
		unitStruct.Level = decryptedDisplayedLevel
		index++

		var specialCooldownByte = readRawBytes(&byteArray, index, 1)
		var decryptedSpecialByte = specialCooldownByte[0] ^ 0x1E
		unitStruct.SpecialCooldown = int8(decryptedSpecialByte)
		index++

		var unk = readRawBytes(&byteArray, index, 1)
		unitStruct.UnknownByte = unk[0]
		index++

		stats, newIndex := readStats(&byteArray, index)
		unitStruct.Stats = stats
		index = newIndex

		var startTurnByte = readRawBytes(&byteArray, index, 1)
		var startTurn = int8(startTurnByte[0] ^ 0xcf)
		unitStruct.StartTurn = startTurn
		index++

		var movementGroupByte = readRawBytes(&byteArray, index, 1)
		unitStruct.MovementGroup = movementGroupByte[0] ^ 0xcf
		index++

		var movementDelay = readRawBytes(&byteArray, index, 1)
		unitStruct.MovementDelay = int8(movementDelay[0] ^ 0x95)
		fmt.Println(unitStruct.MovementDelay)
		index++

		var breakTerrainByte = readRawBytes(&byteArray, index, 1)
		var shouldBreakTerrain = breakTerrainByte[0]^0x71 != 0
		unitStruct.BreakTerrain = shouldBreakTerrain
		index++

		index = 0x1a5
		var tetherByte = readRawBytes(&byteArray, index, 1)
		var shouldGoBackToMainTile = tetherByte[0]^0xb8 != 0
		unitStruct.GoBackToHomeTile = shouldGoBackToMainTile
		index++

		var trueLevelByte = readRawBytes(&byteArray, index, 1)
		var trueLevel = trueLevelByte[0] ^ 0x85
		unitStruct.TrueLevel = trueLevel
		index++

		var isEnemyByte = readRawBytes(&byteArray, index, 1)
		var isEnemy = isEnemyByte[0]^0xd0 != 0
		unitStruct.IsEnemy = isEnemy

		fmt.Println(unitStruct)

		mapData.Units = append(mapData.Units, unitStruct)

		os.Exit(0)
	}

	fmt.Println(mapData)
}

func readStats(byteArray *[]byte, baseIndex int) (Stats, int) {
	var stats Stats = Stats{}
	var hpBytes = readRawBytes(byteArray, baseIndex, 2)
	baseIndex += 2
	var atkBytes = readRawBytes(byteArray, baseIndex, 2)
	baseIndex += 2
	var spdBytes = readRawBytes(byteArray, baseIndex, 2)
	baseIndex += 2
	var defBytes = readRawBytes(byteArray, baseIndex, 2)
	baseIndex += 2
	var resBytes = readRawBytes(byteArray, baseIndex, 2)
	baseIndex += 2
	var xorHP, _ = rawXor(&hpBytes, []byte{0x32, 0xd6})
	var xorAtk, _ = rawXor(&atkBytes, []byte{0xa0, 0x14})
	var xorSpd, _ = rawXor(&spdBytes, []byte{0x5e, 0xa5})
	var xorDef, _ = rawXor(&defBytes, []byte{0x66, 0x85})
	var xorRes, _ = rawXor(&resBytes, []byte{0xe5, 0xae})

	var unk1Bytes = readRawBytes(byteArray, baseIndex, 2)
	stats.Unk1 = byteArrayToInt16(&unk1Bytes)
	baseIndex += 2

	var unk2Bytes = readRawBytes(byteArray, baseIndex, 2)
	stats.Unk2 = byteArrayToInt16(&unk2Bytes)
	baseIndex += 2

	var unk3Bytes = readRawBytes(byteArray, baseIndex, 2)
	stats.Unk3 = byteArrayToInt16(&unk3Bytes)
	baseIndex += 2

	stats.HP = byteArrayToInt16(&xorHP)
	stats.Atk = byteArrayToInt16(&xorAtk)
	stats.Spd = byteArrayToInt16(&xorSpd)
	stats.Def = byteArrayToInt16(&xorDef)
	stats.Res = byteArrayToInt16(&xorRes)

	return stats, baseIndex
}

// Applies a XOR byte array, allowing NULL bytes to appear as a result of the operation.
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

// Read bytes until a NULL byte is hit, or if two consecutive NULL bytes appear,
// returning only the first. Returns the index at which the read stopped,
// as well as the collected byte sequence.
// Adjusts the returned index to prevent off-by-one errors (nth byte is (n - 1) indexed).
func readBytes(byteArray *[]byte, index int) (int, []byte) {
	var subArray []byte = []byte{} // optimization, sort of a guess
	var lastByte byte = 0xff
	var curIndex = index

	for {
		currentByte := (*byteArray)[curIndex]
		if (lastByte == 0 && currentByte == 0) || (lastByte != 0 && currentByte == 0) {
			return curIndex + 1, subArray
		}
		subArray = append(subArray, currentByte)
		lastByte = currentByte
		curIndex++
	}
}

// Skips consecutive NULL bytes and returns the index at which the first non-NULL byte is found.
// Adjusts the returned index to prevent off-by-one errors (nth byte is (n - 1) indexed)
func skipNullBytes(byteArray *[]byte, currentIndex int) int {
	var i = 0
	for {
		if (*byteArray)[currentIndex+i] != 0 {
			return currentIndex + i + 1
		}

		i++
	}
}

// Reads bytes as they are, taking NULL bytes. Returns the read bytes.
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
		var keyByte = key[i%len(encoded)]
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

func byteArrayToInt16(byteArray *[]byte) int16 {
	var value int16
	value |= int16((*byteArray)[0])
	value |= int16((*byteArray)[1]) << 8

	return value
}
