package model

import (
	"time"

	"github.com/go-gorp/gorp"
)

type Toilet struct {
	Id        int64     `json:"id"              db:"id,primarykey,autoincrement"`
	Name      string    `json:"name"            db:"name,notnull,size:200"`
	Comment   string    `json:"comment"         db:"comment,size:400"`
	SandState string    `json:"sandstate"       db:"sandstate,size:50"`
	Created   time.Time `json:"created"         db:"created,notnull"`
	Updated   time.Time `json:"updated"         db:"updated,notnull"`
}

func (t *Toilet) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	t.Created = now
	t.Updated = now
	return nil
}

func (t *Toilet) PreUpdate(s gorp.SqlExecutor) error {
	t.Updated = time.Now()
	return nil
}
