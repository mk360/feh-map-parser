package lookup

import (
	"io"
	"net/http"
	"net/url"
)

func LookupSkill(skill string) string {
	var query = url.Values{}
	query.Add("format", "json")
	query.Add("action", "cargoquery")
	query.Add("tables", "Skills")
	query.Add("fields", "Name")
	query.Add("where", "TagID = \""+skill+"\"")
	var url = "https://feheroes.fandom.com/api.php?"
	resp, _ := http.Get(url + query.Encode())
	byteResponse, _ := io.ReadAll(resp.Body)
	return string(byteResponse)
}
