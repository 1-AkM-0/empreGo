package main

import (
	"log"
	"os"

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

	results, err := search.SearchJobs()
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		if !(db.AlreadyExists(result.Link)) {

			jobToInsert := storage.Job{
				Title: result.Title,
				Link:  result.Link,
			}

			err := db.InsertJob(jobToInsert)
			if err != nil {
				log.Println(err)
			}

			bot.SendMessage(channelID, "Nova vaga: "+result.Title+"\n"+result.Link)
		}
	}
}
