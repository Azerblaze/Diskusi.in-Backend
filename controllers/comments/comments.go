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
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errBind.Error(),
		})
	}

	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errDecodeJWT.Error(),
		})
	}

	postID, errAtoi := strconv.Atoi(c.Param("post_id"))
	if errAtoi != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errAtoi.Error(),
		})
	}

	err := h.ICommentServices.CreateComment(comment, postID, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Comment Created",
	})

}

func (h *CommentHandler) GetAllComment(c echo.Context) error {

	//get post id
	postID, errAtoi := strconv.Atoi(c.Param("post_id"))
	if errAtoi != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errAtoi.Error(),
		})
	}
	//get all coment from post
	comments, err := h.ICommentServices.GetAllComments(postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Succes",
		"data":    comments,
	})
}

func (h *CommentHandler) UpdateComment(c echo.Context) error {
	var comment models.Comment
	errBind := c.Bind(&comment)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errBind.Error(),
		})
	}

	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errDecodeJWT.Error(),
		})
	}

	//check if user who eligibleuntuk param comment
	commentID, errAtoi := strconv.Atoi(c.Param("comment_id"))
	if errAtoi != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errAtoi.Error(),
		})
	}

	comment.ID = uint(commentID)

	err := h.ICommentServices.UpdateComment(comment, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Comment Updated",
	})

}

func (h *CommentHandler) DeleteComment(c echo.Context) error {
	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errDecodeJWT.Error(),
		})
	}

	//check if user who eligible
	commentID, errAtoi := strconv.Atoi(c.Param("comment_id"))
	if errAtoi != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errAtoi.Error(),
		})
	}
	err := h.ICommentServices.DeleteComment(commentID, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Comment Deleted",
	})
}
