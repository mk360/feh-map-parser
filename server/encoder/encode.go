package encoder

import (
	"encoding/binary"
	"errors"
	"feh-map-editor/loader"
	"os"
)

type EncodePayload struct {
	BaseTerrain      int     `json:"BaseTerrain"`
	FileHeader       []int   `json:"FileHeader"`
	Height           int     `json:"Height"`
	ID               string  `json:"Id"`
	TurnsToDefend    int     `json:"TurnsToDefend"`
	TurnsToWin       int     `json:"TurnsToWin"`
	TileLayout       string  `json:"TileLayout"`
	Units            []Units `json:"Units"`
	PlayerPositions  []any   `json:"PlayerPositions"`
	Width            int     `json:"Width"`
	TotalEnemies     int     `json:"TotalEnemies"`
	TotalPlayerUnits int     `json:"TotalPlayerUnits"`
	LastEnemyTurn    bool    `json:"LastEnemyTurn"`
}
type Stats struct {
	Hp  int `json:"hp"`
	Atk int `json:"atk"`
	Spd int `json:"spd"`
	Def int `json:"def"`
	Res int `json:"res"`
}
type Spawning struct {
	DependencyHero          string `json:"dependencyHero"`
	RemainingBeforeSpawning int    `json:"remainingBeforeSpawning"`
	DefeatBeforeSpawning    int    `json:"defeatBeforeSpawning"`
}
type Ai struct {
	BreakTerrain     bool `json:"breakTerrain"`
	GoBackToHomeTile bool `json:"goBackToHomeTile"`
	MovementDelay    int  `json:"movementDelay"`
	MovementGroup    int  `json:"movementGroup"`
	StartTurn        int  `json:"startTurn"`
}
type Units struct {
	Skills          []string `json:"skills"`
	Stats           Stats    `json:"stats"`
	Name            string   `json:"name"`
	Rarity          int      `json:"rarity"`
	IsEnemy         bool     `json:"isEnemy"`
	TrueLevel       int      `json:"trueLevel"`
	Level           int      `json:"level"`
	Column          int      `json:"column"`
	Row             int      `json:"row"`
	Unknown         int      `json:"unknown"`
	Spawning        Spawning `json:"spawning"`
	Ai              Ai       `json:"AI"`
	SpecialCooldown int      `json:"specialCooldown"`
}

func IntsToBytes(ints []int) []byte {
	bytes := make([]byte, len(ints))
	for i, v := range ints {
		bytes[i] = byte(v)
	}
	return bytes
}

// Encode takes a MapData struct and writes it to a file at the specified path
func Encode(mapData EncodePayload, filePath string) error {
	loader.LoadJSONs()

	// Create a buffer to hold the file data
	fileData := make([]byte, 0, 4096) // Initial capacity, will grow as needed

	fileData = append(fileData, IntsToBytes(mapData.FileHeader)...)
	fileData = append(fileData, 0xa8, 0xae, 0x54, 0x16, 0x47, 0x50, 0xe2, 0xae, 0x30, 0, 0, 0, 0, 0, 0, 0, 0x78, 0, 0, 0, 0, 0, 0, 0, 0x98, 0, 0, 0, 0, 0, 0, 0)
	// player positions: currently not supported
	fileData = append(fileData, 0x9a, 0xc7, 0x63, 0x9d)

	var unitCount = byte(mapData.TotalEnemies + mapData.TotalPlayerUnits)
	var unitCountByteArray = []byte{unitCount, 0, 0, 0}
	var xorUnitCount, _ = rawXor(&unitCountByteArray, []byte{0xee, 0x10, 0x67, 0xac})

	fileData = append(fileData, xorUnitCount...)

	// turns to win
	var turnsToWinByte = byte(mapData.TurnsToWin ^ 0xfd)
	fileData = append(fileData, turnsToWinByte)

	// last enemy phase
	var enemyTurnInt = 0
	if mapData.LastEnemyTurn {
		enemyTurnInt = 1
	}

	var lastEnemyPhaseByte = byte(enemyTurnInt ^ 0xc7)
	fileData = append(fileData, lastEnemyPhaseByte)

	var turnsToDefendByte = byte(mapData.TurnsToDefend ^ 0xec)
	var baseTerrainByte = byte(mapData.BaseTerrain ^ 0x41)
	fileData = append(fileData, turnsToDefendByte, baseTerrainByte, 0, 0, 0, 0)

	// missing: pointer to unit data

	fileData = append(fileData, 0, 0, 0, 0, 0, 0)

	// currently we only support a standard 6x8 grid, so i can directly write
	// the xor in the buffer
	fileData = append(fileData, 0x59, 0xd7, 0x7c, 0x6b)
	fileData = append(fileData, 0xdd, 0x12, 0xaa, 0x2b)
	// unknown?
	fileData = append(fileData, 0xbe)

	return os.WriteFile(filePath, fileData, 0644)
	fileData = append(fileData, 0, 0, 0, 0, 0, 0, 0)

}

