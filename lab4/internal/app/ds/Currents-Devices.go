package ds

type CurrentDevices struct {
	CurrDev_ID uint    `gorm:"primaryKey" json:"curr_dev_id"`
	Current_ID uint    `gorm:"not null" json:"current_id"`
	Device_ID  uint    `gorm:"not null" json:"device_id"`
	Amount     int     `json:"amount"`

	//VoltageBord     float64 `gorm:"type:decimal(10,2);default:11.5" json:"voltage_bord"`
	
	Amperage   float64 `gorm:"type:numeric(10,3)" json:"amperage"`

	Device  Device  `gorm:"foreignKey:Device_ID;references:Device_ID"`
	Current Current `gorm:"foreignKey:Current_ID;references:Current_ID"`
}
