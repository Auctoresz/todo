package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func RequiredAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		if !IsAuthenticated(e) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Please login")
		}
		// get session data (user id) from session store (mysql)
		userId := UserIdAuthenticated(e)
		// store it in the context
		e.Set("userId", userId)
		// Otherwise set the "Cache-Control: no-store" header so that pages
		// require authentication are not stored in the users browser cache (or
		// other intermediary cache).
		e.Response().Header().Set("Cache-Control", "no-store")
		return next(e)
	}
}
