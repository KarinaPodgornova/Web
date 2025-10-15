package handler

import (
	"errors"
	"fmt"
	"lab4/internal/app/ds"
	"lab4/internal/app/repository"
	"lab4/internal/app/serializer"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAllCurrents godoc
// @Summary Получить список заявок на расчёт
// @Description Возвращает заявки с возможностью фильтрации по датам и статусу
// @Tags Currents
// @Produce json
// @Param from-date query string false "Начальная дата (YYYY-MM-DD)"
// @Param to-date query string false "Конечная дата (YYYY-MM-DD)"
// @Param status query string false "Статус заявки"
// @Success 200 {array} serializer.CurrentJSON "Список заявок"
// @Failure 400 {object} map[string]string "Неверный формат даты"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /current-calculations/current-calculations [get]	
func (h *Handler) GetAllCurrents(ctx *gin.Context) {
	fromDate := ctx.Query("from")
	var from = time.Time{}
	var to = time.Time{}
	if fromDate != "" {
		from1, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = from1
	}
	fmt.Println(fromDate)

	toDate := ctx.Query("to")
	if toDate != "" {
		to1, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = to1
	}

	status := ctx.Query("status")

	currents, err := h.Repository.GetAllCurrents(from, to, status)
	if err != nil {
		fmt.Printf("Ошибка GetAllCurrents: %v\n", err)
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	currents = h.filterAuthorizedCurrents(currents, ctx)
	resp := make([]serializer.CurrentJSON, 0, len(currents))
	for _, c := range currents {
		creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(c)
		if err != nil {
			fmt.Printf("Ошибка GetModeratorAndCreatorLogin: %v\n", err)
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		resp = append(resp, serializer.CurrentToJSON(c, creatorLogin, moderatorLogin))
	}
	ctx.JSON(http.StatusOK, resp)
	
}

// GetCurrentCart godoc
// @Summary Получить корзину расчёта
// @Description Возвращает информацию о текущей заявке-черновике на расчёт пользователя
// @Tags Currents
// @Produce json
// @Success 200 {object} map[string]interface{} "Данные корзины заявки-черновика"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /current-calculations/current-cart [get]
func (h *Handler) GetCurrentCart(ctx *gin.Context) {
	
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	devices_count := h.Repository.GetCurrentCount(userID)

	if devices_count == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":        "no_draft",
			"devices_count": devices_count,
		})
		return
	}

	current, err := h.Repository.CheckCurrentCurrentDraft(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
		} else if errors.Is(err, repository.ErrNoDraft) {
			ctx.JSON(http.StatusOK, gin.H{
				"status":        "no_draft",
				"devices_count": 0,
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":            current.Current_ID,
		"devices_count": h.Repository.GetCurrentCount(current.Creator_ID),
	})
}

// GetCurrent godoc
// @Summary Получить заявку по ID
// @Description Возвращает полную информацию о заявке
// @Tags Currents
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]interface{} "Данные заявки с устройствами"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /current-calculations/{id} [get]
func (h *Handler) GetCurrent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	devices, current, err := h.Repository.GetCurrentDevices(id)
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

	resp := make([]serializer.DeviceJSON, 0, len(devices))
	for _, r := range devices {
		resp = append(resp, serializer.DeviceToJSON(r))
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(current)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	currentDevices, _ := h.Repository.GetDevicesCurrents(int(current.Current_ID))

	resp2 := make([]serializer.CurrentDeviceJSON, 0, len(currentDevices))
	for _, r := range currentDevices {
		resp2 = append(resp2, serializer.CurrentDeviceToJSON(r))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"current":        serializer.CurrentToJSON(current, creatorLogin, moderatorLogin),
		"devices":        resp,
		"currentDevices": resp2,
	})
}

// FormCurrent godoc
// @Summary Сформировать заявку
// @Description Переводит заявку в статус "formed"
// @Tags Currents
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} serializer.CurrentJSON "Сформированная заявка"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /current-calculations/{id}/form [put]
func (h *Handler) FormCurrent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "formed"

	current, err := h.Repository.FormCurrent(id, status)
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

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(current)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, serializer.CurrentToJSON(current, creatorLogin, moderatorLogin))
}

// EditCurrent godoc
// @Summary Изменить заявка
// @Description Обновляет данные заявки
// @Tags Currents
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param current body serializer.CurrentJSON true "Новые данные заявки"
// @Success 200 {object} serializer.CurrentJSON "Обновленная заявка"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /current-calculations/{id}/edit-current-calculations [put]
func (h *Handler) EditCurrent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var currentJSON serializer.CurrentJSON
	if err := ctx.BindJSON(&currentJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	current, err := h.Repository.EditCurrent(id, currentJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(current)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, serializer.CurrentToJSON(current, creatorLogin, moderatorLogin))
}

// DeleteCurrent godoc
// @Summary Удалить заявка
// @Description Выполняет логическое удаление заявки
// @Tags Currents
// @Produce json
// @Param id path int true "ID заявки"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /current-calculations/{id}/delete-current-calculations [delete]
func (h *Handler) DeleteCurrent(ctx *gin.Context) {
	idStr := ctx.Param("id")
	current_id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "deleted"

	_, err = h.Repository.FormCurrent(current_id, status)
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Current deleted"})
}

// FinishCurrent godoc
// @Summary Завершить заявку
// @Description Изменяет статус заявки (только для модераторов)
// @Tags Currents
// @Accept json
// @Produce json
// @Param id path int true "ID заявки"
// @Param status body serializer.StatusJSON true "Новый статус"
// @Success 200 {object} serializer.CurrentJSON "Результат модерации"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Заявка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /current-calculations/{id}/finish [put]
func (h *Handler) FinishCurrent(ctx *gin.Context) {

	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var statusJSON serializer.StatusJSON
	if err := ctx.BindJSON(&statusJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	if !user.IsModerator {
		h.errorHandler(ctx, http.StatusForbidden, errors.New("требуются права модератора"))
		return
	}

	current, err := h.Repository.FinishCurrent(id, statusJSON.Status, userID)

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

	currentDevices, err := h.Repository.GetDevicesCurrents(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	//РАССЧИТЫВАЕМ ОБЩУЮ СИЛУ ТОКА
	var totalAmperage float64
	if statusJSON.Status == "completed" {
		for _, device := range currentDevices {
			totalAmperage += device.Amperage
		}
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(current)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"current":               serializer.CurrentToJSON(current, creatorLogin, moderatorLogin),
		"devices_with_amperage": serializer.CurrentDevicesArrayToJSON(currentDevices),
		"total_amperage":        totalAmperage,
	})

}

func (h *Handler) filterAuthorizedCurrents(currents []ds.Current, ctx *gin.Context) []ds.Current {
	userID, err := getUserID(ctx)
	if err != nil {
		return []ds.Current{}
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return []ds.Current{}
	}
	if err != nil {
		return []ds.Current{}
	}

	if user.IsModerator {
		return currents
	}

	var userCurrents []ds.Current
	for _, Current := range currents {
		fmt.Println(Current.Current_ID)
		if Current.Creator_ID == userID {
			userCurrents = append(userCurrents, Current)
		}
	}

	return userCurrents

}

func (h *Handler) hasAccessToCurrent(creatorID uuid.UUID, ctx *gin.Context) bool {
	userID, err := getUserID(ctx)
	if err != nil {
		return false
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return false
	}
	if err != nil {
		return false
	}

	return creatorID == userID || user.IsModerator
}
