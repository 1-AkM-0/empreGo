package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/1-AkM-0/empreGo/internal/discord"
	"github.com/1-AkM-0/empreGo/internal/search"
	"github.com/1-AkM-0/empreGo/internal/storage"
)

func main() {
	db, err := storage.NewSQLite()
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := discord.NewBot(os.Getenv("BOT_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	channelID := os.Getenv("CHANNEL_ID")

	defer db.Close()
	defer bot.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex

	sources := []func() ([]search.Job, error){
		search.SearchLinkedin,
		search.SearchGupy,
	}

	allJobs := []search.Job{}

	for _, search := range sources {
		wg.Go(func() {
			jobs, err := search()
			if err != nil {
				log.Println("erro em alguma das fontes", err)
				return
			}
			mu.Lock()
			allJobs = append(allJobs, jobs...)
			mu.Unlock()
		})
	}

	wg.Wait()

	for _, result := range allJobs {
		if !(db.AlreadyExists(result.Link)) {

			jobToInsert := storage.Job{
				Title: result.Title,
				Link:  result.Link,
			}

			err := db.InsertJob(jobToInsert)
			if err != nil {
				log.Println(err)
			}

			_, err = bot.SendMessage(channelID, "Nova vaga: "+result.Title+"\n"+result.Link)
			if err != nil {
				log.Println("erro ao tentar enviar vaga pelo bot", err)
			}
		}
	}
}
