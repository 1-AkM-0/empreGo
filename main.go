package main

import (
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

	jobChannel := make(chan storage.Job, 10)

	var wg sync.WaitGroup

	sources := []func(jobChannel chan storage.Job) error{
		search.SearchLinkedin,
		search.SearchGupy,
	}

	log.Println("Inicinado busca")
	counter := 0

	for _, search := range sources {
		wg.Go(func() {
			err := search(jobChannel)
			if err != nil {
				log.Println("erro em alguma das fontes", err)
				return
			}
		})
	}

	go func() {
		wg.Wait()
		close(jobChannel)
	}()

	for result := range jobChannel {
		if !(db.AlreadyExists(result.Link)) {

			_, err = bot.SendMessage(channelID, "Nova vaga: "+result.Title+"\n"+result.Link)
			if err != nil {
				log.Println("erro ao tentar enviar vaga pelo bot", err)
			}

			err := db.InsertJob(result)
			if err != nil {
				log.Println(err)
			}

			counter++
		}
	}
	log.Printf("%d nova(s) vaga(s) encontrada(s)", counter)
}
