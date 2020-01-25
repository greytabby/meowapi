package model

import (
	"time"

	"github.com/go-gorp/gorp"
)

type UseToilet struct {
	Id       int64     `json:"id"       db:"id,primarykey,autoincrement"`
	UID      int64     `json:"uid"      db:"uid,notnull"`
	ToiletId int64     `json:"toiletid" db:"toiletid,notnull"`
	CatId    int64     `json:"catid"    db:"catid,notnull"`
	Type     string    `json:"type"     db:"type,notnull,size:200"`
	Created  time.Time `json:"created"  db:"created,notnull"`
	Updated  time.Time `json:"updated"  db:"updated,notnull"`
}

func (ut *UseToilet) PreInsert(s gorp.SqlExecutor) error {
	now := time.Now()
	ut.Created = now
	ut.Updated = now
	return nil
}

func (ut *UseToilet) PreUpdate(s gorp.SqlExecutor) error {
	ut.Updated = time.Now()
	return nil
}
