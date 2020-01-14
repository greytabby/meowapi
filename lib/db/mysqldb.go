package db

import (
	"database/sql"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/greytabby/meowapi/lib/model"
)

// MysqlDbAccessor mysqlへの Accessor
// Implementation handler's XXDbAccessor interface
type MysqlDbAccessor struct {
	Db *gorp.DbMap
}

// NewMysqlDbAccessor MysqlDbAccessorを返す
func NewMysqlDbAccessor(dsn string) (*MysqlDbAccessor, error) {
	dbmap, err := newDbMap(dsn)
	if err != nil {
		return nil, err
	}
	return &MysqlDbAccessor{Db: dbmap}, nil
}

// GetAllCats DBからcatテーブルの全てのデータを取得する
func (mda *MysqlDbAccessor) GetAllCats() ([]model.Cat, error) {
	var cats []model.Cat
	_, err := mda.Db.Select(&cats,
		"SELECT * FROM cat ORDER BY created")
	if err != nil {
		return nil, err
	}
	return cats, nil
}

// GetCat DBのcatテーブルからidに合致するcatを1つ返す
// 見つからなかった場合は空のcatとerrorを返す
func (mda *MysqlDbAccessor) GetCat(id int64) (model.Cat, error) {
	var cat model.Cat
	err := mda.Db.SelectOne(&cat, "SELECT * FROM cat WHERE id = ?", id)
	if err != nil {
		return model.Cat{}, err
	}
	return cat, nil
}

// AddCat catテーブルへデータを1件追加する
func (mda *MysqlDbAccessor) AddCat(cat model.Cat) error {
	err := mda.Db.Insert(&cat)
	if err != nil {
		return err
	}
	return nil
}

// UpdateCat catテーブルのデータを1件更新する
func (mda *MysqlDbAccessor) UpdateCat(cat model.Cat) error {
	_, err := mda.Db.Update(&cat)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCat catテーブルのデータを1件削除する
func (mda *MysqlDbAccessor) DeleteCat(cat model.Cat) error {
	_, err := mda.Db.Delete(&cat)
	if err != nil {
		return err
	}
	return nil
}

func newDbMap(dsn string) (*gorp.DbMap, error) {
	db, err := open(dsn)
	if err != nil {
		return nil, err
	}
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Engine: "InnoDB", Encoding: "UTF8"}}
	return dbmap, nil
}

func open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}