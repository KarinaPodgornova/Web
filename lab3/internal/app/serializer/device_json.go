package serializer

import "lab3/internal/app/ds"

// DeviceJSON представляет устройство в формате JSON
type DeviceJSON struct {
	Device_ID       uint    `json:"device_id"`       // Идентификатор устройства
	Name            string  `json:"name"`            // Название устройства
	Type            string  `json:"type"`            // Тип устройства
	PowerNominal    float64 `json:"power_nominal"`   // Номинальная мощность
	Resistance      float64 `json:"resistance"`      // Сопротивление
	VoltageNominal  float64 `json:"voltage_nominal"` // Номинальное напряжение
	CoeffReserve    float64 `json:"coeff_reserve"`   // Коэффициент резерва
	CoeffEfficiency float64 `json:"coeff_efficiency"`// Коэффициент эффективности
	CurrentRequired float64 `json:"current_required"`// Требуемая сила тока
	Description     string  `json:"description"`     // Описание устройства
	Image           string  `json:"image"`           // URL изображения
	InStock         bool    `json:"in_stock"`        // Наличие на складе
	IsDelete        bool    `json:"is_delete"`       // Флаг удаления
}

// DeviceToJSON преобразует ds.Device в DeviceJSON
func DeviceToJSON(device ds.Device) DeviceJSON {
	return DeviceJSON{
		Device_ID:       device.Device_ID,
		Name:            device.Name,
		Type:            device.Type,
		PowerNominal:    device.PowerNominal,
		Resistance:      device.Resistance,
		VoltageNominal:  device.VoltageNominal,
		CoeffReserve:    device.CoeffReserve,
		CoeffEfficiency: device.CoeffEfficiency,
		CurrentRequired: device.CurrentRequired,
		Description:     device.Description,
		Image:           device.Image,
		InStock:         device.InStock,
		IsDelete:        device.IsDelete,
	}
}

// DeviceFromJSON преобразует DeviceJSON в ds.Device
func DeviceFromJSON(deviceJSON DeviceJSON) ds.Device {
	return ds.Device{
		Device_ID:       deviceJSON.Device_ID,
		Name:            deviceJSON.Name,
		Type:            deviceJSON.Type,
		PowerNominal:    deviceJSON.PowerNominal,
		Resistance:      deviceJSON.Resistance,
		VoltageNominal:  deviceJSON.VoltageNominal,
		CoeffReserve:    deviceJSON.CoeffReserve,
		CoeffEfficiency: deviceJSON.CoeffEfficiency,
		CurrentRequired: deviceJSON.CurrentRequired,
		Description:     deviceJSON.Description,
		Image:           deviceJSON.Image,
		InStock:         deviceJSON.InStock,
		IsDelete:        deviceJSON.IsDelete,
	}
}