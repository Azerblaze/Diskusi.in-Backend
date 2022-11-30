package replys

import (
	"discusiin/helper"
	"discusiin/models"
	"discusiin/services/replys"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ReplyHandler struct {
	replys.IReplyServices
}

func (h *ReplyHandler) CreateReply(c echo.Context) error {
	var reply models.Reply
	// c.Bind(&reply)
	errBind := c.Bind(&reply)
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
	// reply.UserID = 1 //untuk percobaan

	//get comment id
	commentId, errAtoi := strconv.Atoi(c.Param("comment_id"))
	if errAtoi != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errAtoi.Error(),
		})
	}

	err := h.IReplyServices.CreateReply(reply, commentId, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Reply Created",
	})
}

func (h *ReplyHandler) GetAllReply(c echo.Context) error {
	//get comment id
	commentId, errAtoi := strconv.Atoi(c.Param("comment_id"))
	if errAtoi != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errAtoi.Error(),
		})
	}

	//get all reply from comment
	replys, err := h.IReplyServices.GetAllReply(commentId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Succes",
		"data":    replys,
	})
}

func (h *ReplyHandler) UpdateReply(c echo.Context) error {
	var newReply models.Reply
	errBind := c.Bind(&newReply)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errBind.Error(),
			// "message": errors.New("hallo semuanya"),
		})
	}

	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errDecodeJWT.Error(),
			// "message": errors.New("hai"),
		})
	}

	//get reply id
	replyId, errAtoi := strconv.Atoi(c.Param("reply_id"))
	if errAtoi != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errAtoi.Error(),
		})
	}

	//update reply
	err := h.IReplyServices.UpdateReply(newReply, replyId, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			// "message": errors.New("hallo semua"),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Reply Updated",
	})
}

func (h *ReplyHandler) DeleteReply(c echo.Context) error {
	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errDecodeJWT.Error(),
			// "message": errors.New("hai"),
		})
	}

	//get reply id
	replyId, errAtoi := strconv.Atoi(c.Param("reply_id"))
	if errAtoi != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": errAtoi.Error(),
		})
	}

	//delete reply
	err := h.IReplyServices.DeleteReply(replyId, token)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Reply Deleted",
	})
}
