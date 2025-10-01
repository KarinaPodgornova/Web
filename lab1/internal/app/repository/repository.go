package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Device struct {
	ID          int
	Title       string
	Description string
	Specs       string
	Power       string
}

// Структура для связи многие-ко-многим
type CurrentDevice struct {
	DeviceID int
	Quantity int
}

// Основная структура заявки
type Current struct {
	ID          int
	CurrentDate string
	Status      string
	DeviceItems []CurrentDevice
}

func (r *Repository) GetDevices() ([]Device, error) {
	devices := []Device{
		{
			ID:          1,
			Title:       "Универсальный адаптер BMW Step-In",
			Power:       "10 Вт",
			Description: "Вставной адаптер Snap-In предназначен для всех подходящих телефонов с разъемом Micro USB. Помимо зарядки телефона обеспечивает также улучшенный прием и надежное крепление. Благодаря разъему USB записанная на телефоне музыка может удобно воспроизводиться через динамики аудиосистемы BMW.",
			Specs:       "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 10 Вт\nРазмеры: 145 x 72 x 22 мм",
		},
		{
			ID:          2,
			Title:       "Адаптер BMW Step-In для мобильных телефонов Apple",
			Power:       "15 Вт",
			Description: "Специально разработан для устройств Apple. Обеспечивает быструю зарядку и стабильное соединение с медиасистемой BMW.",
			Specs:       "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 15 Вт\nРазмеры: 142 x 70 x 20 мм",
		},
		{
			ID:          3,
			Title:       "Автомобильный отопитель BMW",
			Power:       "150 Вт",
			Description: "Дополнительный отопитель салона с дистанционным управлением. Обеспечивает комфортную температуру в холодное время года.",
			Specs:       "Входное напряжение: 12 В\nПотребляемая мощность: 150 Вт\nТепловая мощность: 2 кВт\nРазмеры: 180 x 120 x 80 мм",
		},
		{
			ID:          4,
			Title:       "Автомобильный отопитель BMW Premium",
			Power:       "200 Вт",
			Description: "Усовершенствованная модель отопителя с цифровым управлением и таймером. Поддерживает заданную температуру автоматически.",
			Specs:       "Входное напряжение: 12 В\nПотребляемая мощность: 200 Вт\nТепловая мощность: 2.5 кВт\nРазмеры: 190 x 125 x 85 мм",
		},
		{
			ID:          5,
			Title:       "Переходник Micro USB",
			Power:       "5 Вт",
			Description: "Компактный переходник для подключения устройств с Micro USB к медиасистеме BMW. Обеспечивает зарядку и передачу данных.",
			Specs:       "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 5 Вт\nРазмеры: 45 x 20 x 15 мм",
		},
		{
			ID:          6,
			Title:       "Адаптер Media для Apple iPod / iPhone",
			Power:       "12 Вт",
			Description: "Адаптер для полной интеграции устройств Apple с медиасистемой BMW. Поддержка управления через iDrive.",
			Specs:       "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 12 Вт\nРазмеры: 95 x 55 x 25 мм",
		},
		{
			ID:          7,
			Title:       "Музыкальный/медийный адаптер BMW для Apple iPod/iPhone",
			Power:       "18 Вт",
			Description: "Премиальный адаптер с поддержкой высококачественного звука и видео. Полная интеграция с интерфейсом iDrive.",
			Specs:       "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 18 Вт\nРазмеры: 100 x 60 x 28 мм",
		},
		{
			ID:          8,
			Title:       "Кабель-адаптер BMW для Micro-USB",
			Power:       "10 Вт",
			Description: "Оригинальный кабель-адаптер для устройств Android. Обеспечивает быструю зарядку и стабильное соединение.",
			Specs:       "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 10 Вт\nДлина кабеля: 1.2 м",
		},
		{
			ID:          9,
			Title:       "CD — привод дополнительный",
			Power:       "25 Вт",
			Description: "Встраиваемый дополнительный CD-привод для BMW. Поддержка аудио CD и MP3-дисков. Интеграция со штатной аудиосистемой.",
			Specs:       "Входное напряжение: 12 В\nПотребляемая мощность: 25 Вт\nФорматы: CD, CD-R, CD-RW, MP3\nРазмеры: 180 x 145 x 45 мм",
		},
		{
			ID:          10,
			Title:       "Аксессуары для ключа-браслета",
			Power:       "2 Вт",
			Description: "Стильный браслет-чехол для ключа BMW. Защита от повреждений и удобное ношение. Совместим со всеми моделями ключей BMW.",
			Specs:       "Материал: силикон\nЦвет: черный\nСовместимость: все ключи BMW\nВес: 15 г",
		},
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return devices, nil
}

func (r *Repository) GetDevice(id int) (Device, error) {
	devices, err := r.GetDevices()
	if err != nil {
		return Device{}, err
	}

	for _, device := range devices {
		if device.ID == id {
			return device, nil
		}
	}
	return Device{}, fmt.Errorf("товар не найден")
}

func (r *Repository) GetDevicesByTitle(searchQuery string) ([]Device, error) {
	devices, err := r.GetDevices()
	if err != nil {
		return []Device{}, err
	}

	var result []Device
	for _, device := range devices {
		if strings.Contains(strings.ToLower(device.Title), strings.ToLower(searchQuery)) {
			result = append(result, device)
		}
	}

	return result, nil
}

// Получение заявки по ID
func (r *Repository) GetCurrent(id int) Current {
	currents := map[int]Current{
		1: {
			ID:          1,
			CurrentDate: "16.09.25",
			Status:      "В обработке",
			DeviceItems: []CurrentDevice{
				{DeviceID: 1, Quantity: 2},
				{DeviceID: 3, Quantity: 1},
				{DeviceID: 5, Quantity: 3},
			},
		},
		2: {
			ID:          2,
			CurrentDate: "17.09.25",
			Status:      "Подтверждена",
			DeviceItems: []CurrentDevice{
				{DeviceID: 2, Quantity: 1},
				{DeviceID: 4, Quantity: 2},
				{DeviceID: 6, Quantity: 1},
			},
		},
		3: {
			ID:          3,
			CurrentDate: "18.09.25",
			Status:      "Выполнена",
			DeviceItems: []CurrentDevice{
				{
					DeviceID: 7,
					Quantity: 1,
				},
				{
					DeviceID: 8,
					Quantity: 2,
				},
				{
					DeviceID: 9,
					Quantity: 1,
				},
				{
					DeviceID: 10,
					Quantity: 4,
				},
			},
		},
	}

	return currents[id]
}

// Получение товаров в заявке
func (r *Repository) GetCurrentDevices(id int) []Device {
	current := r.GetCurrent(id)

	var devicesInCurrent []Device
	for _, currentDevice := range current.DeviceItems {
		device, err := r.GetDevice(currentDevice.DeviceID)
		if err == nil {
			devicesInCurrent = append(devicesInCurrent, device)
		}
	}
	return devicesInCurrent
}

// Получение количества товаров в заявке
func (r *Repository) GetCurrentDevicesCount(id int) int {
	current := r.GetCurrent(id)
	return len(current.DeviceItems)
}

// Получение ID активной заявки
func (r *Repository) GetCurrentCurrentID() int {
	return 1
}

// Получение всех заявок
func (r *Repository) GetAllCurrents() []Current {
	return []Current{
		r.GetCurrent(1),
		r.GetCurrent(2),
		r.GetCurrent(3),
	}
}

func (r *Repository) GetCurrentDeviceInfo(currentID, deviceID int) (CurrentDevice, error) {
	current := r.GetCurrent(currentID)
	for _, item := range current.DeviceItems {
		if item.DeviceID == deviceID {
			return item, nil
		}
	}
	return CurrentDevice{}, fmt.Errorf("товар не найден в заявке")
}

func (r *Repository) GetTotalItemsCount(id int) int {
	current := r.GetCurrent(id)
	total := 0
	for _, item := range current.DeviceItems {
		total += item.Quantity
	}
	return total
}
