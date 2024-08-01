package handler

import (
	"errors"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"todo-echo/app/domain/dao"
	"todo-echo/app/domain/dto"
	"todo-echo/app/repository"
)

type TaskHandler interface {
	GetAll(c echo.Context) error
	Get(c echo.Context) error
	Create(c echo.Context) error
	Delete(c echo.Context) error
	Completed(c echo.Context) error
	Edit(c echo.Context) error
}

type TaskHandlerImpl struct {
	TaskRepository repository.TaskRepository
	Validators     *validator.Validate
	SessionManager *scs.SessionManager
}

// NewTaskHandlerImpl is a constructor to ensure all field has been injected
func NewTaskHandlerImpl(taskRepository repository.TaskRepository, validators *validator.Validate, sessionManager *scs.SessionManager) *TaskHandlerImpl {
	return &TaskHandlerImpl{
		TaskRepository: taskRepository,
		Validators:     validators,
		SessionManager: sessionManager,
	}
}

func (t TaskHandlerImpl) GetAll(c echo.Context) error {

	// User id is stored in the context
	userId := c.Get("userId").(int)

	// DAO
	taskDao := dao.Task{
		IdCustomer: int64(userId),
	}

	// Repository
	tasks, err := t.TaskRepository.GetAll(c.Request().Context(), taskDao)
	// error handling
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, "Database conflict")
	}

	// Response (JSON)
	response := dto.TaskResponse{
		Code:   200,
		Status: "Ok",
		Data:   tasks,
	}
	return c.JSON(http.StatusOK, response)
}

func (t TaskHandlerImpl) Get(c echo.Context) error {
	// Request (Path Parameter)
	param := c.Param("taskId")
	id, _ := strconv.ParseInt(param, 10, 64)

	// DTO
	taskDto := dto.TaskGetRequest{Id: id}

	// Validate
	err := t.Validators.Struct(taskDto)
	// error handling
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// DAO
	taskDao := dao.Task{Id: taskDto.Id}

	// Repository
	task, err := t.TaskRepository.Get(c.Request().Context(), taskDao)
	// error handling
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Task not found")
	}

	// Response (JSON)
	flash := t.SessionManager.PopString(c.Request().Context(), "flash")
	if flash == "" {
		flash = "No flash message"
	}

	response := dto.TaskResponse{
		Code:    200,
		Message: flash,
		Status:  "Ok",
		Data:    task,
	}

	return c.JSON(http.StatusOK, response)
}

func (t TaskHandlerImpl) Create(c echo.Context) error {
	// User id is stored in the context
	userId := c.Get("userId").(int)
	// Request (body) and DTO
	taskDto := new(dto.TaskCreateRequest)
	err := c.Bind(taskDto)
	// error handling
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	}

	// Validate
	t.Validators.Struct(taskDto)

	// DAO
	taskDao := dao.Task{IdCustomer: int64(userId), Description: taskDto.Description, DueDate: taskDto.DueDate}

	// Repository
	task, err := t.TaskRepository.Insert(c.Request().Context(), taskDao)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, "Database conflict")
	}

	// Redirect
	t.SessionManager.Put(c.Request().Context(), "flash", "Snippet successfully created!")

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/todo/%d", task.Id))
}

func (t TaskHandlerImpl) Delete(e echo.Context) error {
	userId := e.Get("userId").(int)

	// Path parameter
	param := e.Param("taskId")
	id, _ := strconv.ParseInt(param, 10, 64)

	// Repository
	taskDao := dao.Task{Id: id, IdCustomer: int64(userId)}
	err := t.TaskRepository.Delete(e.Request().Context(), taskDao)
	if err != nil {
		if errors.Is(err, dao.ErrNoRecord) {
			return e.JSON(http.StatusNotFound, map[string]string{"message": "Task not found!"})
		} else {
			return e.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}

	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Delete success!"})
}

func (t TaskHandlerImpl) Completed(e echo.Context) error {
	// User Id
	userId := e.Get("userId").(int)

	// Bind JSON
	taskDto := new(dto.TaskCompletedRequest)
	e.Bind(taskDto)

	// Path parameter
	param := e.Param("taskId")
	id, _ := strconv.ParseInt(param, 10, 64)

	// Repository
	taskDao := dao.Task{Id: id, IdCustomer: int64(userId), IsCompleted: taskDto.IsCompleted}
	err := t.TaskRepository.Completed(e.Request().Context(), taskDao)
	if err != nil {
		if errors.Is(err, dao.ErrNoRecord) {
			return e.JSON(http.StatusNotFound, map[string]string{"message": "Task not found!"})
		} else {
			return e.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}

	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Completed!"})
}

func (t TaskHandlerImpl) Edit(e echo.Context) error {
	// Get user id
	userId := e.Get("userId").(int)

	// Bind JSON
	taskDto := new(dto.TaskCreateRequest)
	e.Bind(taskDto)

	// Path parameter
	param := e.Param("taskId")
	id, _ := strconv.ParseInt(param, 10, 64)

	// Validate
	err := t.Validators.Struct(taskDto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	// Repository
	taskDao := dao.Task{Id: id, IdCustomer: int64(userId), Description: taskDto.Description, DueDate: taskDto.DueDate}
	err = t.TaskRepository.Edit(e.Request().Context(), taskDao)
	if err != nil {
		if errors.Is(err, dao.ErrNoRecord) {
			return e.JSON(http.StatusNotFound, map[string]string{"message": "Task not found!"})
		} else {
			return e.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
		}

	}

	return e.JSON(http.StatusOK, map[string]string{"message": "Edited!"})
}
