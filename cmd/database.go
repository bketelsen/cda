// +build !darwin
// +build !windows

package cmd

import (
	"database/sql"
	"log"
	"strings"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

// Database interface
type Database interface {
	Get(shortcode string) (string, error)
	Save(shortcode, url string) (string, error)
}

type sqlite struct {
	Path string
}

func (s sqlite) Save(shortcode, url string) (string, error) {
	if len(strings.TrimSpace(shortcode)) == 0 && len(strings.TrimSpace(url)) == 0 {
		return shortcode, errors.New("Empty shortcode or URL given")
	}
	db, err := sql.Open("sqlite3", s.Path)
	tx, err := db.Begin()
	if err != nil {
		return shortcode, err
	}
	stmt, err := tx.Prepare("insert into urls(shortcode, url) values(?,?)")
	if err != nil {
		return shortcode, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(shortcode, url)
	if err != nil {
		tx.Rollback()
		return shortcode, err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return shortcode, err
	}
	if ra == 0 {
		return shortcode, errors.New("no rows inserted")
	}
	tx.Commit()
	//result
	return shortcode, nil
}

func (s sqlite) Get(shortcode string) (string, error) {
	db, err := sql.Open("sqlite3", s.Path)
	stmt, err := db.Prepare("select url from urls where shortcode = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	var url string
	err = stmt.QueryRow(shortcode).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s sqlite) Init() {
	c, err := sql.Open("sqlite3", s.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	sqlStmt := `create table if not exists urls (shortcode text not null primary key, url text);`
	_, err = c.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}