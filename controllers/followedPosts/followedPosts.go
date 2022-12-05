package followedPosts

import (
	"discusiin/helper"
	"discusiin/services/followedPosts"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type FollowedPostHandler struct {
	followedPosts.IFollowedPostServices
}

func (h *FollowedPostHandler) AddFollowedPost(c echo.Context) error {
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

	//add followedPost
	err := h.IFollowedPostServices.AddFollowedPost(token, postID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Post followed",
	})
}

func (h *FollowedPostHandler) DeleteFollowedPost(c echo.Context) error {
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

	//delete followedPost
	err := h.IFollowedPostServices.DeleteFollowedPost(token, postID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Post unfollowed",
	})
}

func (h *FollowedPostHandler) GetAllFollowedPost(c echo.Context) error {
	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get all followedPost
	followedPosts, err := h.IFollowedPostServices.GetAllFollowedPost(token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":           "Success",
		"data_followedPost": followedPosts,
	})
}
