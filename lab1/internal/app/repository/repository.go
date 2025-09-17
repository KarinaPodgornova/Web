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

type Order struct {
    ID    int
    Title string
	Description string
	Specs       string
	Power       string
}

type CartItem struct {
    ProductID   int
    Title       string
    Power       string
    Resistance  string
}

type Cart struct {
    Items []CartItem
    Count int
}

func (r *Repository) GetOrders() ([]Order, error) {
    // Реальные товары BMW - ВСЕ 10 штук!
    orders := []Order{
		{
            ID: 1, 
            Title: "Универсальный адаптер BMW Step-In", 
            Power: "10 Вт",
            Description: "Вставной адаптер Snap-In предназначен для всех подходящих телефонов с разъемом Micro USB. Помимо зарядки телефона обеспечивает также улучшенный прием и надежное крепление. Благодаря разъему USB записанная на телефоне музыка может удобно воспроизводиться через динамики аудиосистемы BMW.",
            Specs: "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 10 Вт\nРазмеры: 145 x 72 x 22 мм",
        },
        {
            ID: 2, 
            Title: "Адаптер BMW Step-In для мобильных телефонов Apple", 
            Power: "15 Вт",
            Description: "Специально разработан для устройств Apple. Обеспечивает быструю зарядку и стабильное соединение с медиасистемой BMW.",
            Specs: "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 15 Вт\nРазмеры: 142 x 70 x 20 мм",
        },
        {
            ID: 3, 
            Title: "Автомобильный отопитель BMW", 
            Power: "150 Вт",
            Description: "Дополнительный отопитель салона с дистанционным управлением. Обеспечивает комфортную температуру в холодное время года.",
            Specs: "Входное напряжение: 12 В\nПотребляемая мощность: 150 Вт\nТепловая мощность: 2 кВт\nРазмеры: 180 x 120 x 80 мм",
        },
        {
            ID: 4, 
            Title: "Автомобильный отопитель BMW Premium", 
            Power: "200 Вт",
            Description: "Усовершенствованная модель отопителя с цифровым управлением и таймером. Поддерживает заданную температуру автоматически.",
            Specs: "Входное напряжение: 12 В\nПотребляемая мощность: 200 Вт\nТепловая мощность: 2.5 кВт\nРазмеры: 190 x 125 x 85 мм",
        },
        {
            ID: 5, 
            Title: "Переходник Micro USB", 
            Power: "5 Вт",
            Description: "Компактный переходник для подключения устройств с Micro USB к медиасистеме BMW. Обеспечивает зарядку и передачу данных.",
            Specs: "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 5 Вт\nРазмеры: 45 x 20 x 15 мм",
        },
        {
            ID: 6, 
            Title: "Адаптер Media для Apple iPod / iPhone", 
            Power: "12 Вт",
            Description: "Адаптер для полной интеграции устройств Apple с медиасистемой BMW. Поддержка управления через iDrive.",
            Specs: "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 12 Вт\nРазмеры: 95 x 55 x 25 мм",
        },
        {
            ID: 7, 
            Title: "Музыкальный/медийный адаптер BMW для Apple iPod/iPhone", 
            Power: "18 Вт",
            Description: "Премиальный адаптер с поддержкой высококачественного звука и видео. Полная интеграция с интерфейсом iDrive.",
            Specs: "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 18 Вт\nРазмеры: 100 x 60 x 28 мм",
        },
        {
            ID: 8, 
            Title: "Кабель-адаптер BMW для Micro-USB", 
            Power: "10 Вт",
            Description: "Оригинальный кабель-адаптер для устройств Android. Обеспечивает быструю зарядку и стабильное соединение.",
            Specs: "Входное напряжение: 12 В\nВыходное напряжение: 5 В\nВыходная мощность: 10 Вт\nДлина кабеля: 1.2 м",
        },
        {
            ID: 9, 
            Title: "CD — привод дополнительный", 
            Power: "25 Вт",
            Description: "Встраиваемый дополнительный CD-привод для BMW. Поддержка аудио CD и MP3-дисков. Интеграция со штатной аудиосистемой.",
            Specs: "Входное напряжение: 12 В\nПотребляемая мощность: 25 Вт\nФорматы: CD, CD-R, CD-RW, MP3\nРазмеры: 180 x 145 x 45 мм",
        },
        {
            ID: 10, 
            Title: "Аксессуары для ключа-браслета", 
            Power: "2 Вт",
            Description: "Стильный браслет-чехол для ключа BMW. Защита от повреждений и удобное ношение. Совместим со всеми моделями ключей BMW.",
            Specs: "Материал: силикон\nЦвет: черный\nСовместимость: все ключи BMW\nВес: 15 г",
        },
	
    }

    if len(orders) == 0 {
        return nil, fmt.Errorf("массив пустой")
    }

    return orders, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
    orders, err := r.GetOrders()
    if err != nil {
        return Order{}, err
    }

    for _, order := range orders {
        if order.ID == id {
            return order, nil
        }
    }
    return Order{}, fmt.Errorf("товар не найден")
}


func (r *Repository) GetOrdersByTitle(searchQuery string) ([]Order, error) {
    orders, err := r.GetOrders()
    if err != nil {
        return []Order{}, err
    }

    var result []Order
    for _, order := range orders {
        if strings.Contains(strings.ToLower(order.Title), strings.ToLower(searchQuery)) {
            result = append(result, order)
        }
    }

    return result, nil
}
