package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"lab3/internal/app/ds"
	"lab3/internal/app/serializer"
	"time"
	//"math"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//var errNoDraft = errors.New("no draft for this user")

func (r *Repository) GetAllCurrents(from, to time.Time, status string) ([]ds.Current, error) {
	var currents []ds.Current
	sub := r.db.Where("status != 'deleted' and status != 'draft'")
	if !from.IsZero() {
		sub = sub.Where("forming_date > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("created_at < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("current_id").Find(&currents).Error
	if err != nil {
		return nil, err
	}
	return currents, nil
}

func (r *Repository) GetDevicesCurrents(current_id int) ([]ds.CurrentDevices, error) {
	var currentDevice []ds.CurrentDevices
	err := r.db.Where("current_id = ?", current_id).Find(&currentDevice).Error
	if err != nil {
		return nil, err
	}
	return currentDevice, nil
}

func (r *Repository) GetDevicesCurrent(device_id int, current_id int) (ds.CurrentDevices, error) {
	var currentDevice ds.CurrentDevices
	err := r.db.Where("device_id = ? and current_id = ?", device_id, current_id).First(&currentDevice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.CurrentDevices{}, fmt.Errorf("%w: device current not found", ErrNotFound)
		}
		return ds.CurrentDevices{}, err
	}
	return currentDevice, nil
}

func (r *Repository) GetCurrentDevices(id int) ([]ds.Device, ds.Current, error) {
	current, err := r.GetSingleCurrent(id)
	if err != nil {
		return []ds.Device{}, ds.Current{}, err
	}

	var devices []ds.Device
	sub := r.db.Table("current_devices").Where("current_id = ?", current.Current_ID)
	err = r.db.Order("device_id DESC").Where("device_id IN (?)", sub.Select("device_id")).Find(&devices).Error

	if err != nil {
		return []ds.Device{}, ds.Current{}, err
	}

	return devices, current, nil
}

func (r *Repository) CheckCurrentCurrentDraft(creator_ID uint) (ds.Current, error) {
	// if creatorID == 0 {
	//     return ds.Research{}, fmt.Errorf("%w: user not authenticated", ErrNotAllowed)
	// }

	var current ds.Current
	res := r.db.Where("creator_id = ? AND status = ?", creator_ID, "draft").Limit(1).Find(&current)
	if res.Error != nil {
		return ds.Current{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.Current{}, ErrNoDraft
	}
	return current, nil
}

func (r *Repository) GetCurrentDraft(creator_ID uint) (ds.Current, bool, error) {
	// if creatorID == 0 {
	//     return ds.Research{}, false, fmt.Errorf("%w: user not authenticated", ErrNotAllowed)
	// }

	current, err := r.CheckCurrentCurrentDraft(creator_ID)
	if errors.Is(err, ErrNoDraft) {
		current = ds.Current{
			Status:     "draft",
			Creator_ID: creator_ID,
			
			Created_At: time.Now(),
			VoltageBord: 11.5,
		}
		result := r.db.Create(&current)
		if result.Error != nil {
			return ds.Current{}, false, result.Error
		}
		return current, true, nil
	} else if err != nil {
		return ds.Current{}, false, err
	}
	return current, true, nil
}

func (r *Repository) GetCurrentCount(creator_ID uint) int64 {
	if creator_ID == 0 {
		return 0
	}

	var count int64
	current, err := r.CheckCurrentCurrentDraft(creator_ID)
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.CurrentDevices{}).Where("current_id = ?", current.Current_ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_devices:", err)
	}

	return count
}

func (r *Repository) DeleteCalculation(current_id int) error {
	return r.db.Exec("UPDATE currents SET status = 'deleted' WHERE id = ?", current_id).Error
}

func (r *Repository) GetSingleCurrent(id int) (ds.Current, error) {
	if id < 0 {
		return ds.Current{}, errors.New("неверное id, должно быть >= 0")
	}

	// userId := r.GetUserID()
	// if userId == 0 {
	//     return ds.Research{}, fmt.Errorf("%w: пользователь не авторизирован", ErrNotAllowed)
	// }

	// user, err := r.GetUserByID(userId)
	// if err != nil {
	// 	return ds.Research{}, err
	// }

	var current ds.Current
	err := r.db.Where("current_id = ?", id).First(&current).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Current{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.Current{}, err
	} else if current.Status == "deleted" {
		return ds.Current{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	return current, nil
}

func (r *Repository) FormCurrent(current_id int, status string) (ds.Current, error) {
	current, err := r.GetSingleCurrent(current_id)
	if err != nil {
		return ds.Current{}, err
	}

	if current.Status != "draft" {
		return ds.Current{}, fmt.Errorf("эта заявка не может быть %s", status)
	}

	if status != "deleted" {
		// Проверяем, что есть устройства в заявке
		currentDevices, err := r.GetDevicesCurrents(int(current.Current_ID))
		if err != nil {
			return ds.Current{}, err
		}
		if len(currentDevices) == 0 {
			return ds.Current{}, errors.New("нельзя сформировать пустую заявку")
		}
		
		// Проверяем корректность напряжения бортовой сети
		if current.VoltageBord <= 0 {
			return ds.Current{}, errors.New("вы ввели некорректное напряжение бортовой сети")
		}
	}

	err = r.db.Model(&current).Updates(ds.Current{
		Status: status,
		Forming_Date: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.Current{}, err
	}

	return current, nil
}


func (r *Repository) EditCurrent(id int, currentJSON serializer.CurrentJSON) (ds.Current, error) {
	current := ds.Current{}
	if id < 0 {
		return ds.Current{}, errors.New("неправильное id, должно быть >= 0")
	}
	if currentJSON.VoltageBord <= 0 {
		return ds.Current{}, errors.New("неправильное напряжение бортовой сети")
	}
	
	err := r.db.Where("current_id = ? and status != 'deleted'", id).First(&current).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Current{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.Current{}, err
	}
	
	// Обновляем VoltageBord
	err = r.db.Model(&current).Update("voltage_bord", currentJSON.VoltageBord).Error
	if err != nil {
		return ds.Current{}, err
	}
	
	// Пересчитываем силу тока для всех устройств в заявке
	if current.Status == "completed" {
		err = r.RecalculateCurrentAmperage(current.Current_ID)
		if err != nil {
			return ds.Current{}, err
		}
	}
	
	err = r.db.First(&current, id).Error
	return current, err
}


// CalculateDeviceCurrent рассчитывает силу тока для одного устройства
func (r *Repository) CalculateDeviceCurrent(device *ds.Device, voltageBord float64, amount int) (float64, error) {
	if device.PowerNominal <= 0 || device.Resistance <= 0 || device.VoltageNominal <= 0 || 
	   voltageBord <= 0 || device.CoeffReserve <= 0 || device.CoeffEfficiency <= 0 {
		return 0, errors.New("неверные параметры для расчёта тока")
	}
	
	// I = (P_Hom / R_Hom) * (K_3anaca / (K_ng * (U_6opT / U_Hom)))
	part1 := device.PowerNominal / device.Resistance
	part2 := voltageBord / device.VoltageNominal
	part3 := device.CoeffReserve / (device.CoeffEfficiency * part2)
	
	amperagePerDevice := part1 * part3
	return amperagePerDevice * float64(amount), nil
}



// RecalculateCurrentAmperage пересчитывает силу тока для всех устройств в заявке
func (r *Repository) RecalculateCurrentAmperage(currentID uint) error {
	current, err := r.GetSingleCurrent(int(currentID))
	if err != nil {
		return err
	}

	currentDevices, err := r.GetDevicesCurrents(int(currentID))
	if err != nil {
		return err
	}

	for _, currentDevice := range currentDevices {
		device, err := r.GetDevice(int(currentDevice.Device_ID))
		if err != nil {
			return err
		}
		
		// ПЕРЕДАЕМ УКАЗАТЕЛЬ НА УСТРОЙСТВО!
		amperage, err := r.CalculateDeviceCurrent(device, current.VoltageBord, currentDevice.Amount)
		if err != nil {
			return err
		}
		
		err = r.db.Model(&currentDevice).Update("amperage", amperage).Error
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (r *Repository) FinishCurrent(id int, status string) (ds.Current, error) {
	if status != "completed" && status != "rejected" {
		return ds.Current{}, errors.New("неверный статус")
	}

	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.Current{}, err
	}

	if !user.IsModerator {
		return ds.Current{}, fmt.Errorf("%w: вы не модератор", ErrNotAllowed)
	}

	current, err := r.GetSingleCurrent(id)
	if err != nil {
		return ds.Current{}, err
	} else if current.Status != "formed" {
		return ds.Current{}, fmt.Errorf("эта заявка не может быть %s", status)
	}

	// Обновляем через map чтобы избежать проблем с типами
	moderatorID := user.User_ID
	updates := map[string]interface{}{
		"status": status,
		"finish_date": sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		"moderator_id": &moderatorID,
	}

	err = r.db.Model(&current).Updates(updates).Error
	if err != nil {
		return ds.Current{}, err
	}

	// Если заявка завершена, пересчитываем силу тока
	if status == "completed" {
		err = r.RecalculateCurrentAmperage(current.Current_ID)
		if err != nil {
			return ds.Current{}, err
		}
	}
	
	return current, nil
}