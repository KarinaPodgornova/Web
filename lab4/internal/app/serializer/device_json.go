package serializer

import "lab4/internal/app/ds"

// DeviceJSON представляет устройство в формате JSON
type DeviceJSON struct {
	Device_ID       uint    `json:"device_id"`      
	Name            string  `json:"name"`            
	Type            string  `json:"type"`            
	PowerNominal    float64 `json:"power_nominal"`   
	Resistance      float64 `json:"resistance"`     
	VoltageNominal  float64 `json:"voltage_nominal"` 
	CoeffReserve    float64 `json:"coeff_reserve"`  
	CoeffEfficiency float64 `json:"coeff_efficiency"`
	CurrentRequired float64 `json:"current_required"`
	Description     string  `json:"description"`    
	Image           string  `json:"image"`           
	InStock         bool    `json:"in_stock"`       
	IsDelete        bool    `json:"is_delete"`      
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