// encodeStats encodes unit stats and appends them to the file data
func encodeStats(fileData []byte, stats Stats) []byte {
	// Convert stats to bytes
	hpBytes := make([]byte, 2)
	atkBytes := make([]byte, 2)
	spdBytes := make([]byte, 2)
	defBytes := make([]byte, 2)
	resBytes := make([]byte, 2)

	binary.LittleEndian.PutUint16(hpBytes, uint16(stats.Hp))
	binary.LittleEndian.PutUint16(atkBytes, uint16(stats.Atk))
	binary.LittleEndian.PutUint16(spdBytes, uint16(stats.Spd))
	binary.LittleEndian.PutUint16(defBytes, uint16(stats.Def))
	binary.LittleEndian.PutUint16(resBytes, uint16(stats.Res))

	// XOR the stats
	encryptedHP, _ := rawXor(&hpBytes, []byte{0x32, 0xd6})
	encryptedAtk, _ := rawXor(&atkBytes, []byte{0xa0, 0x14})
	encryptedSpd, _ := rawXor(&spdBytes, []byte{0x5e, 0xa5})
	encryptedDef, _ := rawXor(&defBytes, []byte{0x66, 0x85})
	encryptedRes, _ := rawXor(&resBytes, []byte{0xe5, 0xae})

	// Add stats to file data
	fileData = append(fileData, encryptedHP...)
	fileData = append(fileData, encryptedAtk...)
	fileData = append(fileData, encryptedSpd...)
	fileData = append(fileData, encryptedDef...)
	fileData = append(fileData, encryptedRes...)

	// Add unknown stats
	unk1Bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(unk1Bytes, uint16(0))
	fileData = append(fileData, unk1Bytes...)

	unk2Bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(unk2Bytes, uint16(0))
	fileData = append(fileData, unk2Bytes...)

	unk3Bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(unk3Bytes, uint16(0))
	fileData = append(fileData, unk3Bytes...)

	return fileData
}

// rawXor applies a XOR byte array, matching the decoder's function
func rawXor(byteArray *[]byte, xorKey []byte) ([]byte, error) {
	if len(*byteArray) != len(xorKey) {
		return nil, errors.ErrUnsupported
	}

	withXor := make([]byte, len(*byteArray))
	for i, xorByte := range xorKey {
		withXor[i] = (*byteArray)[i] ^ xorByte
	}

	return withXor, nil
}

// encodeOrDecodeString encodes or decodes a string using the XOR key
func encodeOrDecodeString(plaintext []byte, key []byte) []byte {
	encryptedArray := make([]byte, len(plaintext))
	for i, curByte := range plaintext {
		keyByte := key[i%len(key)]
		if keyByte == curByte {
			encryptedArray[i] = keyByte
		} else {
			encrypted := keyByte ^ curByte
			encryptedArray[i] = encrypted
		}
	}

	return encryptedArray
}
