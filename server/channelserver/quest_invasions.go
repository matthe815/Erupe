package channelserver

import (
	"erupe-ce/common/byteframe"
	"erupe-ce/common/decryption"
	"erupe-ce/network/mhfpacket"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

func doInvasionChance(s *Session) {
	roll := rand.Intn(100)

	if roll < s.server.erupeConfig.GameplayOptions.EnhancedInvasions.InvasionChance {
		s.spawnInvasion = true
	}
}

func canInvasionHappen(s *Session, p mhfpacket.MHFPacket) bool {
	pkt := p.(*mhfpacket.MsgSysGetFile)

	data, _ := os.ReadFile(filepath.Join(s.server.erupeConfig.BinPath, fmt.Sprintf("quests/%s.bin", pkt.Filename)))
	decrypted := decryption.UnpackSimple(data)

	fileBytes := byteframe.NewByteFrameFromBytes(decrypted)
	fileBytes.SetLE()

	fileBytes.Seek(270, 0)
	questLevel := fileBytes.ReadInt16()

	return questLevel >= s.server.erupeConfig.GameplayOptions.EnhancedInvasions.MinimumInvasionRank
}

func getOriginalArea(s *Session, questId string) int {
	data, _ := os.ReadFile(filepath.Join(s.server.erupeConfig.BinPath, fmt.Sprintf("quests/%s.bin", questId)))
	decrypted := decryption.UnpackSimple(data)

	fileBytes := byteframe.NewByteFrameFromBytes(decrypted)
	fileBytes.SetLE()
	fileBytes.Seek(228, 0)

	area := fileBytes.ReadBytes(1)

	return int(area[0])
}

func overwriteInvasion(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysGetFile)
	area := getOriginalArea(s, pkt.Filename)
	questId := ""

	switch int(area) {
	case 2: // forest and hills
		questId = "26613d0"
		break
	case 16: // forest and hills - night
		questId = "26613n0"
		break
	case 3: // desert
		questId = "26616d0"
		break
	case 17: // desert - night
		questId = "26616n0"
		break
	case 4: // swamp
		questId = "70032d0"
		break
	case 18: // swamp - night
		questId = "70032n0"
		break
	case 5: // volcano
		questId = "70033d0"
		break
	case 19: // volcano - night
		questId = "70033n0"
		break
	case 6: // jungle
		questId = "26619d0"
		break
	case 20: // jungle - night
		questId = "26619n0"
		break
	case 11: // snowy mountains
		questId = "70035d0"
		break
	case 21: // snowy mountains - night
		questId = "70035n0"
		break
	case 26: // great forest
		questId = "26622d0"
		break
	case 27: // great forest - night
		questId = "26622n0"
		break
	case 31: // gorge
		questId = "26625d0"
		break
	case 32: // gorge - night
		questId = "26625n0"
		break
	case 50: // highlands
		questId = "70034d0"
		break
	case 51: // highlands - night
		questId = "70034n0"
		break
	case 57: // tidal islands
		questId = "70031d0"
	case 58: // tidal islands - night
		questId = "70031n0"
		break
	}

	// if we don't have a questId, then we don't have an invasion for this area
	if questId != "" {
		pkt.Filename = questId
		s.spawnInvasion = false // only turn off the invasion if the previous quest supported it.
	}
}
