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
func (mda *MysqlDbAccessor) GetAllCats(uid int64) ([]model.Cat, error) {
	var cats []model.Cat
	_, err := mda.Db.Select(&cats,
		"SELECT * FROM cat WHERE uid = ? ORDER BY created", uid)
	if err != nil {
		return nil, err
	}
	return cats, nil
}

// GetCat DBのcatテーブルからidに合致するcatを1つ返す
// 見つからなかった場合は空のcatとerrorを返す
func (mda *MysqlDbAccessor) GetCat(id, uid int64) (model.Cat, error) {
	var cat model.Cat
	err := mda.Db.SelectOne(&cat, "SELECT * FROM cat WHERE id = ? AND uid = ?", id, uid)
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

// GetAllToilets DBからcatテーブルの全てのデータを取得する
func (mda *MysqlDbAccessor) GetAllToilets(uid int64) ([]model.Toilet, error) {
	var toilets []model.Toilet
	_, err := mda.Db.Select(&toilets,
		"SELECT * FROM toilet WHERE uid = ? ORDER BY created", uid)
	if err != nil {
		return nil, err
	}
	return toilets, nil
}

// GetToilet DBのtoiletテーブルからidに合致するtoiletを1つ返す
// 見つからなかった場合は空のtoiletとerrorを返す
func (mda *MysqlDbAccessor) GetToilet(id, uid int64) (model.Toilet, error) {
	var toilet model.Toilet
	err := mda.Db.SelectOne(&toilet, "SELECT * FROM toilet WHERE id = ? AND uid = ?", id, uid)
	if err != nil {
		return model.Toilet{}, err
	}
	return toilet, nil
}

// AddToilet toiletテーブルへデータを1件追加する
func (mda *MysqlDbAccessor) AddToilet(toilet model.Toilet) error {
	err := mda.Db.Insert(&toilet)
	if err != nil {
		return err
	}
	return nil
}

// UpdateToilet toiletテーブルのデータを1件更新する
func (mda *MysqlDbAccessor) UpdateToilet(toilet model.Toilet) error {
	_, err := mda.Db.Update(&toilet)
	if err != nil {
		return err
	}
	return nil
}

// DeleteToilet toiletテーブルのデータを1件削除する
func (mda *MysqlDbAccessor) DeleteToilet(toilet model.Toilet) error {
	_, err := mda.Db.Delete(&toilet)
	if err != nil {
		return err
	}
	return nil
}

// GetAllUseToilets DBからcatテーブルの全てのデータを取得する
func (mda *MysqlDbAccessor) GetAllUseToilets(uid int64) ([]model.UseToilet, error) {
	var usetoilets []model.UseToilet
	_, err := mda.Db.Select(&usetoilets,
		"SELECT * FROM usetoilet WHERE uid = ? ORDER BY created", uid)
	if err != nil {
		return nil, err
	}
	return usetoilets, nil
}

// GetUseToilet DBのusetoiletテーブルからidに合致するusetoiletを1つ返す
// 見つからなかった場合は空のusetoiletとerrorを返す
func (mda *MysqlDbAccessor) GetUseToilet(id, uid int64) (model.UseToilet, error) {
	var usetoilet model.UseToilet
	err := mda.Db.SelectOne(&usetoilet, "SELECT * FROM usetoilet WHERE id = ? AND uid = ?", id, uid)
	if err != nil {
		return model.UseToilet{}, err
	}
	return usetoilet, nil
}

// AddUseToilet usetoiletテーブルへデータを1件追加する
func (mda *MysqlDbAccessor) AddUseToilet(usetoilet model.UseToilet) error {
	err := mda.Db.Insert(&usetoilet)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUseToilet usetoiletテーブルのデータを1件更新する
func (mda *MysqlDbAccessor) UpdateUseToilet(usetoilet model.UseToilet) error {
	_, err := mda.Db.Update(&usetoilet)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUseToilet usetoiletテーブルのデータを1件削除する
func (mda *MysqlDbAccessor) DeleteUseToilet(usetoilet model.UseToilet) error {
	_, err := mda.Db.Delete(&usetoilet)
	if err != nil {
		return err
	}
	return nil
}

// GetAllWashes washテーブルの全てのデータを取得する
func (mda *MysqlDbAccessor) GetAllWashes(uid int64) ([]model.Wash, error) {
	var ws []model.Wash
	_, err := mda.Db.Select(&ws,
		"SELECT * FROM wash WHERE uid = ? ORDER BY created", uid)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// GetWashesByToiletId washテーブルから特定のToiletIdの全てのデータを取得する
func (mda *MysqlDbAccessor) GetWashesByToiletId(toiletid, uid int64) ([]model.Wash, error) {
	var ws []model.Wash
	_, err := mda.Db.Select(&ws,
		"SELECT * FROM wash WHERE toiletid = ? AND uid = ? ORDER BY created", toiletid, uid)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

// GetWash DBのwashテーブルからidに合致するwashを1つ返す
// 見つからなかった場合は空のwashとerrorを返す
func (mda *MysqlDbAccessor) GetWash(id, uid int64) (model.Wash, error) {
	var w model.Wash
	err := mda.Db.SelectOne(&w, "SELECT * FROM wash WHERE id = ? AND uid = ?", id, uid)
	if err != nil {
		return model.Wash{}, err
	}
	return w, nil
}

// AddWash washテーブルへデータを1件追加する
func (mda *MysqlDbAccessor) AddWash(wash model.Wash) error {
	err := mda.Db.Insert(&wash)
	if err != nil {
		return err
	}
	return nil
}

// UpdateWash washテーブルのデータを1件更新する
func (mda *MysqlDbAccessor) UpdateWash(wash model.Wash) error {
	_, err := mda.Db.Update(&wash)
	if err != nil {
		return err
	}
	return nil
}

// DeleteWash usetoiletテーブルのデータを1件削除する
func (mda *MysqlDbAccessor) DeleteWash(wash model.Wash) error {
	_, err := mda.Db.Delete(&wash)
	if err != nil {
		return err
	}
	return nil
}

// FindUser userテーブルからnameに合致するデータを1件取得する
func (mda *MysqlDbAccessor) FindUser(name string) (model.User, error) {
	var u model.User
	err := mda.Db.SelectOne(&u, "Select * FROM user WHERE name = ?", name)
	if err != nil {
		return model.User{}, err
	}
	return u, nil
}

func (mda *MysqlDbAccessor) AddUser(user model.User) error {
	err := mda.Db.Insert(&user)
	if err != nil {
		return err
	}
	return nil
}

func (mda *MysqlDbAccessor) DeleteUser(user model.User) error {
	_, err := mda.Db.Delete(&user)
	if err != nil {
		return err
	}
	return nil
}
