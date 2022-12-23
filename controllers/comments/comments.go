package comments

import (
	"discusiin/helper"
	"discusiin/models"
	comments "discusiin/services/comments"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CommentHandler struct {
	comments.ICommentServices
}

func (h *CommentHandler) CreateComment(c echo.Context) error {
	var comment models.Comment
	errBind := c.Bind(&comment)
	if errBind != nil {
		return echo.NewHTTPError(http.StatusUnsupportedMediaType, errBind.Error())
	}

	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	postId, errAtoi := strconv.Atoi(c.Param("postId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	err := h.ICommentServices.CreateComment(comment, postId, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Comment created",
	})

}

func (h *CommentHandler) GetAllCommentByPostID(c echo.Context) error {

	//get post id
	if c.Param("postId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "postId should not be empty")
	}
	postId, errAtoi := strconv.Atoi(c.Param("postId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}
	//get all coment from post
	comments, err := h.ICommentServices.GetAllComments(postId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
		"data":    comments,
	})
}

func (h *CommentHandler) UpdateComment(c echo.Context) error {
	var comment models.Comment
	errBind := c.Bind(&comment)
	if errBind != nil {
		return errBind
	}
	if comment.Body == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Comment should not be empty")
	}

	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//check if user who eligible untuk param comment
	if c.Param("commentId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "commentId should not be empty")
	}
	commentId, errAtoi := strconv.Atoi(c.Param("commentId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	comment.ID = uint(commentId)

	err := h.ICommentServices.UpdateComment(comment, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Comment updated",
	})

}

func (h *CommentHandler) DeleteComment(c echo.Context) error {
	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//check if user who eligible
	if c.Param("commentId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "commentId should not be empty")
	}
	commentId, errAtoi := strconv.Atoi(c.Param("commentId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}
	err := h.ICommentServices.DeleteComment(commentId, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Comment deleted",
	})
}
