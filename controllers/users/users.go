package users

import (
	"discusiin/helper"
	"discusiin/models"
	"discusiin/services/users"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	users.IUserServices
}

func (h *UserHandler) Register(c echo.Context) error {
	var u models.User
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	// isEmailKosong?
	if u.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "email should not be empty",
		})
	}
	// isEmailValid?
	valid := helper.IsEmailValid(u.Email)
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "email invalid",
		})
	}
	// isUsernameKosong?
	if u.Username == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "username should not be empty",
		})
	}
	// isPasswordKosong?
	if u.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "password should not be empty",
		})
	}
	// isPasswordLessThan8?
	if len(u.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "password can not less than 8",
		})
	}
	// isUsernameLessThan3?
	if len(u.Username) < 3 {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "username can not less than 3",
		})
	}
	// isPasswordContainNumber
	if !helper.IsStringContainNumber(u.Password) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "password should be contain number",
		})
	}
	// isPasswordContainUppercaseLetter
	if !helper.IsStringContainUppercaseLetter(u.Password) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "password should be contain uppercase letter",
		})
	}

	err = h.IUserServices.Register(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "created",
	})
}

func (h *UserHandler) Login(c echo.Context) error {

	var u models.User
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if u.Email == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "email should not be empty",
		})
	}
	// isEmailValid?
	valid := helper.IsEmailValid(u.Email)
	if !valid {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "email invalid",
		})
	}
	if u.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "password should not be empty",
		})
	}

	result, err := h.IUserServices.Login(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    result,
	})
}
