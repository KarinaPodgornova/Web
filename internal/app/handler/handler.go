package handler

import (
  "github.com/gin-gonic/gin"
  "github.com/sirupsen/logrus"
  "lab2/internal/app/repository"
)

type Handler struct {
  Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
  return &Handler{
    Repository: r,
  }
}

// RegisterHandler Функция, в которой мы отдельно регистрируем маршруты
func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/", h.GetDevices)
	router.GET("/device/:id", h.GetDevice)
	router.GET("/request", h.GetRequest)  
	router.GET("/cart", h.GetCart)       
  router.GET("/devices", h.GetAllDevices)
  router.POST("/delete-device", h.DeleteDevice)

  

}

// RegisterStatic регистрируем статику
func (h *Handler) RegisterStatic(router *gin.Engine) {
  router.LoadHTMLGlob("templates/*")
  
  // ДОБАВЬ эту строку для стилей из internal/resources/styles/
  //router.Static("/styles", "./internal/resources/styles")
  
  // Или если у тебя есть другие статические файлы в resources:
  router.Static("/static", "./resources")
}

// errorHandler для более удобного вывода ошибок 
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}