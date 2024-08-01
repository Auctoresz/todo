package handler

import (
	"errors"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"todo-echo/app/domain/dao"
	"todo-echo/app/domain/dto"
	"todo-echo/app/repository"
)

type UserHandler interface {
	Login(e echo.Context) error
	Signup(e echo.Context) error
	Logout(e echo.Context) error
}

type UserHandlerImpl struct {
	UserRepository repository.UserRepository
	Validators     *validator.Validate
	SessionManager *scs.SessionManager
}

func NewUserHandlerImpl(userRepository repository.UserRepository, validators *validator.Validate, sessionManager *scs.SessionManager) UserHandler {
	return &UserHandlerImpl{UserRepository: userRepository, Validators: validators, SessionManager: sessionManager}
}

func (u UserHandlerImpl) Login(e echo.Context) error {
	// Read from request body and Bind JSON
	userDto := new(dto.UserLoginRequest)
	e.Bind(userDto)

	// Validate
	err := u.Validators.Struct(userDto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// Repository
	userDao := dao.User{Email: userDto.Email, Password: []byte(userDto.Password)}
	userId, err := u.UserRepository.Authenticate(e.Request().Context(), userDao)
	if err != nil {
		if errors.Is(err, dao.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusNotAcceptable, "Email or password incorrect")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// privilege levels changes for the user (e.g. login
	// and logout operations).
	err = u.SessionManager.RenewToken(e.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	// Add session data (userId)
	u.SessionManager.Put(e.Request().Context(), "authenticatedUserId", int(userId))

	// Redirect
	return e.Redirect(http.StatusSeeOther, "/todo")

}

func (u UserHandlerImpl) Signup(e echo.Context) error {
	// Request (body) and DTO
	userDto := new(dto.UserSignupRequest)
	e.Bind(userDto)

	// Validate
	err := u.Validators.Struct(userDto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	// DAO
	userDao := dao.User{
		Email:    userDto.Email,
		Password: []byte(userDto.Password),
		Name:     userDto.Name,
	}

	// Repository
	err = u.UserRepository.Insert(e.Request().Context(), userDao)
	if err != nil {
		if errors.Is(err, dao.ErrDuplicateEmail) {
			return echo.NewHTTPError(http.StatusNotAcceptable, "Email already used")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "Random error on server")
		}

	}

	// Response (redirect)
	return e.String(http.StatusOK, "Signup success")
}

func (u UserHandlerImpl) Logout(e echo.Context) error {
	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// privilege levels changes for the user (e.g. login
	// and logout operations).
	err := u.SessionManager.RenewToken(e.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	u.SessionManager.Remove(e.Request().Context(), "authenticatedUserId")

	return e.String(http.StatusOK, "Logout")
}
