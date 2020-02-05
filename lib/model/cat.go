package model

import (
	"time"

	"github.com/go-gorp/gorp"
)

type Cat struct {
	Id      int64     `json:"id"      db:"id,primarykey,autoincrement"`
	UID     int64     `json:"uid"     db:"uid,notnull"`
	Name    string    `json:"name"    db:"name,notnull,size:200"`
	Breed   string    `json:"breed"   db:"breed,size:200"`
	Gender  string    `json:"gender"  db:"gender,size:200"`
	Age     int64     `json:"age"     db:"age"`
	Created time.Time `json:"created" db:"created,notnull"`
	Updated time.Time `json:"updated" db:"updated,notnull"`
}

func (c *Cat) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	c.Created = now
	c.Updated = now
	return nil
}

func (c *Cat) PreUpdate(s gorp.SqlExecutor) error {
	c.Updated = time.Now()
	return nil
}
