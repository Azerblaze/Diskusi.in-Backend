package replies

import (
	"discusiin/helper"
	"discusiin/models"
	"discusiin/services/replies"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ReplyHandler struct {
	replies.IReplyServices
}

func (h *ReplyHandler) CreateReply(c echo.Context) error {
	var reply models.Reply
	// c.Bind(&reply)
	errBind := c.Bind(&reply)
	if errBind != nil {
		return errBind
	}
	if reply.Body == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Reply should not be empty")
	}

	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get comment id
	if c.Param("commentId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "commentId should not be empty")
	}
	commentId, errAtoi := strconv.Atoi(c.Param("commentId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	err := h.IReplyServices.CreateReply(reply, commentId, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Reply created",
	})
}

func (h *ReplyHandler) GetAllReplyByCommentID(c echo.Context) error {
	//get comment id
	if c.Param("commentId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "commentId should not be empty")
	}
	commentId, errAtoi := strconv.Atoi(c.Param("commentId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	//get all reply from comment
	replies, err := h.IReplyServices.GetAllReply(commentId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
		"data":    replies,
	})
}

func (h *ReplyHandler) UpdateReply(c echo.Context) error {
	var newReply models.Reply
	errBind := c.Bind(&newReply)
	if errBind != nil {
		return errBind
	}

	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get reply id
	if c.Param("replyId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "replyId should not be empty")
	}
	replyId, errAtoi := strconv.Atoi(c.Param("replyId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	//update reply
	err := h.IReplyServices.UpdateReply(newReply, replyId, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Reply updated",
	})
}

func (h *ReplyHandler) DeleteReply(c echo.Context) error {
	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get reply id
	if c.Param("replyId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "replyId should not be empty")
	}
	replyId, errAtoi := strconv.Atoi(c.Param("replyId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	//delete reply
	err := h.IReplyServices.DeleteReply(replyId, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Reply deleted",
	})
}
