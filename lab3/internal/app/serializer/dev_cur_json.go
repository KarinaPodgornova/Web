package serializer

import "lab3/internal/app/ds"

// CurrentDeviceJSON представляет связь заявки и устройства в формате JSON
type CurrentDeviceJSON struct {
	CurrDev_ID  uint    `json:"curr_dev_id"`   // Идентификатор связи
	Current_ID  uint    `json:"current_id"`    // Идентификатор заявки
	Device_ID   uint    `json:"device_id"`     // Идентификатор устройства
	Amount      int     `json:"amount"`        // Количество устройств
	Amperage    float64 `json:"amperage"`      // Сила тока
}

// CurrentDeviceToJSON преобразует ds.CurrentDevices в CurrentDeviceJSON
func CurrentDeviceToJSON(currentDevice ds.CurrentDevices) CurrentDeviceJSON {
	return CurrentDeviceJSON{
		CurrDev_ID:  currentDevice.CurrDev_ID,
		Current_ID:  currentDevice.Current_ID,
		Device_ID:   currentDevice.Device_ID,
		Amount:      currentDevice.Amount,
		Amperage:    currentDevice.Amperage,
	}
}

// CurrentDeviceFromJSON преобразует CurrentDeviceJSON в ds.CurrentDevices
func CurrentDeviceFromJSON(deviceJSON CurrentDeviceJSON) ds.CurrentDevices {
	return ds.CurrentDevices{
		CurrDev_ID:  deviceJSON.CurrDev_ID,
		Current_ID:  deviceJSON.Current_ID,
		Device_ID:   deviceJSON.Device_ID,
		Amount:      deviceJSON.Amount,
		Amperage:    deviceJSON.Amperage,
	}
}