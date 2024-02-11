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

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

func NewDB(path string) (*DB, error) {
	mux := &sync.RWMutex{}
	db := &DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]User),
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	jsonEncoder := json.NewEncoder(f)
	err = jsonEncoder.Encode(db)
	if err != nil {
		return nil, err
	}

	return &DB{
		path: path,
		mux:  mux,
	}, nil
}

func (db *DB) CreateChirp(body string) (std Chirp, err error) {
	curr, err := db.loadDB()
	if err != nil {
		return std, err
	}
	nextID := len(curr.Chirps) + 1
	c := Chirp{
		ID:   nextID,
		Body: body,
	}
	curr.Chirps[nextID] = c

	db.writeDB(curr)
	return c, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	curr, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := []Chirp{}
	for _, c := range curr.Chirps {
		chirps = append(chirps, c)
	}

	return chirps, nil
}

func (db *DB) CreateUser(email string) (std User, err error) {
	curr, err := db.loadDB()
	if err != nil {
		return std, err
	}
	nextID := len(curr.Users) + 1
	u := User{
		ID:    nextID,
		Email: email,
	}
	curr.Users[nextID] = u

	db.writeDB(curr)
	return u, nil
}

func (db *DB) GetUsers() ([]User, error) {
	curr, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	users := []User{}
	for _, u := range curr.Users {
		users = append(users, u)
	}

	return users, nil
}

func (db *DB) ensureDB() error {
	d, err := os.Open(db.path)
	if err != nil {
		return err
	}

	defer d.Close()
	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(db.path)
		if err != nil {
			return err
		}
		defer f.Close()
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
	defer db.mux.Unlock()
	err := db.ensureDB()
	if err != nil {
		return err
	}

	res, err := json.Marshal(dbStructure)
	err = os.WriteFile(db.path, res, 0666)
	if err != nil {
		return err
	}
	return nil
}
