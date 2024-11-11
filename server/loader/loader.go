package loader

import (
	"encoding/json"
	"feh-map-editor/updater"
	"os"
)

var SkillToId = make(map[string]string)
var HeroToId = make(map[string]string)

var IdToSkill = make(map[string]updater.SkillStruct)
var IdToHero = make(map[string]string)

func LoadJSONs() {
	skillBytes, _ := os.ReadFile("skills_to_ids.json")
	json.Unmarshal(skillBytes, &SkillToId)

	reverseSkillBytes, _ := os.ReadFile("ids_to_skills.json")
	json.Unmarshal(reverseSkillBytes, &IdToSkill)

	heroBytes, _ := os.ReadFile("character_to_id.json")
	json.Unmarshal(heroBytes, &HeroToId)

	reverseHeroBytes, _ := os.ReadFile("id_to_character.json")
	json.Unmarshal(reverseHeroBytes, &IdToHero)
}
