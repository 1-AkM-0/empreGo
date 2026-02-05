package main

import (
	"fmt"
	"log"

	"github.com/1-AkM-0/empreGo/internal/search"
	"github.com/1-AkM-0/empreGo/internal/storage"
)

func main() {
	db, err := storage.NewSQLite()
	db.CreateTable()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	results, err := search.SearchJobs()
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		db.InsertJob(result)
	}
	jobs, err := db.GetJobs()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jobs)

}
