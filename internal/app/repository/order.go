package repository

import (
	"fmt"
	"lab2/internal/app/ds"
	"github.com/sirupsen/logrus"
	"math"
	"database/sql"
	"errors"

)


// ДОБАВЬ ЭТОТ МЕТОД:
func (r *Repository) CalculateRequiredCurrent(device *ds.Device) float64 {
	// L_требуемая = (√(P_ном / R_ном)) * (K_запаса / (K_пд * (U_борт / U_ном)))
	
	// Защита от деления на ноль
	if device.Resistance == 0 || device.CoeffEfficiency == 0 || device.VoltageBord == 0 {
		return 0
	}
	
	// √(P_ном / R_ном)
	sqrtPart := math.Sqrt(device.PowerNominal / device.Resistance)
	
	// (U_борт / U_ном)
	voltageRatio := device.VoltageBord / device.VoltageNominal
	
	// K_пд * (U_борт / U_ном)
	efficiencyPart := device.CoeffEfficiency * voltageRatio
	
	// K_запаса / (K_пд * (U_борт / U_ном))
	reservePart := device.CoeffReserve / efficiencyPart
	
	// Итоговый расчет
	requiredCurrent := sqrtPart * reservePart
	
	return math.Round(requiredCurrent*100) / 100 // Округляем до 2 знаков
}



// Переименовываем методы под устройства
func (r *Repository) GetDevices() ([]ds.Device, error) {
	var devices []ds.Device
	err := r.db.Where("is_delete = false").Find(&devices).Error // добавили условие
	if err != nil {
		return nil, err
	}
	return devices, nil
}






/*
func (r *Repository) GetDevice(id int) (ds.Device, error) {
	device := ds.Device{}
	err := r.db.Where("device_id = ? AND in_stock = ?", id, true).First(&device).Error
	if err != nil {
		return ds.Device{}, err
	}
	return device, nil
}
*/

func (r *Repository) GetDeviceByID(id int) (*ds.Device, error) {
    
    query := "SELECT device_id, name, type, power_nominal, resistance, voltage_nominal, voltage_bord, coeff_reserve, coeff_efficiency, COALESCE(current_required, 0) as current_required, description, image, in_stock FROM devices WHERE device_id = $1 AND is_delete = false AND in_stock = true"
    // Создание курсора (строковый указатель)
    row := r.db.Raw(query, id).Row()

    // Создание объекта для хранения данных
    device := &ds.Device{}

    // Сканирование строки в структуру
    err := row.Scan(
        &device.Device_ID,
        &device.Name,
        &device.Type,
        &device.PowerNominal,
        &device.Resistance,
        &device.VoltageNominal,
        &device.VoltageBord,
        &device.CoeffReserve,
        &device.CoeffEfficiency,
        &device.CurrentRequired,
        &device.Description,
        &device.Image,
        &device.InStock,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // Возвращаем nil, если записи нет
        }
        return nil, err
    }

    return device, nil
}






func (r *Repository) GetDevicesByName(name string) ([]ds.Device, error) {
	var devices []ds.Device
	err := r.db.Where("name ILIKE ? AND in_stock = ?", "%"+name+"%", true).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

// Метод для расчета силы тока
func (r *Repository) CalculateCurrent(power, voltage float64) float64 {
	if voltage == 0 {
		return 0
	}
	return power / voltage
}




func (r *Repository)  SearchDevicesByName(name string) ([]ds.Device, error) {
	var devices []ds.Device
	err := r.db.Where("name ILIKE ? and is_delete = ?", "%"+name+"%", false).Find(&devices).Error // добавили условие
	if err != nil {
		return nil, err
	}
	return devices, nil
}



// GetCartCount для получения количества устройств в заявке
func (r *Repository) GetCartCount() int64 {
	var application ds.Application
	var count int64
	creatorID := 1
 
	// Находим заявку по creator_id и статусу
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&application).Error
	if err != nil {
	   return 0
	}
 
	// Считаем устройства в заявке
	err = r.db.Model(&ds.ApplicationDevices{}).Where("application_id = ?", application.Application_ID).Count(&count).Error
	if err != nil {
	   logrus.Println("Error counting records in application_devices:", err)
	}
 
	return count
 }


 
func (r *Repository) DeleteDevice(deviceID uint) error {
    err := r.db.Model(&ds.Device{}).Where("device_id = ?", deviceID).UpdateColumn("is_delete", true).Error
    fmt.Println(deviceID)
    if err != nil {
        return fmt.Errorf("ошибка при удалении устройства с id %d: %w", deviceID, err)
    }
    return nil
}

