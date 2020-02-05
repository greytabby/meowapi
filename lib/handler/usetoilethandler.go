package handler

import (
	"net/http"

	"github.com/greytabby/meowapi/lib/model"
	"github.com/labstack/echo"
)

type UseToiletReader interface {
	GetAllUseToilets(uid int64) ([]model.UseToilet, error)
	GetUseToilet(id, uid int64) (model.UseToilet, error)
}

type UseToiletManipulator interface {
	AddUseToilet(ut model.UseToilet) error
	UpdateUseToilet(ut model.UseToilet) error
	DeleteUseToilet(ut model.UseToilet) error
}

// UseToiletDbAccessor usetoiletテーブルを操作するinterface
type UseToiletDbAccessor interface {
	UseToiletReader
	UseToiletManipulator
}

// UseToiletHandler /api/usetoiletへのリクエストを処理する
type UseToiletHandler struct {
	Db UseToiletDbAccessor
}

// GetAllUseToilets UseToiletテーブルから全てのUseToiletを返す
func (th *UseToiletHandler) GetAllUseToilets(c echo.Context) error {
	uid := UserIdFromToken(c)
	usetoilets, err := th.Db.GetAllUseToilets(uid)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, usetoilets)
}

// AddUseToilet UseToiletテーブルへusetoiletを1追加する
func (th *UseToiletHandler) AddUseToilet(c echo.Context) error {
	var usetoilet model.UseToilet

	if err := c.Bind(&usetoilet); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	uid := UserIdFromToken(c)
	usetoilet.UID = uid
	if err := th.Db.AddUseToilet(usetoilet); err != nil {
		c.Logger().Errorf("Insert: ", err)
		return c.String(http.StatusBadRequest, "Could not add new usetoilet.")
	}
	c.Logger().Infof("Added: %#v", usetoilet)
	return c.String(http.StatusOK, "")
}

// UpdateUseToilet UseToiletの情報を1件更新する
func (th *UseToiletHandler) UpdateUseToilet(c echo.Context) error {
	var usetoilet, selectedUseToilet model.UseToilet

	if err := c.Bind(&usetoilet); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if usetoilet.Id == 0 {
		return c.String(http.StatusBadRequest, "UseToilet id not specified.")
	}

	// Get cat from db for confirming wheather the user specified cat exist.
	uid := UserIdFromToken(c)
	selectedUseToilet, err := th.Db.GetUseToilet(usetoilet.Id, uid)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "No your specified cat in the database.")
	}

	// Update information.
	selectedUseToilet.ToiletId = usetoilet.ToiletId
	selectedUseToilet.CatId = usetoilet.CatId
	selectedUseToilet.Type = usetoilet.Type
	if err := th.Db.UpdateUseToilet(selectedUseToilet); err != nil {
		c.Logger().Errorf("Update: ", err)
		return c.String(http.StatusBadRequest, "Could not update usetoilet info.")
	}
	c.Logger().Infof("Updated: %#v", selectedUseToilet)
	return c.String(http.StatusOK, "")
}

// DeleteUseToilet UseToiletを1件削除する
func (th *UseToiletHandler) DeleteUseToilet(c echo.Context) error {
	var usetoilet model.UseToilet

	if err := c.Bind(&usetoilet); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if usetoilet.Id == 0 {
		return c.String(http.StatusBadRequest, "UseToilet id is not specified.")
	}

	uid := UserIdFromToken(c)
	selectedUseToilet, err := th.Db.GetUseToilet(usetoilet.Id, uid)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "No your specified usetoilet.")
	}

	if err := th.Db.DeleteUseToilet(selectedUseToilet); err != nil {
		c.Logger().Errorf("Delete: ", err)
		return c.String(http.StatusBadRequest, "Failed delete the usetoilet.")
	}
	c.Logger().Infof("Deleted: %#v", selectedUseToilet)
	return c.String(http.StatusOK, "")
}
