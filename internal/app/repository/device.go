package repository

import (
	"fmt"
	"lab2/internal/app/ds"

	"time"

	"github.com/sirupsen/logrus"

	//"database/sql"
	"errors"

	"gorm.io/gorm"
)

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

func (r *Repository) GetCurrentCount() int64 {
	var CurrentID uint
	var count int64
	creatorID := 1

	err := r.db.Model(&ds.Current{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("current_id").First(&CurrentID).Error
	if err != nil {
		return 0
	}

	err = r.db.Model(&ds.CurrentDevices{}).Where("current_id = ?", CurrentID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_chats:", err)
	}

	return count
}

func (r *Repository) GetActiveCurrentID() uint {
	var CurrentID uint
	err := r.db.Model(&ds.Current{}).Where("status = ?", "черновик").Select("current_id").First(&CurrentID).Error
	if err != nil {
		return 0
	}
	return CurrentID
}

func (r *Repository) GetCurrent(id int) ([]ds.CurrentDevices, error) {
	var currentItems []ds.CurrentDevices
	err := r.db.Where("current_id = ?", id).Preload("Device").Find(&currentItems).Error
	if err != nil {
		return nil, err
	}

	return currentItems, nil
}

func (r *Repository) AddDevice(deviceID uint, creatorID uint) error {
	var app ds.Current

	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").
		First(&app).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		app = ds.Current{
			Status:       "черновик",
			Created_At:   time.Now(),
			Creator_ID:   creatorID,
			Moderator_ID: 2,
		}
		if err := r.db.Create(&app).Error; err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	var count int64
	r.db.Model(&ds.CurrentDevices{}).
		Where("current_id = ? AND device_id = ?", app.Current_ID, deviceID).Preload("Device").
		Count(&count)

	if count == 0 {
		var device ds.Device
		if err := r.db.First(&device, deviceID).Error; err != nil {
			return err
		}

		appDev := ds.CurrentDevices{
			Current_ID: app.Current_ID,
			Device_ID:  deviceID,
			Amperage:   0,
			Amount:     1,
		}
		if err := r.db.Create(&appDev).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) DeleteCurrent(appID uint) error {
	device_query := `
		UPDATE currents 
		SET status = 'удалён'
		WHERE current_id = $1;
	`
	result := r.db.Exec(device_query, appID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("current with id %d not found", appID)
	}
	return nil
}

func (r *Repository) IsDraftCurrent(appID int) (bool, error) {
	var app ds.Current
	err := r.db.Select("status").Where("current_id = ?", appID).First(&app).Error
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

    device_query := "SELECT device_id, name, type, power_nominal, resistance, voltage_nominal, voltage_bord, coeff_reserve, coeff_efficiency, COALESCE(current_required, 0) as current_required, description, image, in_stock FROM devices WHERE device_id = $1 AND is_delete = false AND in_stock = true"
    // Создание курсора (строковый указатель)
    row := r.db.Raw(device_query, id).Row()

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
	var current ds.Current
	var count int64
	creatorID := 1

	// Находим заявку по creator_id и статусу
	err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").First(&current).Error
	if err != nil {
	   return 0
	}

	// Считаем устройства в заявке
	err = r.db.Model(&ds.CurrentDevices{}).Where("current_id = ?", current.Current_ID).Count(&count).Error
	if err != nil {
	   logrus.Println("Error counting records in application_devices:", err)
	}

	return count
 }

*/
