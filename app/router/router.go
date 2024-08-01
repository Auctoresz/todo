package router

import (
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"todo-echo/app/handler"
	"todo-echo/app/utils"
)

func Routers(e *echo.Echo, taskHandler handler.TaskHandler, userHandler handler.UserHandler, sessionManager *scs.SessionManager) {
	// LoadAndSave() middleware,
	// takes care of loading and committing session data to the session store,
	// and communicating the session token to/from the client in a cookie as necessary.

	// middleware -> router -> handler (general)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echo.WrapMiddleware(sessionManager.LoadAndSave))

	// router -> middleware -> handler (specific)
	e.GET("/todo", taskHandler.GetAll, utils.RequiredAuth)
	e.POST("/todo", taskHandler.Create, utils.RequiredAuth)
	e.GET("/todo/:taskId", taskHandler.Get, utils.RequiredAuth)
	e.DELETE("/todo/:taskId", taskHandler.Delete, utils.RequiredAuth)
	e.PATCH("/todo/:taskId", taskHandler.Completed, utils.RequiredAuth)
	e.PUT("/todo/:taskId", taskHandler.Edit, utils.RequiredAuth)

	e.POST("/user/signup", userHandler.Signup)
	e.POST("/user/login", userHandler.Login)
	e.POST("/user/logout", userHandler.Logout, utils.RequiredAuth)

}
