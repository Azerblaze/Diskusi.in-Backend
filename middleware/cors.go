package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"
)

var CustomCORSConfigs = middleware.CORSConfig{
	Skipper:      middleware.DefaultSkipper,
	AllowOrigins: []string{"*"},
	AllowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodPost,
		http.MethodDelete,
	},
}
