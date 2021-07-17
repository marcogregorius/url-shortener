package models

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/lib/pq"
	"github.com/lithammer/shortuuid"
)

const (
	// with shortuuid, we can effectively have 57 ^ 8 = 111e12 unique URLs
	Length = 8
	Table  = "tb_shortlinks"
)

// use db as global variable, so we don't have to Close() after every query
var db *sql.DB
var err error

type Formatter interface {
	Format() error
}

type Shortlink struct {
	// Id is not validated as it is generated on method CreateShortlink()
	Id            string      `json:"id"`
	SourceUrl     string      `json:"source_url" validate:"required,url"`
	Visited       int         `json:"visited"`
	LastVisitedAt pq.NullTime `json:"last_visited_at"`
	ExpiredAt     pq.NullTime `json:"expired_at"`
	CreatedAt     pq.NullTime `json:"created_at"`
}

func (s *Shortlink) Format() (map[string]string, error) {
	fields := reflect.TypeOf(s)
	values := reflect.ValueOf(s)
	fmt.Println(fields, values)
	out := map[string]string{"a": "a"}
	return out, nil
	//if s.LastVisitedAt.Valid {
	//s.LastVisitedAt = s.LastVisitedAt.Time
	//}
	//if s.ExpiredAt.Valid {
	//s.ExpiredAt = s.ExpiredAt.Time
	//}
	//if s.CreatedAt.Valid {
	//s.CreatedAt = s.CreatedAt.Time
	//}
}

//type NullTime struct {
//Time time.Time
//}

//func (nt *NullTime) Scan(value interface{}) error {
//nt.Time = value.(time.Time)
//return nil
//}

//func (nt NullTime) Value() (driver.Value, error) {
//if !nt.Valid {
//return nil, nil
//}
//return nt.Time, nil
//}

func InitDb() *sql.DB {
	db, err = sql.Open("postgres", "postgres://localhost:5432/url-shortener?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (s *Shortlink) CreateShortlink() (err error) {
	s.Id = shortuuid.New()[Length:]
	_, err = db.Exec("INSERT INTO tb_shortlinks(id,source_url,expired_at,created_at) VALUES($1, $2, $3, $4)",
		s.Id, s.SourceUrl, s.ExpiredAt, time.Now())
	if err == nil {
		s.Format()
	}
	return
}

func (s *Shortlink) GetShortlink() (err error) {
	err = db.QueryRow("SELECT source_url,visited,last_visited_at,expired_at,created_at FROM tb_shortlinks WHERE id=$1", s.Id).
		Scan(&s.SourceUrl,
			&s.Visited,
			&s.LastVisitedAt,
			&s.ExpiredAt,
			&s.CreatedAt)
	if err == nil {
		s.Format()
	}
	return
}

func (s *Shortlink) UpdateShortlink() (err error) {
	_, err = db.Exec("UPDATE tb_shortlinks SET visited = $1, last_visited_at = $2 WHERE id = $3",
		s.Visited, s.LastVisitedAt, s.Id)
	if err == nil {
		s.Format()
	}
	return
}
