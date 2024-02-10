package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	mux := &sync.RWMutex{}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return &DB{
		path: path,
		mux:  mux,
	}, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {

}

func (db *DB) GetChirps() ([]Chirp, error) {

}

func (db *DB) ensureDB() error {
	d, err := os.Open(db.path)
	defer d.Close()

	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(db.path)
		defer f.Close()

		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (db *DB) loadDB() (str DBStructure, err error) {
	db.mux.RLock()
	err = db.ensureDB()
	if err != nil {
		return str, err
	}

	f, err := os.ReadFile(db.path)
	var data DBStructure
	err = json.Unmarshal(f, &data)
	db.mux.RUnlock()
	if err != nil {
		return str, err
	}
	return data, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	err := db.ensureDB()
	if err != nil {
		return err
	}

	f, err := os.ReadFile(db.path)
	var data DBStructure
	err = json.Unmarshal(f, &data)
	if err != nil {
		return str, err
	}
	return data, nil
}
