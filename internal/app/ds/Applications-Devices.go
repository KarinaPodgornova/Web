
package ds

type ApplicationDevices struct {
    AppDev_ID       uint    `gorm:"primaryKey" json:"app_dev_id"`
    Application_ID  uint    `gorm:"not null" json:"application_id"`
    Device_ID       uint    `gorm:"not null" json:"device_id"`
    Amount          int     `json:"amount"`
    Notes           string  `gorm:"type:varchar(255)" json:"notes"`
    Amperage        float64 `gorm:"type:numeric(10,3)" json:"amperage"`
}