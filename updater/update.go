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

func fetchCharacters() []Character {
	var query = url.Values{}
	var data []Character = []Character{}
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
			data = append(data, el.Title)
		}
		if len(responseStruct.CargoQuery) == 500 {
			var intOffset, _ = strconv.Atoi(query.Get("offset"))
			query.Set("offset", strconv.Itoa(intOffset+500))
		} else {
			break
		}
	}

	return data
}

func fetchSkills() []Character {
	var query = url.Values{}
	var data []Character = []Character{}
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
			data = append(data, el.Title)
		}
		if len(responseStruct.CargoQuery) == 500 {
			var intOffset, _ = strconv.Atoi(query.Get("offset"))
			query.Set("offset", strconv.Itoa(intOffset+500))
		} else {
			break
		}
	}

	return data
}

func Update() {
	var chars = fetchCharacters()
	var byteChars, _ = json.Marshal(chars)
	os.WriteFile("characters.json", byteChars, 0644)

	var skills = fetchSkills()
	var byteSkills, _ = json.Marshal(skills)
	os.WriteFile("skills.json", byteSkills, 0644)

}
