package serializer

import "lab3/internal/app/ds"

// CurrentDeviceJSON представляет связь заявки и устройства в формате JSON
type CurrentDeviceJSON struct {
	CurrDev_ID  uint    `json:"curr_dev_id"`   // Идентификатор связи
	Current_ID  uint    `json:"current_id"`    // Идентификатор заявки
	Device_ID   uint    `json:"device_id"`     // Идентификатор устройства
	Amount      int     `json:"amount"`        // Количество устройств
	VoltageBord float64 `json:"voltage_bord"`  // Напряжение бортовой сети
	Amperage    float64 `json:"amperage"`      // Сила тока
}

// CurrentDeviceToJSON преобразует ds.CurrentDevices в CurrentDeviceJSON
func CurrentDeviceToJSON(app_dev ds.CurrentDevices) CurrentDeviceJSON {
	return CurrentDeviceJSON{
		CurrDev_ID:  app_dev.CurrDev_ID,
		Current_ID:  app_dev.Current_ID,
		Device_ID:   app_dev.Device_ID,
		Amount:      app_dev.Amount,
		VoltageBord: app_dev.VoltageBord,
		Amperage:    app_dev.Amperage,
	}
}

// CurrentDeviceFromJSON преобразует CurrentDeviceJSON в ds.CurrentDevices
func CurrentDeviceFromJSON(deviceJSON CurrentDeviceJSON) ds.CurrentDevices {
	return ds.CurrentDevices{
		CurrDev_ID:  deviceJSON.CurrDev_ID,
		Current_ID:  deviceJSON.Current_ID,
		Device_ID:   deviceJSON.Device_ID,
		Amount:      deviceJSON.Amount,
		VoltageBord: deviceJSON.VoltageBord,
		Amperage:    deviceJSON.Amperage,
	}
}