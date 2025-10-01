package handler

import (
	"lab1/internal/app/repository"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func (h *Handler) GetDevices(ctx *gin.Context) {
	var devices []repository.Device
	var err error

	searchQuery := ctx.Query("device_query")
	if searchQuery == "" {
		devices, err = h.Repository.GetDevices()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		devices, err = h.Repository.GetDevicesByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	// Получаем количество товаров в текущей заявке
	currentCurrentID := h.Repository.GetCurrentCurrentID()
	cartCount := h.Repository.GetCurrentDevicesCount(currentCurrentID)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":         time.Now().Format("15:04:05"),
		"devices":      devices,
		"cartCount":    cartCount,
		"device_query": searchQuery,
	})
}

func (h *Handler) GetDevice(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	device, err := h.Repository.GetDevice(id)
	if err != nil {
		logrus.Error(err)
	}

	specsArray := strings.Split(device.Specs, "\n")

	currentCurrentID := h.Repository.GetCurrentCurrentID()
	cartCount := h.Repository.GetCurrentDevicesCount(currentCurrentID)

	ctx.HTML(http.StatusOK, "device.html", gin.H{
		"device":     device,
		"specsArray": specsArray,
		"cartCount":  cartCount,
	})
}

func (h *Handler) GetCurrent(ctx *gin.Context) {
	currentCurrentID := h.Repository.GetCurrentCurrentID()
	current := h.Repository.GetCurrent(currentCurrentID)
	devicesInCurrent := h.Repository.GetCurrentDevices(currentCurrentID)

	// Создаем мапу для передачи в шаблон с дополнительной информацией
	currentItems := make(map[int]repository.CurrentDevice)
	for _, item := range current.DeviceItems {
		currentItems[item.DeviceID] = item
	}

	ctx.HTML(http.StatusOK, "current.html", gin.H{
		"current":      current,
		"devices":      devicesInCurrent,
		"currentItems": currentItems,
		"cartCount":    len(devicesInCurrent),
	})
}

func (h *Handler) GetCart(ctx *gin.Context) {
	// Теперь корзина - это текущая заявка
	currentCurrentID := h.Repository.GetCurrentCurrentID()
	current := h.Repository.GetCurrent(currentCurrentID)
	devicesInCurrent := h.Repository.GetCurrentDevices(currentCurrentID)

	currentItems := make(map[int]repository.CurrentDevice)
	for _, item := range current.DeviceItems {
		currentItems[item.DeviceID] = item
	}

	ctx.HTML(http.StatusOK, "current.html", gin.H{
		"current":      current,
		"devices":      devicesInCurrent,
		"currentItems": currentItems,
		"cartCount":    len(devicesInCurrent),
	})
}

func (h *Handler) GetAllCurrents(ctx *gin.Context) {
	currents := h.Repository.GetAllCurrents()

	ctx.HTML(http.StatusOK, "currents.html", gin.H{
		"currents": currents,
	})
}
