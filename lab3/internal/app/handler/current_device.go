// current_device.go
package handler

import (
	"errors"
	"net/http"
	"strconv"
	"lab3/internal/app/repository"
	"lab3/internal/app/serializer"
	"github.com/gin-gonic/gin"
)

func (h *Handler) DeleteDeviceFromCurrent(ctx *gin.Context) {
	current_id, err := strconv.Atoi(ctx.Param("current_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("неверный ID заявки"))
		return
	}

	device_id, err := strconv.Atoi(ctx.Param("device_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("неверный ID устройства"))
		return
	}

	// Удаляем устройство из заявки
	err = h.Repository.DeleteDeviceFromCurrent(current_id, device_id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	// Возвращаем только статус успеха
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "Устройство удалено из заявки",
	})
}

func (h *Handler) EditDeviceFromCurrent(ctx *gin.Context) {
	// Берем параметры из URL
	current_id, err := strconv.Atoi(ctx.Param("current_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("неверный ID заявки"))
		return
	}

	device_id, err := strconv.Atoi(ctx.Param("device_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, errors.New("неверный ID устройства"))
		return
	}

	// Берем только обновляемые поля из body
	var currentDeviceJSON serializer.CurrentDeviceJSON
	if err := ctx.BindJSON(&currentDeviceJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	currentDevice, err := h.Repository.EditDeviceFromCurrent(current_id, device_id, currentDeviceJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, serializer.CurrentDeviceToJSON(currentDevice))
}

