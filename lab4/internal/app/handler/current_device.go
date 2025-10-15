// current_device.go
package handler

import (
	"errors"
	"lab4/internal/app/ds"
	"lab4/internal/app/repository"
	"lab4/internal/app/serializer"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteDeviceFromCurrent godoc
// @Summary Удалить устройство из заявки
// @Description Удаляет связь устройства и заявки
// @Tags current_devices
// @Produce json
// @Param device_id path int true "ID устройства"
// @Param current_id path int true "ID заявки"
// @Success 200 {object} serializer.CurrentJSON "Обновленная заявка"
// @Failure 400 {object} map[string]string "Неверные ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /current-devices/{current_id}/{device_id} [delete]


// calculateAmperage - функция расчета силы тока для устройства в заявке
// calculateAmperage - ПРАВИЛЬНАЯ функция расчета силы тока по формуле
func (h *Handler) calculateAmperage(device ds.Device, current ds.Current, amount int) float64 {
	// Базовые значения по умолчанию
	voltageBord := current.VoltageBord
	if voltageBord == 0 {
		voltageBord = 11.5 // минимальное напряжение бортовой сети
	}

	power := device.PowerNominal
	resistance := device.Resistance
	voltageNominal := device.VoltageNominal
	efficiency := device.CoeffEfficiency
	reserve := device.CoeffReserve

	// Проверка на нулевые значения
	if resistance == 0 || voltageNominal == 0 || efficiency == 0 || reserve == 0 {
		return 0
	}

	// ПРАВИЛЬНАЯ ФОРМУЛА:
	// I_требуемая = √(P_ном / R_ном) * (K_запаса / (K_пд * (U_борт / U_ном)))

	// 1. Вычисляем √(P_ном / R_ном)
	part1 := math.Sqrt(power / resistance)

	// 2. Вычисляем (U_борт / U_ном)
	voltageRatio := voltageBord / voltageNominal

	// 3. Вычисляем (K_пд * (U_борт / U_ном))
	denominator := efficiency * voltageRatio

	// 4. Вычисляем (K_запаса / denominator)
	part2 := reserve / denominator

	// 5. Итоговая сила тока для одного устройства
	amperagePerDevice := part1 * part2

	// 6. Умножаем на количество устройств
	return amperagePerDevice * float64(amount)
}

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
		"status":  "success",
		"message": "Устройство удалено из заявки",
	})
}

// EditDeviceFromCurrent godoc
// @Summary Изменить данные устройства в заявке
// @Description Обновляет параметры устройства в конкретной заявке
// @Tags current_devices
// @Accept json
// @Produce json
// @Param device_id path int true "ID устройства"
// @Param current_id path int true "ID заявки"
// @Param data body serializer.CurrentDeviceJSON true "Новые данные"
// @Success 200 {object} serializer.CurrentDeviceJSON "Обновленные данные"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router  /current-devices/{current_id}/{device_id} [put]

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

	// 1. Получаем текущую связь с предзагрузкой Device и Current
	var currentDevice ds.CurrentDevices
	err = h.Repository.DB().Preload("Device").Preload("Current").
		Where("current_id = ? AND device_id = ?", current_id, device_id).
		First(&currentDevice).Error

	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound, errors.New("связь устройства с заявкой не найдена"))
		return
	}

	// 2. Обновляем количество
	currentDevice.Amount = currentDeviceJSON.Amount

	// 3. ← РАССЧИТЫВАЕМ Amperage ТОЛЬКО ДЛЯ ЗАВЕРШЕННЫХ ЗАЯВОК
	if currentDevice.Current.Status == "completed" {
		currentDevice.Amperage = h.calculateAmperage(
			currentDevice.Device,
			currentDevice.Current,
			currentDeviceJSON.Amount,
		)
	} else {
		// ← ДЛЯ ВСЕХ ДРУГИХ СТАТУСОВ - сила тока = 0
		currentDevice.Amperage = 0
	}

	// 4. Сохраняем обновленную связь
	if err := h.Repository.DB().Save(&currentDevice).Error; err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// 5. Если заявка завершена, пересчитываем общую силу тока
	if currentDevice.Current.Status == "completed" {
		err = h.Repository.RecalculateCurrentAmperage(uint(current_id))
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusOK, serializer.CurrentDeviceToJSON(currentDevice))
}
