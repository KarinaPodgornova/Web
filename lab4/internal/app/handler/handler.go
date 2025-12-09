package handler

import (
	"errors"
	"lab4/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

// RegisterHandler godoc
// @title Amperage Calculation API
// @version 1.0
// @description API для управления расчётами силы тока
// @contact.name API Support
// @contact.url http://localhost
// @contact.email support@current.com
// @license.name MIT
// @host localhost
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (h *Handler) RegisterHandler(router *gin.Engine) {
	//Devices
	api := router.Group("/api/v1")
	//Currents
	unauthorized := api.Group("/")
	unauthorized.POST("/users/signup", h.CreateUser)
	unauthorized.POST("/users/signin", h.SignIn)
	unauthorized.GET("/devices", h.GetDevices)
	unauthorized.GET("/devices/:id", h.GetDevice)
	//unauthorized.GET("/current-calculations/current-cart", h.GetCurrentCart)
	unauthorized.PUT("/current-calculations/:id/device_amperage", h.UpdateDeviceAmperage)

	optionalauthorized := api.Group("/")
	optionalauthorized.Use(h.WithOptionalAuthCheck())
	optionalauthorized.GET("/current-calculations/current-cart", h.GetCurrentCart)
	//M:M
	authorized := api.Group("/")
	authorized.Use(h.ModeratorMiddleware(false))
	authorized.POST("/devices/create-device", h.CreateDevice)
	authorized.PUT("/devices/:id/edit-device", h.EditDevice)
	authorized.DELETE("/devices/:id/delete-device", h.DeleteDevice)
	authorized.POST("/devices/:id/add-to-current-calculation", h.AddToCurrent)
	authorized.POST("/devices/:id/image", h.AddPhoto)
	// Users

	authorized.GET("current-calculations/current-calculations", h.GetAllCurrents)
	authorized.GET("/current-calculations/:id", h.GetCurrent)
	authorized.PUT("/current-calculations/:id/edit-current-calculations", h.EditCurrent)
	authorized.PUT("/current-calculations/:id/form", h.FormCurrent)
	//authorized.PUT("/current-calculations/:id/finish", h.FinishCurrent)
	authorized.DELETE("/current-calculations/:id/delete-current-calculations", h.DeleteCurrent)

	authorized.DELETE("/current-devices/:current_id/:device_id", h.DeleteDeviceFromCurrent)
	authorized.PUT("/current-devices/:current_id/:device_id", h.EditDeviceFromCurrent)

	authorized.GET("/users/:login/me", h.GetInfo)
	authorized.PUT("/users/:login/me", h.EditInfo)
	authorized.POST("/users/signout", h.SignOut)

	moderator := api.Group("/")
	moderator.Use(h.ModeratorMiddleware(true))
	//moderator.PUT("/current-calculation/:id/form", h.FormCurrent)
	moderator.PUT("/current-calculations/:id/finish", h.FinishCurrent)

	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
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