package posts

import (
	"discusiin/helper"
	"discusiin/models"
	"discusiin/services/posts"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	posts.IPostServices
}

func (h *PostHandler) CreateNewPost(c echo.Context) error {
	var p models.Post
	errBind := c.Bind(&p)
	if errBind != nil {
		return errBind
	}

	url_param_value := c.Param("topic_name")
	if url_param_value == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topic_name should not be empty")
	}
	topicName := helper.URLDecodeReformat(url_param_value)

	if p.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "post title should not be empty")
	}
	if p.Body == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "post body should not be empty")
	}

	//get logged userId
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	err := h.IPostServices.CreatePost(p, topicName, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Post created",
	})
}

func (h *PostHandler) GetAllPost(c echo.Context) error {
	if c.Param("topic_name") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topic name should not be empty")
	}
	topicName := helper.URLDecodeReformat(c.Param("topic_name"))
	if topicName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topic name should not be empty")
	}
	//check if page exist
	if c.QueryParam("page") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "page query parameter should not be empty")
	}
	page, errAtoi := strconv.Atoi(c.QueryParam("page"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	posts, numberOfPage, err := h.IPostServices.GetPosts(topicName, page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":        "Success",
		"data":           posts,
		"number_of_page": numberOfPage,
		"page":           page,
	})
}

func (h *PostHandler) GetPost(c echo.Context) error {

	if c.Param("post_id") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "post_id parameter should not be empty")
	}
	id, errAtoi := strconv.Atoi(c.Param("post_id"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	p, err := h.IPostServices.GetPost(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
		"data":    p,
	})
}

func (h *PostHandler) EditPost(c echo.Context) error {
	var newPost models.Post
	errBind := c.Bind(&newPost)
	if errBind != nil {
		return errBind
	}

	//get user id from logged user
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}
	if c.Param("post_id") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "post_id parameter should not be empty")
	}
	id, errAtoi := strconv.Atoi(c.Param("post_id"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	err := h.IPostServices.UpdatePost(newPost, id, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Post updated",
	})
}

func (h *PostHandler) DeletePost(c echo.Context) error {

	//get user id from logged user
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}
	if c.Param("post_id") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "post_id parameter should not be empty")
	}
	postID, errAtoi := strconv.Atoi(c.Param("post_id"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}
	err := h.IPostServices.DeletePost(postID, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Post deleted",
	})
}

func (h *PostHandler) GetRecentPost(c echo.Context) error {
	//check if page exist
	pageStr := c.QueryParam("page")
	if pageStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "page parameter should not be empty")
	}
	page, _ := strconv.Atoi(pageStr)

	posts, numberOfPage, err := h.IPostServices.GetRecentPost(page)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message":        "Success",
		"data":           posts,
		"number_of_page": numberOfPage,
		"page":           page,
	})
}
