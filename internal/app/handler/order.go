package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"lab2/internal/app/ds"
)

func (h *Handler) GetDevices(ctx *gin.Context) {
	var devices []ds.Device
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		devices, err = h.Repository.GetDevices()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		devices, err = h.Repository.GetDevicesByName(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	// Рассчитываем силу тока для каждого устройства
	for i := range devices {
		if devices[i].CurrentRequired == 0{
			devices[i].CurrentRequired = h.Repository.CalculateRequiredCurrent(&devices[i])
		}
	}

	cartCount := h.Repository.GetCartCount()

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":    time.Now().Format("15:04:05"),
		"devices": devices, // МЕНЯЕМ orders на devices
		"query":   searchQuery,
		"cart_count": cartCount, // ДОБАВЛЯЕМ КОЛИЧЕСТВО В КОРЗИНЕ
	})
}

func (h *Handler) GetDevice(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	//device, err := h.Repository.GetDevice(id)
	device, err := h.Repository.GetDeviceByID(id)
	if err != nil {
		logrus.Error(err)
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// Рассчитываем силу тока если не рассчитана
	//if device.CurrentRequired == 0 {
		//device.CurrentRequired = h.Repository.CalculateRequiredCurrent(&device)
	//}

	specsArray := strings.Split(device.Description, "\n")

	// ДОБАВЛЯЕМ cart_count ДЛЯ СТРАНИЦЫ УСТРОЙСТВА
	cartCount := h.Repository.GetCartCount()

	ctx.HTML(http.StatusOK, "order.html", gin.H{ // МЕНЯЕМ order.html на device.html
		"device": device,
		"specsArray": specsArray, // ДОБАВЬ ЭТУ СТРОКУ
		"cart_count":  cartCount, // ДОБАВЛЯЕМ И ЗДЕСЬ
	})
}

// ДОБАВЬТЕ ЭТИ ДВА МЕТОДА ДЛЯ КОРЗИНЫ:

func (h *Handler) GetCart(ctx *gin.Context) {
    // Здесь будет логика отображения корзины
    cartCount := h.Repository.GetCartCount()
    
    ctx.HTML(http.StatusOK, "cart.html", gin.H{
        "cart_count": cartCount,
    })
}

func (h *Handler) GetRequest(ctx *gin.Context) {
    // Здесь будет логика отображения калькулятора
    cartCount := h.Repository.GetCartCount()
    
    ctx.HTML(http.StatusOK, "request.html", gin.H{
        "cart_count": cartCount,
    })
}

// ДОБАВЛЯЕМ АНАЛОГИ МЕТОДОВ ИЗ CHAT.GO:

// GetAllDevicesPage - аналог GetAllChats, но для страницы со списком устройств
func (h *Handler) GetAllDevices(ctx *gin.Context) {
	var devices []ds.Device
	var err error

	search := ctx.Query("search")
	if search == "" {
		devices, err = h.Repository.GetDevices()
	} else {
		devices, err = h.Repository.GetDevicesByName(search)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	// Рассчитываем силу тока для каждого устройства
	for i := range devices {
		if devices[i].CurrentRequired == 0 {
			devices[i].CurrentRequired = h.Repository.CalculateRequiredCurrent(&devices[i])
		}
	}

	ctx.HTML(http.StatusOK, "devices.html", gin.H{
		"devices":    devices,
		"cart_count": h.Repository.GetCartCount(),
		"search":     search,
	})
}


func (h *Handler) DeleteDevice(ctx *gin.Context) {
	// считываем значение из формы, которую мы добавим в наш шаблон
	strId := ctx.PostForm("device_id")
	id, err := strconv.Atoi(strId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	// Вызов функции добавления чата в заявку
	err = h.Repository.DeleteDevice(uint(id))
	if err != nil && !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return
	}
	
	// после вызова сразу произойдет обновление страницы
	ctx.Redirect(http.StatusFound, "/devices")
}
