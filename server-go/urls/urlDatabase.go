package urls

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrNotFound2 = errors.New("not found")
)

type urlDb struct {
	db *sql.DB
}

func NewUrlDb() *urlDb {

	db, _ := sql.Open("sqlite3", "./urls.db")

	return &urlDb{
		db,
	}
}

func (m MemStore) Add(name string, urlEntry UrlEntry) error {
	m.list[name] = urlEntry
	return nil
}

func (m MemStore) Get(name string) (UrlEntry, error) {

	if val, ok := m.list[name]; ok {
		return val, nil
	}

	return UrlEntry{}, ErrNotFound
}

func (m MemStore) List() (map[string]UrlEntry, error) {
	return m.list, nil
}

func (m MemStore) Update(name string, urlEntry UrlEntry) error {

	if _, ok := m.list[name]; ok {
		m.list[name] = urlEntry
		return nil
	}

	return ErrNotFound
}

func (m MemStore) Remove(name string) error {
	delete(m.list, name)
	return nil
}
