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
	Name string `json:"name"`
	Id   string `json:"id"`
}

type ResponseCharacter struct {
	CargoQuery []struct {
		Title Character `json:"title"`
	} `json:"cargoquery"`
}

func fetchCharacters() (map[string]string, map[string]string) {
	var id_to_char = make(map[string]string)
	var char_to_id = make(map[string]string)
	var query = url.Values{}
	query.Add("action", "cargoquery")
	query.Add("format", "json")
	query.Add("tables", "Units")
	query.Add("fields", "_pageName=name, TagID=id")
	query.Add("limit", "500")
	query.Add("offset", "0")

	for {
		req, _ := http.Get("https://feheroes.fandom.com/api.php?" + query.Encode())
		byteResponse, _ := io.ReadAll(req.Body)
		var responseStruct ResponseCharacter = ResponseCharacter{}
		json.Unmarshal(byteResponse, &responseStruct)
		for _, el := range responseStruct.CargoQuery {
			id_to_char[el.Title.Id] = el.Title.Name
			char_to_id[el.Title.Name] = el.Title.Id
		}
		if len(responseStruct.CargoQuery) == 500 {
			var intOffset, _ = strconv.Atoi(query.Get("offset"))
			query.Set("offset", strconv.Itoa(intOffset+500))
		} else {
			break
		}
	}

	return id_to_char, char_to_id
}

func fetchSkills() (map[string]string, map[string]string) {
	var id_to_skill = make(map[string]string)
	var skill_to_id = make(map[string]string)
	var query = url.Values{}
	query.Add("action", "cargoquery")
	query.Add("format", "json")
	query.Add("tables", "Skills")
	query.Add("fields", "Name=name, TagID=id")
	query.Add("where", "RefinePath is null")
	query.Add("limit", "500")
	query.Add("offset", "0")

	for {
		req, _ := http.Get("https://feheroes.fandom.com/api.php?" + query.Encode())
		byteResponse, _ := io.ReadAll(req.Body)
		var responseStruct ResponseCharacter = ResponseCharacter{}
		json.Unmarshal(byteResponse, &responseStruct)
		for _, el := range responseStruct.CargoQuery {
			id_to_skill[el.Title.Id] = el.Title.Name
			skill_to_id[el.Title.Name] = el.Title.Id
		}
		if len(responseStruct.CargoQuery) == 500 {
			var intOffset, _ = strconv.Atoi(query.Get("offset"))
			query.Set("offset", strconv.Itoa(intOffset+500))
		} else {
			break
		}
	}

	return id_to_skill, skill_to_id
}

func Update() {
	var id_to_char, char_to_id = fetchCharacters()
	var byteChars, _ = json.Marshal(char_to_id)
	var byteIds, _ = json.Marshal(id_to_char)
	os.WriteFile("id_to_character.json", byteIds, 0644)
	os.WriteFile("character_to_id.json", byteChars, 0644)

	var id_to_skill, skill_to_id = fetchSkills()
	var byteSkills, _ = json.Marshal(skill_to_id)
	os.WriteFile("skills_to_ids.json", byteSkills, 0644)

	var byteIdSkills, _ = json.Marshal(id_to_skill)
	os.WriteFile("ids_to_skills.json", byteIdSkills, 0644)
}
