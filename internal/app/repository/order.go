package repository

import (
	"fmt"
	"lab2/internal/app/ds"
	"github.com/sirupsen/logrus"
	"math"
	"time"
	//"database/sql"
	"errors"
	"gorm.io/gorm"

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



func (r *Repository) GetDevices() ([]ds.Device, error) {
	var devices []ds.Device
	err := r.db.Find(&devices).Error
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return devices, nil
}

func (r *Repository) GetDevice(id int) (ds.Device, error) {
	device := ds.Device{}
	err := r.db.Where("device_id = ?", id).Find(&device).Error
	if err != nil {
		return ds.Device{}, err
	}
	return device, nil
}


func (r *Repository) GetDevicesByTitle(title string) ([]ds.Device, error) {
	var devices []ds.Device
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *Repository) GetApplicationCount() int64 {
	var ApplicationID uint
	var count int64
	creatorID := 1
 
	err := r.db.Model(&ds.Application{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("application_id").First(&ApplicationID).Error
	if err != nil {
	   return 0
	}
 
	err = r.db.Model(&ds.ApplicationDevices{}).Where("application_id = ?", ApplicationID).Count(&count).Error
	if err != nil {
	   logrus.Println("Error counting records in lists_chats:", err)
	}
 
	return count
 }

 func (r *Repository) GetActiveApplicationID() uint {
	var ApplicationID uint
	err := r.db.Model(&ds.Application{}).Where("status = ?", "черновик").Select("application_id").First(&ApplicationID).Error
	if err != nil {
		return 0
	}
	return ApplicationID
}

func (r *Repository) GetApplication(id int) ([]ds.ApplicationDevices, error) {
    var applicationItems []ds.ApplicationDevices
    err := r.db.Where("application_id = ?", id).Preload("Device").Find(&applicationItems).Error
    if err != nil {
        return nil, err
    }

    return applicationItems, nil
}

func (r *Repository) AddDevice(deviceID uint, creatorID uint) (error) {
    var app ds.Application

    err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").
        First(&app).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        app = ds.Application{
            Status:     "черновик",
            Created_At: time.Now(),
            Creator_ID: creatorID,
			Moderator_ID: 2,
        }
        if err := r.db.Create(&app).Error; err != nil {
            return err
        }
    } else if err != nil {
        return err
    }

    var count int64
    r.db.Model(&ds.ApplicationDevices{}).
        Where("application_id = ? AND device_id = ?", app.Application_ID, deviceID).Preload("Device").
        Count(&count)

    if count == 0 {
		var device ds.Device
		  if err := r.db.First(&device, deviceID).Error; err != nil {
            return err
        }

        appDev := ds.ApplicationDevices{
            Application_ID: app.Application_ID,
            Device_ID:      deviceID,
			Amperage: 		0,
            Amount:         1,
        }
        if err := r.db.Create(&appDev).Error; err != nil {
            return err
        }
    }

    return nil
}


func (r *Repository) DeleteApplication(appID uint) error {
	query := `
		UPDATE applications 
		SET status = 'удалён'
		WHERE application_id = $1;
	`
	result := r.db.Exec(query, appID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("application with id %d not found", appID)
	}
	return nil
}

func (r *Repository) IsDraftApplication(appID int) (bool, error) {
	var app ds.Application
	err := r.db.Select("status").Where("application_id = ?", appID).First(&app).Error
	if err != nil {
		return false, err
	}
	return app.Status == "черновик", nil
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

/*func (r *Repository) GetDeviceByID(id int) (*ds.Device, error) {
    
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

*/
 