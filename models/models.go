package models

import (
	"database/sql"
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

type Shortlink struct {
	// Id is not validated as it is generated on method CreateShortlink()
	Id            string      `json:"id"`
	SourceUrl     string      `json:"source_url" validate:"required,url"`
	Visited       int         `json:"visited"`
	LastVisitedAt pq.NullTime `json:"last_visited_at"`
	CreatedAt     pq.NullTime `json:"created_at"`
}

func (s *Shortlink) formatForOutput() map[string]interface{} {
	// returns the JSON representation in map, with fields following the json tag
	values := reflect.ValueOf(s).Elem()
	typeOfS := values.Type()
	out := map[string]interface{}{}
	for i := 0; i < values.NumField(); i++ {
		key := typeOfS.Field(i).Tag.Get("json")
		value := values.Field(i).Interface()

		// Special handling for pq.NullTime
		if res, ok := value.(pq.NullTime); ok {
			if res.Valid {
				value = res.Time
			} else {
				value = nil
			}
		}
		out[key] = value
	}
	return out
}

func (s *Shortlink) CreateShortlink(db *sql.DB) (map[string]interface{}, error) {
	s.Id = shortuuid.New()[Length:]
	s.CreatedAt.Time, s.CreatedAt.Valid = time.Now(), true
	_, err := db.Exec("INSERT INTO tb_shortlinks(id,source_url,created_at) VALUES($1, $2, $3)",
		s.Id, s.SourceUrl, time.Now())
	if err != nil {
		return nil, err
	}
	return s.formatForOutput(), nil
}

func (s *Shortlink) GetShortlink(db *sql.DB) (map[string]interface{}, error) {
	err := db.QueryRow("SELECT source_url,visited,last_visited_at,created_at FROM tb_shortlinks WHERE id=$1", s.Id).
		Scan(&s.SourceUrl,
			&s.Visited,
			&s.LastVisitedAt,
			&s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s.formatForOutput(), nil
}

func (s *Shortlink) VisitShortlink(db *sql.DB) error {
	_, err := db.Exec("UPDATE tb_shortlinks SET visited = visited + 1, last_visited_at = $1 WHERE id = $2",
		s.LastVisitedAt, s.Id)
	if err != nil {
		return err
	}
	return nil
}
