package storage

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite3Storage struct {
	db *sql.DB
}

func (s *Sqlite3Storage) Initialize(args ...interface{}) error {
	// Try open and connect to database
	var err error
	var dbName string
	if len(args) == 0 {
		dbName = "database.db"
	} else {
		dbName = args[0].(string)
	}
	s.db, err = sql.Open("sqlite3", dbName)
	return err
}

func (s *Sqlite3Storage) CreateUser(name string, group string, publicKey string) error {
	// Create a new user
	_, err := s.db.Exec("INSERT INTO voters (name, grouptype, public_key) VALUES (?, ?, ?)", name, group, publicKey)
	return err
}

func (s *Sqlite3Storage) FetchUser(name string) (User, error) {
	// Fetch a user
	var user User
	row := s.db.QueryRow("SELECT name, grouptype, public_key FROM voters WHERE name = ?", name)
	err := row.Scan(&user.Name, &user.Group, &user.PublicKey)
	return user, err
}

func (s *Sqlite3Storage) RemoveUser(name string) error {
	// Remove a user
	_, err := s.db.Exec("DELETE FROM voters WHERE name = ?", name)
	return err
}

func (s *Sqlite3Storage) CreateElection(name string, groups []string, choices []string, endDate time.Time) error {
	// Todo
	return nil
}

func (s *Sqlite3Storage) FetchElection(name string) (Election, error) {
	// Todo
	return Election{}, nil
}

func (s *Sqlite3Storage) VoteElection(electionName string, voterName string, choice string) error {
	// Todo
	return nil
}

func (s *Sqlite3Storage) FetchElectionResults(electionName string) (ElectionResults, error) {
	// Todo
	return ElectionResults{}, nil
}
