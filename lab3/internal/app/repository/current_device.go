package repository

import (
	"errors"
	"fmt"
	"lab3/internal/app/ds"
	"lab3/internal/app/serializer"

	"gorm.io/gorm"
)

func (r *Repository) DeleteDeviceFromCurrent(current_id int, device_id int) (ds.Current, error) {
	// userId := r.userId
	// if userId == 0 {
	//     return ds.Research{}, fmt.Errorf("%w: пользователь не авторизирован", ErrNotAllowed)
	// }

	// user, err := r.GetUserByID(userId)
	// if err != nil {
	// 	return ds.Research{}, err
	// }

	var current ds.Current
	err := r.db.Where("current_id = ?", current_id).First(&current).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Current{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, current_id)
		}
		return ds.Current{}, err
	}

	// if research.CreatorID != r.userId && !user.IsModerator{
	// 	return ds.Research{}, fmt.Errorf("%w: Вы не создатель этого исследования", ErrNotAllowed)
	// }

	err = r.db.Where("device_id = ? and current_id = ?", device_id, current_id).Delete(&ds.CurrentDevices{}).Error
	if err != nil {
		return ds.Current{}, err
	}
	return current, nil
}

func (r *Repository) EditDeviceFromCurrent(current_id int, device_id int, currentDeviceJSON serializer.CurrentDeviceJSON) (ds.CurrentDevices, error) {
	var currentsDevice ds.CurrentDevices
	err := r.db.Model(&currentsDevice).Where("device_id = ? and current_id = ?", device_id, current_id).Updates(serializer.CurrentDeviceFromJSON(currentDeviceJSON)).First(&currentsDevice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.CurrentDevices{}, fmt.Errorf("%w: устройства в заявке", ErrNotFound)
		}
		return ds.CurrentDevices{}, err
	}
	return currentsDevice, nil
}
