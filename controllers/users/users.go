package users

import (
	"discusiin/helper"
	"discusiin/models"
	"discusiin/services/users"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	users.IUserServices
}

func (h *UserHandler) Register(c echo.Context) error {
	// validation
	var u models.User

	errBind := c.Bind(&u)
	if errBind != nil {
		return errBind
	}

	// isEmailKosong?
	if u.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email should not be empty")
	}
	// isEmailValid?
	valid := helper.IsEmailValid(u.Email)
	if !valid {
		return echo.NewHTTPError(http.StatusBadRequest, "email invalid")
	}
	// isUsernameKosong?
	if u.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username should not be empty")
	}
	// isPasswordKosong?
	if u.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "password should not be empty")
	}
	// isPasswordLessThan8?
	if len(u.Password) < 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "password should not lower than 8")
	}

	err := h.IUserServices.Register(u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Created",
	})
}

func (h *UserHandler) RegisterAdmin(c echo.Context) error {
	// validation
	var u models.User

	errBind := c.Bind(&u)
	if errBind != nil {
		return errBind
	}

	// isEmailKosong?
	if u.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email should not be empty")
	}
	// isEmailValid?
	valid := helper.IsEmailValid(u.Email)
	if !valid {
		return echo.NewHTTPError(http.StatusBadRequest, "email invalid")
	}
	// isUsernameKosong?
	if u.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username should not be empty")
	}
	// isPasswordKosong?
	if u.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "password should not be empty")
	}
	// isPasswordLessThan8?
	if len(u.Password) < 8 {
		return echo.NewHTTPError(http.StatusBadRequest, "password should not lower than 8")
	}

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	err := h.IUserServices.RegisterAdmin(u, token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Created",
	})
}

func (h *UserHandler) Login(c echo.Context) error {

	var u models.User
	err := c.Bind(&u)
	if err != nil {
		return err
	}

	if u.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email should not be empty")
	}
	// isEmailValid?
	valid := helper.IsEmailValid(u.Email)
	if !valid {
		return echo.NewHTTPError(http.StatusBadRequest, "email invalid")
	}
	if u.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "password should not be empty")
	}

	result, err := h.IUserServices.Login(u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"data":    result,
	})
}

func (h *UserHandler) GetUsersAdminNotIncluded(c echo.Context) error {
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}
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

	result, numberOfPage, err := h.IUserServices.GetUsersAdminNotIncluded(token, page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Success",
		"data":         result,
		"numberOfPage": numberOfPage,
		"page":         page,
	})
}

func (h *UserHandler) GetProfile(c echo.Context) error {
	var user models.User
	errBind := c.Bind(&user)
	if errBind != nil {
		return errBind
	}

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	result, err := h.IUserServices.GetProfile(token, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
		"data":    result,
	})
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	// validation
	var u models.User

	errBind := c.Bind(&u)
	if errBind != nil {
		return errBind
	}

	// isUsernameKosong?
	if u.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username should not be empty")
	}

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	err := h.IUserServices.UpdateProfile(token, u)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Profile Updated",
	})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	userID, errAtoi := strconv.Atoi(c.Param("userId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	err := h.IUserServices.DeleteUser(token, userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User deleted",
	})
}

func (h *UserHandler) GetPostByUserIdForAdmin(c echo.Context) error {
	userId, errAtoi := strconv.Atoi(c.Param("userId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

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

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	user, result, numberOfPage, err := h.IUserServices.GetPostAsAdmin(token, userId, page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Success",
		"user":         user,
		"data":         result,
		"numberOfPage": numberOfPage,
		"page":         page,
	})
}

func (h *UserHandler) GetCommentByUserIdForAdmin(c echo.Context) error {
	userId, errAtoi := strconv.Atoi(c.Param("userId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

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

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	user, result, numberOfPage, err := h.IUserServices.GetCommentAsAdmin(token, userId, page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Success",
		"user":         user,
		"data":         result,
		"numberOfPage": numberOfPage,
		"page":         page,
	})
}

func (h *UserHandler) GetPostByUserIdAsUser(c echo.Context) error {
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

	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	result, numberOfPage, err := h.IUserServices.GetPostAsUser(token, page)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Success",
		"data":         result,
		"numberOfPage": numberOfPage,
		"page":         page,
	})
}

func (h *UserHandler) BanUser(c echo.Context) error {
	var user models.User

	//bind
	errBind := c.Bind(&user)
	if errBind != nil {
		return errBind
	}

	//get user id from param
	userId, errAtoi := strconv.Atoi(c.Param("userId"))
	if errAtoi != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errAtoi.Error())
	}

	//get token from logged user
	token, errDecodeJWT := helper.DecodeJWT(c)
	if errDecodeJWT != nil {
		return errDecodeJWT
	}

	//ban user
	userBanned, err := h.IUserServices.BanUser(token, userId, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
		"data":    userBanned,
	})
}
