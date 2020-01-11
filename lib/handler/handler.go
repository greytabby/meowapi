package handler

import (
	"net/http"

	"github.com/greytabby/meowapi/lib/model"
	"github.com/labstack/echo"
)

// DbAccessor Handlerに必要なDBアクセスを提供するインタフェース
type DbAccessor interface {
	GetAllItems() ([]model.Item, error)
	GetItem(id int64) (model.Item, error)
	InsertItem(item model.Item) error
	UpdateItem(item model.Item) error
	DeleteItem(item model.Item) error
}

// Handler /itemへのapiリクエストを処理するHandler
// DBアクセスのためのAccessorをもつ
type Handler struct {
	Db DbAccessor
}

// GetItems 全てのitemを返す
func (h *Handler) GetItems(c echo.Context) error {
	items, err := h.Db.GetAllItems()
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}
	return c.JSON(http.StatusOK, items)
}

// InsertItem itemを1件追加する
func (h *Handler) InsertItem(c echo.Context) error {
	var item model.Item

	if err := c.Bind(&item); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}
	if err := h.Db.InsertItem(item); err != nil {
		c.Logger().Errorf("Insert: ", err)
		return c.String(http.StatusBadRequest, "Insert: "+err.Error())
	}
	c.Logger().Infof("Added: %v", item.Id)
	return c.String(http.StatusOK, "")
}

// UpdateItem itemを1件更新する
func (h *Handler) UpdateItem(c echo.Context) error {
	var item, selectedItem model.Item

	if err := c.Bind(&item); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if item.Id == 0 {
		return c.String(http.StatusBadRequest, "Item id not specified.")
	}

	selectedItem, err := h.Db.GetItem(item.Id)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}

	selectedItem.Name = item.Name
	selectedItem.Text = item.Text
	if err := h.Db.UpdateItem(selectedItem); err != nil {
		c.Logger().Errorf("Update: ", err)
		return c.String(http.StatusBadRequest, "Update: "+err.Error())
	}
	c.Logger().Infof("Updated: %v", item.Id)
	return c.String(http.StatusOK, "")
}

// DeleteItem itemを1件削除する
func (h *Handler) DeleteItem(c echo.Context) error {
	var item model.Item

	if err := c.Bind(&item); err != nil {
		c.Logger().Errorf("Bind: ", err)
		return c.String(http.StatusBadRequest, "Bind: "+err.Error())
	}

	if item.Id == 0 {
		return c.String(http.StatusBadRequest, "Item id is not specified.")
	}

	selectedItem, err := h.Db.GetItem(item.Id)
	if err != nil {
		c.Logger().Errorf("Select: ", err)
		return c.String(http.StatusBadRequest, "Select: "+err.Error())
	}

	if err := h.Db.DeleteItem(selectedItem); err != nil {
		c.Logger().Errorf("Delete: ", err)
		return c.String(http.StatusBadRequest, "Delete: "+err.Error())
	}
	c.Logger().Infof("Deleted: %v", item.Id)
	return c.String(http.StatusOK, "")
}
