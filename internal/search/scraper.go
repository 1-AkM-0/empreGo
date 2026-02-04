package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SerperResponse struct {
	Jobs []Jobs `json:"organic"`
}

type Jobs struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func SearchJobs() (*SerperResponse, error) {
	url := fmt.Sprint("https://google.serper.dev/search?q=(site%3Abr.linkedin.com%2Fjobs%2Fview+OR+site%3Agupy.io)+%22est%C3%A1gio%22+(%22desenvolvimento%22+OR+%22backend%22+OR+%22ti%22+OR+%22web%22+OR+%22fullstack%22)&gl=br&hl=pt-br&tbs=qdr%3Aw&apiKey=" + os.Getenv("API_KEY"))
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("erro na tentativa de fazer o wrapper do request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer o request: %v", err)

	}
	defer res.Body.Close()
	serperResponse := &SerperResponse{}

	_ = json.NewDecoder(res.Body).Decode(&serperResponse)

	return serperResponse, nil
}
