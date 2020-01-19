package handler

import (
	"net/http"
	"strconv"

	"github.com/greytabby/meowapi/lib/model"
	"github.com/labstack/echo"
)

// WashReader washテーブルを参照する
type WashReader interface {
	GetAllWashes() ([]model.Wash, error)
	GetWashesByToiletId(toiletid int64) ([]model.Wash, error)
	GetWash(id int64) (model.Wash, error)
}

// WashManipulater washテーブルを操作する
type WashManipulator interface {
	AddWash(wash model.Wash) error
	UpdateWash(wash model.Wash) error
	DeleteWash(wash model.Wash) error
}

// WashDbAccessor washテーブルの参照/操作を行う
type WashDbAccessor interface {
	WashReader
	WashManipulator
}

// WashHandler /api/washへのリクエストを処理する
type WashHandler struct {
	Db WashDbAccessor
}

// GetAllWashed 全てのwashを取得する
func (wh *WashHandler) GetAllWashes(c echo.Context) error {
	washes, err := wh.Db.GetAllWashes()
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusInternalServerError, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, washes)
}

// GetWashedByToiletId 特定のToiletIdのwashを取得する
func (wh *WashHandler) GetWashesByToiletId(c echo.Context) error {
	toiletid, err := strconv.ParseInt(c.Param("toiletid"), 10, 64)
	if err != nil {
		c.Logger().Errorf("Param parse: ", err)
		return c.String(http.StatusBadRequest, "Param parse: "+err.Error())
	}
	washes, err := wh.Db.GetWashesByToiletId(toiletid)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusInternalServerError, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, washes)
}

// AddWash washを1件登録する
func (wh *WashHandler) AddWash(c echo.Context) error {
	var w model.Wash
	if err := c.Bind(&w); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}
	if err := wh.Db.AddWash(w); err != nil {
		c.Logger().Error("Insert: ", err)
		return c.String(http.StatusInternalServerError, "Could not add wash.")
	}
	c.Logger().Infof("Added: %#v", w)
	return c.JSON(http.StatusOK, "")
}

// UpdateWash washを1件更新する
func (wh *WashHandler) UpdateWash(c echo.Context) error {
	var w, selected model.Wash
	if err := c.Bind(&w); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if w.Id == 0 {
		return c.String(http.StatusBadRequest, "Wash id was not specifyed.")
	}

	selected, err := wh.Db.GetWash(w.Id)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusInternalServerError, "Select: "+err.Error())
	}

	selected.ToiletId = w.ToiletId
	selected.Comment = w.Comment
	if err := wh.Db.UpdateWash(selected); err != nil {
		c.Logger().Error("Update: ", err)
		return c.String(http.StatusInternalServerError, "Could not update wash.")
	}
	return c.JSON(http.StatusOK, w)
}

// UpdateWash washを1件削除する
func (wh *WashHandler) DeleteWash(c echo.Context) error {
	var w, selected model.Wash
	if err := c.Bind(&w); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if w.Id == 0 {
		return c.String(http.StatusBadRequest, "Wash id was not specifyed.")
	}

	selected, err := wh.Db.GetWash(w.Id)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusInternalServerError, "Select: "+err.Error())
	}

	if err := wh.Db.DeleteWash(selected); err != nil {
		c.Logger().Error("Delete: ", err)
		return c.String(http.StatusInternalServerError, "Could not delete wash.")
	}
	return c.JSON(http.StatusOK, w)
}
