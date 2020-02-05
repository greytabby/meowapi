package handler

import (
	"net/http"

	"github.com/greytabby/meowapi/lib/model"
	"github.com/labstack/echo"
)

type CatReader interface {
	GetAllCats(uid int64) ([]model.Cat, error)
	GetCat(id, uid int64) (model.Cat, error)
}

type CatManipulator interface {
	AddCat(cat model.Cat) error
	UpdateCat(cat model.Cat) error
	DeleteCat(cat model.Cat) error
}

// CatDbAccessor catテーブルを操作するinterface
type CatDbAccessor interface {
	CatReader
	CatManipulator
}

// CatHandler /api/catへのリクエストを処理する
type CatHandler struct {
	Db CatDbAccessor
}

// GetAllCats catテーブルから全てのcatを返す
func (ch *CatHandler) GetAllCats(c echo.Context) error {
	uid := UserIdFromToken(c)
	cats, err := ch.Db.GetAllCats(uid)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, cats)
}

// AddCat catテーブルへcatを1匹追加する
func (ch *CatHandler) AddCat(c echo.Context) error {
	var cat model.Cat
	uid := UserIdFromToken(c)

	if err := c.Bind(&cat); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	cat.UID = uid
	if err := ch.Db.AddCat(cat); err != nil {
		c.Logger().Errorf("Insert: ", err)
		return c.String(http.StatusBadRequest, "Could not add new cat.")
	}
	c.Logger().Infof("Added: %#v", cat)
	return c.String(http.StatusOK, "")
}

// UpdateCat catの情報を1件更新する
func (ch *CatHandler) UpdateCat(c echo.Context) error {
	var cat, selectedCat model.Cat

	if err := c.Bind(&cat); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if cat.Id == 0 {
		return c.String(http.StatusBadRequest, "Cat id not specified.")
	}

	// Get cat from db for confirming wheather the user specified cat exist.
	uid := UserIdFromToken(c)
	selectedCat, err := ch.Db.GetCat(cat.Id, uid)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "No your specified cat in the database.")
	}

	// Update cat information.
	selectedCat.Name = cat.Name
	selectedCat.Breed = cat.Breed
	selectedCat.Gender = cat.Gender
	selectedCat.Age = cat.Age
	if err := ch.Db.UpdateCat(selectedCat); err != nil {
		c.Logger().Errorf("Update: ", err)
		return c.String(http.StatusBadRequest, "Could not update cat info.")
	}
	c.Logger().Infof("Updated: %#v", selectedCat)
	return c.String(http.StatusOK, "")
}

// DeleteCat catを1匹登録削除する
func (ch *CatHandler) DeleteCat(c echo.Context) error {
	var cat model.Cat

	if err := c.Bind(&cat); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if cat.Id == 0 {
		return c.String(http.StatusBadRequest, "Cat id is not specified.")
	}

	uid := UserIdFromToken(c)
	selectedCat, err := ch.Db.GetCat(cat.Id, uid)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "No your specified cat.")
	}

	if err := ch.Db.DeleteCat(selectedCat); err != nil {
		c.Logger().Errorf("Delete: ", err)
		return c.String(http.StatusBadRequest, "Failed delete the cat.")
	}
	c.Logger().Infof("Deleted: %#v", selectedCat)
	return c.String(http.StatusOK, "")
}
