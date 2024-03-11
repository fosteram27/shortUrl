package urls

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DBStore struct {
	db *sql.DB
}

func NewDBStore(filename string) (*DBStore, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	sqlStmt := `
	create table if not exists urlEntries (hash varchar(16) not null primary key, urlLong text, 
	urlShort text, date text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return &DBStore{
		db: db,
	}, nil
}

func (db *DBStore) Close() error {
	return db.db.Close()
}

func (db *DBStore) Add(hash string, urlEntry UrlEntry) error {
	_, err := db.db.Exec("insert into urlEntries(hash, urlLong, urlShort) values(?, ?, ?)", hash, urlEntry.UrlLong, urlEntry.UrlShort)
	return err
}

func (db *DBStore) Get(name string) (UrlEntry, error) {
	// db.db.QueryRow()
	return UrlEntry{}, nil
}

func (db *DBStore) List() (map[string]UrlEntry, error) {
	rows, err := db.db.Query("select hash, urlLong, urlShort from urlEntries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urlEntries := make(map[string]UrlEntry)
	for rows.Next() {
		var urlEntry UrlEntry
		err = rows.Scan(&urlEntry.Id, &urlEntry.UrlLong, &urlEntry.UrlShort)
		if err != nil {
			return nil, err
		}
		urlEntries[urlEntry.Id] = urlEntry
	}
	return urlEntries, rows.Err()
}

func (db *DBStore) Remove(name string) error {
	return nil
}
