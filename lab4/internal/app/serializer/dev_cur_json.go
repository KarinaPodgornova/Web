package serializer

import "lab4/internal/app/ds"

// CurrentDeviceJSON представляет связь заявки и устройства в формате JSON
type CurrentDeviceJSON struct {
	CurrDev_ID  uint    `json:"curr_dev_id"`   
	Current_ID  uint    `json:"current_id"`    
	Device_ID   uint    `json:"device_id"`     
	Amount      int     `json:"amount"`        
	Amperage    float64 `json:"amperage"`      
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

// CurrentDevicesArrayToJSON преобразует массив ds.CurrentDevices в массив CurrentDeviceJSON
func CurrentDevicesArrayToJSON(currentDevices []ds.CurrentDevices) []CurrentDeviceJSON {
	result := make([]CurrentDeviceJSON, 0, len(currentDevices))
	for _, cd := range currentDevices {
		result = append(result, CurrentDeviceToJSON(cd))
	}
	return result
}