package storage

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type Job struct {
	ID    int
	Title string
	Link  string
}

type SQLiteJobsStore struct {
	db *sql.DB
}

func NewSQLite() (*SQLiteJobsStore, error) {
	db, err := sql.Open("sqlite", "vagas.db")
	if err != nil {
		return nil, fmt.Errorf("erro ao tentar abrir conexao com o database: %w", err)
	}
	return &SQLiteJobsStore{db: db}, nil

}

func (s *SQLiteJobsStore) Close() {
	s.db.Close()
}

func (s *SQLiteJobsStore) InsertJob(job Job) error {

	query := `
	INSERT INTO jobs (title, link)
	VALUES (? , ?)
	`
	_, err := s.db.Exec(query, job.Title, job.Link)
	if err != nil {
		return fmt.Errorf("erro ao tentar inserir job no database: %w", err)
	}

	return nil
}

func (s *SQLiteJobsStore) GetJobs() ([]Job, error) {
	query := `
	SELECT title, link FROM jobs
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer consulta na tabela jobs: %w", err)
	}
	defer rows.Close()

	jobs := []Job{}

	for rows.Next() {
		job := Job{}
		err := rows.Scan(&job.Title, &job.Link)
		if err != nil {
			return nil, fmt.Errorf("erro ao tentar scanear os valores: %w", err)
		}
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (s *SQLiteJobsStore) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS jobs(
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"link" TEXT UNIQUE
	);
	`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("erro ao tentar criar a tabela jobs: %w", err)
	}
	fmt.Println("tabela jobs criada")
	return nil
}

func (s *SQLiteJobsStore) AlreadyExists(link string) bool {
	var lin string
	query := `
	SELECT link FROM jobs
	WHERE link = (?)
	
	`
	rows, err := s.db.Query(query, link)
	if err != nil {
		fmt.Println("erro na pesquisa")
		return true
	}
	defer rows.Close()

	for rows.Next() {

		rows.Scan(&lin)

		fmt.Println(lin)
		return true
	}

	return false

}
