package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsRed    bool   `json:"is_chirpy_red"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp `json:"chirps"`
	Users         map[int]User  `json:"users"`
	RevokedTokens map[string]time.Time
}

func NewDB(path string) (*DB, error) {
	mux := &sync.RWMutex{}
	db := &DBStructure{
		Chirps:        make(map[int]Chirp),
		Users:         make(map[int]User),
		RevokedTokens: make(map[string]time.Time),
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

func (db *DB) CreateChirp(body string, author int) (std Chirp, err error) {
	curr, err := db.loadDB()
	if err != nil {
		return std, err
	}
	nextID := len(curr.Chirps) + 1
	c := Chirp{
		ID:       nextID,
		Body:     body,
		AuthorID: author,
	}
	curr.Chirps[nextID] = c

	db.writeDB(curr)
	return c, nil
}

func (db *DB) GetChirps(id int) ([]Chirp, error) {
	curr, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := []Chirp{}
	for _, c := range curr.Chirps {
		if id > 0 {
			if c.AuthorID == id {
				chirps = append(chirps, c)
				continue
			}
		} else {
			chirps = append(chirps, c)
		}
	}

	return chirps, nil
}

func (db *DB) DeleteChirp(id int) error {
	curr, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(curr.Chirps, id)

	db.writeDB(curr)
	return nil
}

func (db *DB) CreateUser(email string, hash []byte) (std User, err error) {
	curr, err := db.loadDB()
	if err != nil {
		return std, err
	}
	// checking if email exists
	for _, u2 := range curr.Users {
		if u2.Email == email {
			return std, errors.New("duplicated email")
		}
	}

	nextID := len(curr.Users) + 1
	u := User{
		ID:       nextID,
		Email:    email,
		Password: string(hash),
		IsRed:    false,
	}
	curr.Users[nextID] = u

	db.writeDB(curr)
	return u, nil
}

func (db *DB) UpdateUser(id int, email string, hash []byte) (std User, err error) {
	curr, err := db.loadDB()
	if err != nil {
		return std, err
	}

	// checking if user exists
	_, ok := curr.Users[id]
	if !ok {
		return std, errors.New("the user does not exist")
	}
	u := User{
		ID:       id,
		Email:    email,
		Password: string(hash),
	}
	curr.Users[id] = u

	db.writeDB(curr)
	return u, nil
}

func (db *DB) UpdateUserMembership(id int) error {
	curr, err := db.loadDB()
	if err != nil {
		return err
	}

	// checking if user exists
	user, ok := curr.Users[id]
	if !ok {
		return errors.New("the user does not exist")
	}
	user.IsRed = true
	curr.Users[id] = user

	db.writeDB(curr)
	return nil
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

func (db *DB) GetUserByEmail(email string) (us User, err error) {
	curr, err := db.loadDB()
	if err != nil {
		return us, err
	}
	for _, u := range curr.Users {
		if u.Email == email {
			return u, nil
		}
	}

	return us, errors.New("user does not exist")
}

func (db *DB) IsRevoked(token string) bool {
	curr, err := db.loadDB()
	if err != nil {
		return true
	}

	if _, ok := curr.RevokedTokens[token]; !ok {
		return false
	}

	db.writeDB(curr)

	return true
}

func (db *DB) RevokeToken(token string) error {
	curr, err := db.loadDB()
	if err != nil {
		return err
	}

	if _, ok := curr.RevokedTokens[token]; ok {
		return errors.New("it's already revoked")
	}

	curr.RevokedTokens[token] = time.Now()

	db.writeDB(curr)

	return nil
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
