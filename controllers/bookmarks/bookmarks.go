package bookmarks

import (
	"discusiin/helper"
	"discusiin/services/bookmarks"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type BookmarkHandler struct {
	bookmarks.IBookmarkServices
}

func (h *BookmarkHandler) AddBookmark(c echo.Context) error {
	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get postId
	postID, errAtoi := strconv.Atoi(c.Param("post_id"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	//add bookmark
	err := h.IBookmarkServices.AddBookmark(token, postID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Bookmark added",
	})
}

func (h *BookmarkHandler) DeleteBookmark(c echo.Context) error {
	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get bookmarkID
	bookmarkID, errAtoi := strconv.Atoi(c.Param("bookmark_id"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	//delete bookmark
	err := h.IBookmarkServices.DeleteBookmark(token, bookmarkID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Bookmark deleted",
	})
}

func (h *BookmarkHandler) GetAllBookmark(c echo.Context) error {
	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get all bookmark
	bookmarks, err := h.IBookmarkServices.GetAllBookmark(token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Success",
		"data":    bookmarks,
	})
}
