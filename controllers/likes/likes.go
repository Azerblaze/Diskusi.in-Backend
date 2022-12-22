package likes

import (
	"discusiin/helper"
	"discusiin/services/likes"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type LikeHandler struct {
	likes.ILikeServices
}

func (h *LikeHandler) LikePost(c echo.Context) error {

	//get logged userid
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get post id
	postIDStr := c.Param("postId")
	if postIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "postId parameter should not be emtpy")
	}
	postId, errAtoi := strconv.Atoi(postIDStr)
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	err := h.ILikeServices.LikePost(token, postId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
	})
}
func (h *LikeHandler) DislikePost(c echo.Context) error {

	//get logged userid
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get post id
	postId, errAtoi := strconv.Atoi(c.Param("postId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	err := h.ILikeServices.DislikePost(token, postId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
	})
}
