package handler

import (
	"net/http"

	"github.com/greytabby/meowapi/lib/model"
	"github.com/labstack/echo"
)

type ToiletReader interface {
	GetAllToilets() ([]model.Toilet, error)
	GetToilet(id int64) (model.Toilet, error)
}

type ToiletManipulater interface {
	AddToilet(cat model.Toilet) error
	UpdateToilet(cat model.Toilet) error
	DeleteToilet(cat model.Toilet) error
}

// ToiletDbAccessor toiletテーブルを操作するinterface
type ToiletDbAccessor interface {
	ToiletReader
	ToiletManipulater
}

// ToiletHandler /api/toiletへのリクエストを処理する
type ToiletHandler struct {
	Db ToiletDbAccessor
}

// GetAllToilets Toiletテーブルから全てのToiletを返す
func (th *ToiletHandler) GetAllToilets(c echo.Context) error {
	toilets, err := th.Db.GetAllToilets()
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, toilets)
}

// AddToilet Toiletテーブルへtoiletを1追加する
func (th *ToiletHandler) AddToilet(c echo.Context) error {
	var toilet model.Toilet

	if err := c.Bind(&toilet); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}
	if err := th.Db.AddToilet(toilet); err != nil {
		c.Logger().Errorf("Insert: ", err)
		return c.String(http.StatusBadRequest, "Could not add new toilet.")
	}
	c.Logger().Infof("Added: %#v", toilet)
	return c.String(http.StatusOK, "")
}

// UpdateToilet Toiletの情報を1件更新する
func (th *ToiletHandler) UpdateToilet(c echo.Context) error {
	var toilet, selectedToilet model.Toilet

	if err := c.Bind(&toilet); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if toilet.Id == 0 {
		return c.String(http.StatusBadRequest, "Toilet id not specified.")
	}

	// Get cat from db for confirming wheather the user specified cat exist.
	selectedToilet, err := th.Db.GetToilet(toilet.Id)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "No your specified cat in the database.")
	}

	// Update information.
	selectedToilet.Name = toilet.Name
	selectedToilet.Comment = toilet.Comment
	selectedToilet.SandState = toilet.SandState
	if err := th.Db.UpdateToilet(selectedToilet); err != nil {
		c.Logger().Errorf("Update: ", err)
		return c.String(http.StatusBadRequest, "Could not update toilet info.")
	}
	c.Logger().Infof("Updated: %#v", selectedToilet)
	return c.String(http.StatusOK, "")
}

// DeleteToilet Toiletを1件削除する
func (th *ToiletHandler) DeleteToilet(c echo.Context) error {
	var toilet model.Toilet

	if err := c.Bind(&toilet); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if toilet.Id == 0 {
		return c.String(http.StatusBadRequest, "Toilet id is not specified.")
	}

	selectedToilet, err := th.Db.GetToilet(toilet.Id)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "No your specified toilet.")
	}

	if err := th.Db.DeleteToilet(selectedToilet); err != nil {
		c.Logger().Errorf("Delete: ", err)
		return c.String(http.StatusBadRequest, "Failed delete the toilet.")
	}
	c.Logger().Infof("Deleted: %#v", selectedToilet)
	return c.String(http.StatusOK, "")
}
