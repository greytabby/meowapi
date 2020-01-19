package model

import (
	"time"

	"github.com/go-gorp/gorp"
)

type User struct {
	Id       int64     `json:"id"       db:"id,primarykey,autoincrement"`
	Name     string    `json:"name"     db:"name,notnull,size:200"`
	Password string    `json:"password" db:"password,notnull,size:400"`
	Created  time.Time `json:"created"  db:"created,notnull"`
	Updated  time.Time `json:"updated"  db:"updated,notnull"`
}

func (u *User) PreInsert(s gorp.SqlExecutor) error {
	u.Created = time.Now()
	u.Updated = u.Created
	return nil
}

func (u *User) PreUpdate(s gorp.SqlExecutor) error {
	u.Updated = time.Now()
	return nil
}
