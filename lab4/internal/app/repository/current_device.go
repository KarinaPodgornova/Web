package repository

import (
	"errors"
	"fmt"
	"lab4/internal/app/ds"
	"lab4/internal/app/serializer"

	"gorm.io/gorm"
)

func (r *Repository) DeleteDeviceFromCurrent(current_id int, device_id int) error {
	// Проверяем существование заявки
	var current ds.Current
	err := r.db.Where("current_id = ?", current_id).First(&current).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: заявка с id %d", ErrNotFound, current_id)
		}
		return err
	}

	// Удаляем связь
	result := r.db.Where("device_id = ? and current_id = ?", device_id, current_id).Delete(&ds.CurrentDevices{})
	if result.Error != nil {
		return result.Error
	}

	// Проверяем, была ли удалена хотя бы одна запись
	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: устройство %d не найдено в заявке %d", ErrNotFound, device_id, current_id)
	}

	return nil
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