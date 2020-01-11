package model

import (
	"time"

	"github.com/go-gorp/gorp"
)

type Item struct {
	Id         int64     `json:"id" db:"id,primarykey,autoincrement"`
	Name       string    `json:"name" db:"name,notnull,size:200"`
	Text       string    `json:"text" db:"text,size:400"`
	PreUpdated time.Time `json:"preupdated" db:"preupdated,notnull"`
	Created    time.Time `json:"created" db:"created,notnull"`
	Updated    time.Time `json:"updated" db:"updated,notnull"`
}

func (i *Item) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	i.Created = now
	i.Updated = now
	i.PreUpdated = now
	return nil
}

func (i *Item) PreUpdate(s gorp.SqlExecutor) error {
	i.PreUpdated = i.Updated
	i.Updated = time.Now()
	return nil
}
