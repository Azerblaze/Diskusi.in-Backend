package dashboard

import (
	"discusiin/helper"
	"discusiin/services/dashboard"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	dashboard.IDashboardServices
}

func (h *DashboardHandler) GetTotalCountOfUserAndTopicAndPost(c echo.Context) error {
	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//get all total of user, topic and post
	userCount, topicCount, postCount, errTotal := h.IDashboardServices.GetTotalCountOfUserAndTopicAndPost(token)
	if errTotal != nil {
		return errTotal
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":    "Success",
		"userTotal":  userCount,
		"topicTotal": topicCount,
		"postTotal":  postCount,
	})
}
