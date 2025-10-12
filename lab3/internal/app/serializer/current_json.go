package serializer

import (
	"database/sql"
	"lab3/internal/app/ds"
	"time"
)

// CurrentJSON представляет заявку на расчет силы тока в формате JSON
type CurrentJSON struct {
	ID              uint       `json:"current_id"`      // Идентификатор заявки
	Status          string     `json:"status"`          // Статус заявки
	Created_At      time.Time  `json:"created_at"`      // Дата создания
	Creator_Login   string     `json:"creator_login"`   // Логин создателя
	Moderator_Login *string    `json:"moderator_login"` // Логин модератора (опционально)
	Forming_Date    *time.Time `json:"form_date"`       // Дата формирования (опционально)
	Finish_Date     *time.Time `json:"finish_date"`     // Дата завершения (опционально)
	VoltageBord     float64    `json:"voltage_bord"`    // Бортовое напряжение для расчета
}

// StatusJSON представляет статус для обновления заявки
type StatusJSON struct {
	Status string `json:"status"` // Статус заявки
}

// CurrentToJSON преобразует структуру ds.Current в JSON-формат с учетом логинов
func CurrentToJSON(current ds.Current, creator_login string, moderator_login string) CurrentJSON {
	var form_date, finish_date *time.Time
	// Если Forming_Date валиден, создаем указатель на время
	if current.Forming_Date.Valid {
		form_date = &current.Forming_Date.Time
	}
	// Если Finish_Date валиден, создаем указатель на время
	if current.Finish_Date.Valid {
		finish_date = &current.Finish_Date.Time
	}
	var m_login *string
	// Если moderator_login не пустой, создаем указатель на него
	if moderator_login != "" {
		m_login = &moderator_login
	}

	return CurrentJSON{
		ID:              current.Current_ID,
		Status:          current.Status,
		Created_At:      current.Created_At,
		Creator_Login:   creator_login,
		Moderator_Login: m_login,
		Forming_Date:    form_date,
		Finish_Date:     finish_date,
		VoltageBord:     current.VoltageBord,
	}
}

// CurrentFromJSON преобразует JSON-данные в структуру ds.Current
func CurrentFromJSON(currentJSON CurrentJSON) ds.Current {
	// Инициализируем структуру с учетом всех полей
	current := ds.Current{
		Current_ID:  currentJSON.ID,
		Status:      currentJSON.Status,
		Created_At:  currentJSON.Created_At,
		VoltageBord: currentJSON.VoltageBord,
	}
	// Преобразуем Forming_Date из *time.Time в sql.NullTime
	if currentJSON.Forming_Date != nil {
		current.Forming_Date = sql.NullTime{Time: *currentJSON.Forming_Date, Valid: true}
	}
	// Преобразуем Finish_Date из *time.Time в sql.NullTime
	if currentJSON.Finish_Date != nil {
		current.Finish_Date = sql.NullTime{Time: *currentJSON.Finish_Date, Valid: true}
	}
	// Creator_ID и Moderator_ID остаются неизменными, так как логины не преобразуются в ID напрямую
	return current
}