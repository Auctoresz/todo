package main

import (
	"crypto/tls"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
	"todo-echo/app/handler"
	"todo-echo/app/repository"
	"todo-echo/app/router"
	"todo-echo/app/utils"
	"todo-echo/config"
)

func main() {
	// CONFIG
	// Echo
	e := echo.New()
	// Database
	Db := config.OpenDB()
	defer Db.Close()
	// Web Security (TLS)
	cert, err := tls.LoadX509KeyPair("localhost.crt", "localhost.key")
	if err != nil {
		log.Fatal(err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// SECONDARY
	// flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	//logger
	errorLog := utils.ErrorLog()
	infoLog := utils.InfoLog()
	// validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	_ = validate.RegisterValidation("ISO8601date", utils.IsISO8601Date)
	// session manager (middleware)
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(Db)
	sessionManager.Lifetime = 12 * time.Hour
	// utils
	utils.SetSessionManager(sessionManager)

	// PRIMARY
	// repository
	taskRepository := repository.NewTaskRepositoryImpl(Db)
	userRepository := repository.NewUserRepositoryImpl(Db)
	// handler
	taskHandler := handler.NewTaskHandlerImpl(taskRepository, validate, sessionManager)
	userHandler := handler.NewUserHandlerImpl(userRepository, validate, sessionManager)
	// router
	router.Routers(e, taskHandler, userHandler, sessionManager)
	// server
	server := http.Server{
		Addr:      *addr,
		Handler:   e,
		ErrorLog:  errorLog,
		TLSConfig: tlsConfig,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServeTLS("", "")
	errorLog.Fatal(err)

}
