package updater

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Character struct {
	Name     string `json:"name"`
	Id       string `json:"id"`
	WikiName string `json:"wikiName,omitempty"`
	Slot     string `json:"slot,omitempty"`
}

type SkillStruct struct {
	Name string `json:"name"`
	Slot string `json:"slot"`
}

type ResponseCharacter struct {
	CargoQuery []struct {
		Title Character `json:"title"`
	} `json:"cargoquery"`
}

type CharacterData struct {
	Id       string `json:"id"`
	WikiName string `json:"wikiName"`
}

type LearnsetData struct {
	Page  string `json:"Page"`
	Skill string `json:"Skill"`
}

type LearnsetResponse struct {
	CargoQuery []struct {
		Title LearnsetData `json:"title"`
	} `json:"cargoquery"`
}

func fetchCharacters(channel chan bool, id_to_char *map[string]string, char_to_id *map[string]CharacterData) {
	var query = url.Values{}
	query.Add("action", "cargoquery")
	query.Add("format", "json")
	query.Add("tables", "Units")
	query.Add("fields", "_pageName=name, TagID=id, WikiName=wikiName")
	query.Add("limit", "500")
	query.Add("offset", "0")

	for {
		req, _ := http.Get("https://feheroes.fandom.com/api.php?" + query.Encode())
		byteResponse, _ := io.ReadAll(req.Body)
		var responseStruct ResponseCharacter = ResponseCharacter{}
		json.Unmarshal(byteResponse, &responseStruct)
		for _, el := range responseStruct.CargoQuery {
			var charData CharacterData = CharacterData{
				Id:       el.Title.Id,
				WikiName: el.Title.WikiName,
			}
			(*id_to_char)[el.Title.Id] = el.Title.Name
			(*char_to_id)[el.Title.Name] = charData
		}
		if len(responseStruct.CargoQuery) == 500 {
			var intOffset, _ = strconv.Atoi(query.Get("offset"))
			query.Set("offset", strconv.Itoa(intOffset+500))
		} else {
			break
		}
	}

	channel <- true
}

func fetchSkills(channel chan bool, skill_to_id *map[string]SkillStruct, id_to_skill *map[string]SkillStruct) {
	var query = url.Values{}
	query.Add("action", "cargoquery")
	query.Add("format", "json")
	query.Add("tables", "Skills")
	query.Add("fields", "Name=name, TagID=id, Scategory=slot")
	query.Add("where", "RefinePath is null")
	query.Add("limit", "500")
	query.Add("offset", "0")

	for {
		req, _ := http.Get("https://feheroes.fandom.com/api.php?" + query.Encode())
		byteResponse, _ := io.ReadAll(req.Body)
		var responseStruct ResponseCharacter = ResponseCharacter{}
		json.Unmarshal(byteResponse, &responseStruct)
		for _, el := range responseStruct.CargoQuery {
			(*id_to_skill)[el.Title.Id] = SkillStruct{
				Name: el.Title.Name,
				Slot: el.Title.Slot,
			}
			(*skill_to_id)[el.Title.Name] = SkillStruct{
				Name: el.Title.Id,
				Slot: el.Title.Slot,
			}
		}
		if len(responseStruct.CargoQuery) == 500 {
			var intOffset, _ = strconv.Atoi(query.Get("offset"))
			query.Set("offset", strconv.Itoa(intOffset+500))
		} else {
			break
		}
	}

	channel <- true
}

func fetchLearnsets(channel chan bool, char_to_skills *map[string][]string) {
	var query = url.Values{}
	query.Add("action", "cargoquery")
	query.Add("format", "json")
	query.Add("tables", "UnitSkills, Skills, Units")
	query.Add("fields", "UnitSkills._pageName=Page, Skills.Name=Skill")
	query.Add("join_on", "UnitSkills.skill = Skills.WikiName, UnitSkills._pageName = Units._pageName")
	query.Add("where", "RefinePath is null")
	query.Add("order_by", "ReleaseDate ASC, Exclusive DESC, Skills.Name ASC")
	query.Add("limit", "500")
	query.Add("offset", "0")

	for {
		req, _ := http.Get("https://feheroes.fandom.com/api.php?" + query.Encode())
		byteResponse, _ := io.ReadAll(req.Body)
		var responseStruct LearnsetResponse = LearnsetResponse{}
		json.Unmarshal(byteResponse, &responseStruct)
		for _, el := range responseStruct.CargoQuery {
			_, ok := (*char_to_skills)[el.Title.Page]
			if !ok {
				(*char_to_skills)[el.Title.Page] = []string{}
			}
			(*char_to_skills)[el.Title.Page] = append((*char_to_skills)[el.Title.Page], el.Title.Skill)
		}

		if len(responseStruct.CargoQuery) == 500 {
			var intOffset, _ = strconv.Atoi(query.Get("offset"))
			query.Set("offset", strconv.Itoa(intOffset+500))
		} else {
			break
		}
	}

	channel <- true
}

func Update() {
	var charactersChannel = make(chan bool, 1)
	var skillsChannel = make(chan bool, 1)
	var quit = make(chan bool, 1)
	var learnsetsChannel = make(chan bool, 1)

	var id_to_char = make(map[string]string)
	var char_to_id = make(map[string]CharacterData)

	var id_to_skill = make(map[string]SkillStruct)
	var skill_to_id = make(map[string]SkillStruct)

	var char_to_skills = make(map[string][]string)

	go fetchCharacters(charactersChannel, &id_to_char, &char_to_id)
	go fetchSkills(skillsChannel, &skill_to_id, &id_to_skill)
	go fetchLearnsets(learnsetsChannel, &char_to_skills)

	var finishedChannels = 0
	for {
		select {
		case <-charactersChannel:
			var byteChars, _ = json.Marshal(char_to_id)
			var byteIds, _ = json.Marshal(id_to_char)
			os.WriteFile("id_to_character.json", byteIds, 0644)
			os.WriteFile("character_to_id.json", byteChars, 0644)
			var jsUnitsVariable = []byte("const UNITS = ")
			jsUnitsVariable = append(jsUnitsVariable, byteChars...)

			os.WriteFile("../units.js", jsUnitsVariable, 0644)
			finishedChannels++
			if finishedChannels == 3 {
				quit <- true
			}
		case <-skillsChannel:
			var byteSkills, _ = json.Marshal(skill_to_id)
			var jsSkillsVariable = []byte("const SKILLS = ")
			jsSkillsVariable = append(jsSkillsVariable, byteSkills...)
			os.WriteFile("skills_to_ids.json", byteSkills, 0644)
			os.WriteFile("../skills.js", jsSkillsVariable, 0644)

			var byteIdSkills, _ = json.Marshal(id_to_skill)
			os.WriteFile("ids_to_skills.json", byteIdSkills, 0644)
			finishedChannels++
			if finishedChannels == 3 {
				quit <- true
			}

		case <-learnsetsChannel:
			var byteLearnsets, _ = json.Marshal(char_to_skills)
			var jsLearnsets = []byte("const LEARNSETS = ")
			jsLearnsets = append(jsLearnsets, byteLearnsets...)
			os.WriteFile("../learnsets.js", jsLearnsets, 0644)
			finishedChannels++
			if finishedChannels == 3 {
				quit <- true
			}
		case <-quit:
			return
		}

	}
}
