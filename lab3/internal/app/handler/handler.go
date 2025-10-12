package handler

import (
	"errors"
	"lab3/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/api/devices", h.GetDevices)  // Список с фильтром
	router.GET("/api/devices/:id", h.GetDevice)
	router.POST("/api/devices", h.CreateDevice)  // Без -create-device
	router.PUT("/api/devices/:id", h.EditDevice)  // Без -edit-device
	router.DELETE("/api/devices/:id", h.DeleteDevice)  // Без -delete-device
	router.POST("/api/devices/:id/add-to-calculation", h.AddToCurrent)  // По теме: add to calculation (draft)
	router.POST("/api/devices/:id/image", h.AddPhoto)  // Отдельный для изображения

	router.GET("/api/calculations/cart", h.GetCurrentCart)  // Иконка корзины: id draft + count
	router.GET("/api/calculations", h.GetAllCurrents)  // Список с фильтром по forming_date и status
	router.GET("/api/calculations/:id", h.GetCurrent)
	router.PUT("/api/calculations/:id", h.EditCurrent)  // Изменение тематических полей
	router.PUT("/api/calculations/:id/form", h.FormCurrent)  // Формировать (creator)
	router.PUT("/api/calculations/:id/finish", h.FinishCurrent)  // Завершить/отклонить (moderator, с расчётом)
	router.DELETE("/api/calculations/:id", h.DeleteCurrent)  //

	router.DELETE("/api/calculation-devices/:current_id/:device_id", h.DeleteDeviceFromCurrent)
	router.PUT("/api/calculation-devices/:current_id/:device_id", h.EditDeviceFromCurrent)
	// Users
	router.POST("/api/users", h.CreateUser)  // Регистрация
	router.GET("/api/users/me", h.GetInfo)  // После auth
	router.PUT("/api/users/me", h.EditInfo)
	router.POST("/api/users/login", h.SignIn)
	router.POST("/api/users/logout", h.SignOut)

}



func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())

	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}

	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}


