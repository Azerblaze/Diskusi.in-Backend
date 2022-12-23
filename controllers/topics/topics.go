package topics

import (
	"discusiin/helper"
	"discusiin/models"
	"discusiin/services/topics"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TopicHandler struct {
	topics.ITopicServices
}

func (h *TopicHandler) CreateNewTopic(c echo.Context) error {
	// validation
	var topic models.Topic

	errBind := c.Bind(&topic)
	if errBind != nil {
		return errBind
	}

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//is title empty
	if topic.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topic name should not be empty")
	}

	//is description empty
	if topic.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "description name should not be empty")
	}

	// is description to short
	if len(topic.Description) < 25 {
		return echo.NewHTTPError(http.StatusBadRequest, "description to short, at least cpntain 25 character")
	}

	result, err := h.ITopicServices.CreateTopic(topic, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Topic created",
		"data":    result,
	})
}
func (h *TopicHandler) GetNumberOfPostOnATopicByTopicName(c echo.Context) error {
	topicName := c.Param("topicName")
	if topicName == "" {
		err := echo.NewHTTPError(http.StatusBadRequest, "topicName should not be empty")
		if err == nil {
			panic("unexpected nil error")
		}
		return err
	}

	numberOfPost, err := h.ITopicServices.GetNumberOfPostOnATopicByTopicName(topicName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"data":    numberOfPost,
	})
}
func (h *TopicHandler) GetTopTopics(c echo.Context) error {
	topics, err := h.ITopicServices.GetTopTopics()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
		"data":    topics,
	})
}
func (h *TopicHandler) GetAllTopics(c echo.Context) error {
	topics, err := h.ITopicServices.GetTopics()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"data":    topics,
	})
}

func (h *TopicHandler) GetTopic(c echo.Context) error {

	topicIdStr := c.Param("topicId")
	if topicIdStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topicId should not be empty")
	}
	topicId, errAtoi := strconv.Atoi(topicIdStr)
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	topic, err := h.ITopicServices.GetTopic(topicId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"data":    topic,
	})
}

func (h *TopicHandler) UpdateTopicDescription(c echo.Context) error {

	// validation
	var newTopic models.Topic
	errBind := c.Bind(&newTopic)
	if errBind != nil {
		return errBind
	}

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	idStr := c.Param("topicId")
	if idStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topicId should not be empty")
	}
	id, errAtoi := strconv.Atoi(idStr)
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	//is description empty
	if newTopic.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "description name should not be empty")
	}

	// is description to short
	if len(newTopic.Description) < 25 {
		return echo.NewHTTPError(http.StatusBadRequest, "description to short, at least contain 25 character")
	}

	newTopic.ID = uint(id)

	topic, err := h.ITopicServices.UpdateTopicDescription(newTopic, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Topic updated",
		"data":    topic,
	})
}

func (h *TopicHandler) DeleteTopic(c echo.Context) error {
	idStr := c.Param("topicId")
	if idStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topicId should not be empty")
	}
	id, errAtoi := strconv.Atoi(idStr)
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	err := h.ITopicServices.RemoveTopic(token, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Topic deleted",
	})
}
