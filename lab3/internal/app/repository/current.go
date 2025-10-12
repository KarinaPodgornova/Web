package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"lab3/internal/app/ds"
	"lab3/internal/app/serializer"
	"time"
	"math"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var errNoDraft = errors.New("no draft for this user")

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

	// user, err := r.GetUserByID(r.GetUserID())
	// if err != nil{
	// 	return ds.Research{}, fmt.Errorf("%w: пользователь на авторизирован", ErrNotAllowed)
	// }

	// if research.CreatorID != r.userId && !user.IsModerator{
	// 	return ds.Research{}, fmt.Errorf("%w: у вас нет прав чтобы эта заявка имела статус %s", ErrNotAllowed, status)
	// }

	if current.Status != "draft" {
		return ds.Current{}, fmt.Errorf("эта заявка не может быть %s", status)
	}

	if status != "deleted" {
		if current.Amperage < 0 {
			return ds.Current{}, errors.New("вы не написали нагрузку системы")
		}
		currentDevices, _ := r.GetDevicesCurrents(int(current.Current_ID))
		for _, currentDevices := range currentDevices {
			if currentDevices.VoltageBord < 0 {
				return ds.Current{}, errors.New("вы ввели некорректное напряжение бортовой сети")
			}
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
	if currentJSON.Amperage < 0 {
		return ds.Current{}, errors.New("неправильная нагрузка")
	}
	err := r.db.Where("current_id = ? and status != 'deleted'", id).First(&current).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Current{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.Current{}, err
	}
	err = r.db.Model(&current).Updates(serializer.CurrentFromJSON(currentJSON)).Error
	if err != nil {
		return ds.Current{}, err
	}
	return current, nil
}

func CalculateDeviceCurrent(device *ds.Device, cd *ds.CurrentDevices) (float64, error) {
	if device.PowerNominal <= 0 || device.Resistance <= 0 || device.VoltageNominal <= 0 || cd.VoltageBord <= 0 || device.CoeffReserve <= 0 || device.CoeffEfficiency <= 0 {
		return 0, errors.New("неверные параметры для расчёта тока")
	}
	sqrtTerm := math.Sqrt(device.PowerNominal / device.Resistance)
	denom := device.CoeffEfficiency * (cd.VoltageBord / device.VoltageNominal)
	adjustTerm := device.CoeffReserve / denom
	return sqrtTerm * adjustTerm, nil
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
		return ds.Current{}, fmt.Errorf("это исследование не может быть %s", status)
	}

	err = r.db.Model(&current).Updates(ds.Current{
		Status: status,
		Finish_Date: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Moderator_ID: uint(user.User_ID),
	}).Error
	if err != nil {
		return ds.Current{}, err
	}

	if status == "completed" {
		currentsDevice, err := r.GetDevicesCurrents(int(current.Current_ID))
		if err != nil {
			return ds.Current{}, err
		}
		var totalAmperage float64 = 0
		for _, currentDevice := range currentsDevice {
			device, err := r.GetDevice(int(currentDevice.Device_ID))
			if err != nil {
				return ds.Current{}, err
			}
			deviceAmperage, err := CalculateDeviceCurrent(device, &currentDevice)
			if err != nil {
				return ds.Current{}, err
			}
			amperageWithAmount := deviceAmperage * float64(currentDevice.Amount)
			err = r.db.Model(&currentDevice).Updates(ds.CurrentDevices{
				Amperage: amperageWithAmount,
			}).Error
			if err != nil {
				return ds.Current{}, err
			}
			totalAmperage += amperageWithAmount
		}
		// Сохрани общий ток в заявке
		err = r.db.Model(&current).Update("amperage", totalAmperage).Error
		if err != nil {
			return ds.Current{}, err
		}
	}
	return current, nil

}
