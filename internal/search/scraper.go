package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Job struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type GupyJobs struct {
	Data []struct {
		Title string `json:"name"`
		Link  string `json:"jobUrl"`
	} `json:"data"`
}

func SearchGupy() ([]Job, error) {
	rawUrl := "https://employability-portal.gupy.io/api/v1/jobs?jobName=est%C3%A1gio&limit=10&offset=0&workplaceType=remote"
	method := "GET"

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(method, rawUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("erro na tentativa de fazer o wrapper do request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer o request: %v", err)
	}
	defer res.Body.Close()
	gupyResponse := &GupyJobs{}
	err = json.NewDecoder(res.Body).Decode(&gupyResponse)
	if err != nil {
		return nil, fmt.Errorf("erro ao tentar decodar o json: %v", err)
	}
	jobs := []Job{}

	for _, result := range gupyResponse.Data {
		jobToInsert := Job{
			Title: result.Title,
			Link:  result.Link,
		}
		jobs = append(jobs, jobToInsert)
	}

	return jobs, nil
}

func SearchLinkedin() ([]Job, error) {
	rawUrl := "https://www.linkedin.com/jobs/search?keywords=%22est%C3%A1gio%22+%28%22ti%22+OR+%22desenvolvimento%22+OR+%22web%22+OR+%22backend%22+OR+%22fullstack%22%29&location=Brasil&geoId=106057199&f_TPR=r86400&f_WT=2"
	method := "GET"

	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest(method, rawUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("erro na tentativa de fazer o wrapper do request: %v", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, Like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer o request: %v", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("linkedin nao retornou 200")
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao tentar parsear o html: %v", err)
	}

	var jobs []Job

	doc.Find("ul.jobs-search__results-list > li").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h3.base-search-card__title").Text())
		link, exists := (s.Find("a.base-card__full-link").Attr("href"))
		u, err := url.Parse(link)
		if err != nil {
			return
		}
		u.RawQuery = ""
		link = u.String()

		if title != "" && exists {
			job := Job{
				Title: title,
				Link:  link,
			}
			jobs = append(jobs, job)
		}

	})

	return jobs, nil

}
