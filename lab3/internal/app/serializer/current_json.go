package serializer

import (
	"database/sql"
	"lab3/internal/app/ds"
	"time"
)

// CurrentJSON представляет заявку на расчет силы тока в формате JSON
type CurrentJSON struct {
	ID             uint       `json:"current_id"` // Идентификатор заявки
	Status         string     `json:"status"`                  // Статус заявки
	Created_At     time.Time  `json:"created_at"`              // Дата создания
	Creator_Login  string     `json:"creator_login"`           // Логин создателя
	Moderator_Login *string   `json:"moderator_login"`         // Логин модератора (опционально)
	Forming_Date   *time.Time `json:"form_date"`               // Дата формирования (опционально)
	Finish_Date    *time.Time `json:"finish_date"`             // Дата завершения (опционально)
	Amperage       float64    `json:"amperage"`                // Сила тока
}

// StatusJSON представляет статус для обновления заявки
type StatusJSON struct {
	Status string `json:"status"` // Статус заявки
}

// CurrentToJSON преобразует структуру ds.Current в JSON-формат с учетом логинов
func CurrentToJSON(app ds.Current, creator_login string, moderator_login string) CurrentJSON {
	var form_date, finish_date *time.Time
	// Если Forming_Date валиден, создаем указатель на время
	if app.Forming_Date.Valid {
		form_date = &app.Forming_Date.Time
	}
	// Если Finish_Date валиден, создаем указатель на время
	if app.Finish_Date.Valid {
		finish_date = &app.Finish_Date.Time
	}
	var m_login *string
	// Если moderator_login не пустой, создаем указатель на него
	if moderator_login != "" {
		m_login = &moderator_login
	}

	return CurrentJSON{
		ID:             app.Current_ID,
		Status:         app.Status,
		Created_At:     app.Created_At,
		Creator_Login:  creator_login,
		Moderator_Login: m_login,
		Forming_Date:   form_date,
		Finish_Date:    finish_date,
		Amperage:       app.Amperage,
	}
}

// CurrentFromJSON преобразует JSON-данные в структуру ds.Current
func CurrentFromJSON(appJSON CurrentJSON) ds.Current {
	// Если Amperage равен 0, возвращаем пустую структуру
	if appJSON.Amperage == 0 {
		return ds.Current{}
	}
	// Инициализируем структуру с учетом всех полей
	current := ds.Current{
		Current_ID: appJSON.ID,
		Status:     appJSON.Status,
		Created_At: appJSON.Created_At,
		Amperage:   appJSON.Amperage,
	}
	// Преобразуем Forming_Date из *time.Time в sql.NullTime
	if appJSON.Forming_Date != nil {
		current.Forming_Date = sql.NullTime{Time: *appJSON.Forming_Date, Valid: true}
	}
	// Преобразуем Finish_Date из *time.Time в sql.NullTime
	if appJSON.Finish_Date != nil {
		current.Finish_Date = sql.NullTime{Time: *appJSON.Finish_Date, Valid: true}
	}
	// Creator_ID и Moderator_ID остаются неизменными, так как логины не преобразуются в ID напрямую
	return current
}