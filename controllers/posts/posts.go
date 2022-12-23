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

	urlParamValue := c.Param("topicName")
	if urlParamValue == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topicName should not be empty")
	}
	topicName := helper.URLDecodeReformat(urlParamValue)

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

func (h *PostHandler) GetAllPostByTopicName(c echo.Context) error {
	if c.Param("topicName") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topic name should not be empty")
	}
	topicName := helper.URLDecodeReformat(c.Param("topicName"))
	if topicName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topic name should not be empty")
	}

	//insert search query param
	search := (c.QueryParam("search"))

	//check if page exist
	var page int
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		var errAtoi error
		page, errAtoi = strconv.Atoi(c.QueryParam("page"))
		if errAtoi != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
		}
	}

	posts, numberOfPage, err := h.IPostServices.GetPosts(topicName, page, search)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Success",
		"data":         posts,
		"numberOfPage": numberOfPage,
		"page":         page,
	})
}

func (h *PostHandler) GetAllPostByTopicByLike(c echo.Context) error {
	if c.Param("topicName") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topic name should not be empty")
	}
	topicName := helper.URLDecodeReformat(c.Param("topicName"))
	if topicName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "topic name should not be empty")
	}

	//check if page exist
	var page int
	if c.QueryParam("page") == "" {
		page = 1
	} else {
		var errAtoi error
		page, errAtoi = strconv.Atoi(c.QueryParam("page"))
		if errAtoi != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
		}
	}

	posts, numberOfPage, err := h.IPostServices.GetPostsByTopicByLike(topicName, page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Success",
		"data":         posts,
		"numberOfPage": numberOfPage,
		"page":         page,
	})
}

func (h *PostHandler) GetPostByPostID(c echo.Context) error {

	if c.Param("postId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "postId parameter should not be empty")
	}
	id, errAtoi := strconv.Atoi(c.Param("postId"))
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

	var postId int
	if c.Param("postId") == "" {
		postId = 1
	} else {
		var errAtoi error
		postId, errAtoi = strconv.Atoi(c.Param("postId"))
		if errAtoi != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
		}
	}

	err := h.IPostServices.UpdatePost(newPost, postId, token)
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
	if c.Param("postId") == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "postId parameter should not be empty")
	}
	postId, errAtoi := strconv.Atoi(c.Param("postId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}
	err := h.IPostServices.DeletePost(postId, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Post deleted",
	})
}

func (h *PostHandler) GetAllRecentPost(c echo.Context) error {
	//check if page exist
	pageStr := c.QueryParam("page")
	var page int
	if pageStr == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	search := c.QueryParam("search")

	posts, numberOfPage, err := h.IPostServices.GetRecentPost(page, search)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Success",
		"data":         posts,
		"numberOfPage": numberOfPage,
		"page":         page,
	})
}

func (h *PostHandler) GetAllPostSortByLike(c echo.Context) error {
	//check if page exist
	pageStr := c.QueryParam("page")
	var page int
	if pageStr == "" {
		page = 1
	} else {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	search := c.QueryParam("search")

	posts, numberOfPage, err := h.IPostServices.GetAllPostByLike(page, search)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Success",
		"data":         posts,
		"numberOfPage": numberOfPage,
		"page":         page,
	})
}

func (h *PostHandler) SuspendPost(c echo.Context) error {
	//check if page exist
	postId, errAtoi := strconv.Atoi(c.Param("postId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	err := h.IPostServices.SuspendPost(token, postId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
	})
}
