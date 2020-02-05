package model

import (
	"time"

	"github.com/go-gorp/gorp"
)

type Wash struct {
	Id       int64     `json:"id"       db:"id,primarykey,autoincrement"`
	UID      int64     `json:"uid"      db:"uid,notnull"`
	ToiletId int64     `json:"toiletid" db:"toiletid,notnull"`
	Comment  string    `json:"comment"  db:"comment,size:400"`
	Created  time.Time `json:"created"  db:"created,notnull"`
	Updated  time.Time `json:"updated"  db:"updated,notnull"`
}

func (w *Wash) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	w.Created = now
	w.Updated = now
	return nil
}

func (w *Wash) PreUpdate(s gorp.SqlExecutor) error {
	w.Updated = time.Now()
	return nil
}
