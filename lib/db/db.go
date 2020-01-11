package db

import (
	"database/sql"

	"github.com/greytabby/meowapi/lib/model"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
)

// MysqlDbAccessor mysqlへの Accessor
// handlerのDBAccessorインタフェースを実装している。
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

// GetAllItems DBからitemsテーブルの全てのデータを取得する
func (mda *MysqlDbAccessor) GetAllItems() ([]model.Item, error) {
	var items []model.Item
	_, err := mda.Db.Select(&items,
		"SELECT * FROM items ORDER BY created")
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetItem DBのitemsテーブルからidに合致するitemを1つ返す
// 見つからなかった場合は空のitemとerrorを返す
func (mda *MysqlDbAccessor) GetItem(id int64) (model.Item, error) {
	var item model.Item
	err := mda.Db.SelectOne(&item, "SELECT * FROM items WHERE id = ?", id)
	if err != nil {
		return model.Item{}, err
	}
	return item, nil
}

// InsertItem itemsテーブルへデータを1件追加する
func (mda *MysqlDbAccessor) InsertItem(item model.Item) error {
	err := mda.Db.Insert(&item)
	if err != nil {
		return err
	}
	return nil
}

// UpdateItem itemsテーブルのデータを1件更新する
func (mda *MysqlDbAccessor) UpdateItem(item model.Item) error {
	_, err := mda.Db.Update(&item)
	if err != nil {
		return err
	}
	return nil
}

// UpdateItem itemsテーブルのデータを1件削除する
func (mda *MysqlDbAccessor) DeleteItem(item model.Item) error {
	_, err := mda.Db.Delete(&item)
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
