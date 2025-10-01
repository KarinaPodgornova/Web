package ds

type Device struct {
    Device_ID       uint    `gorm:"primaryKey" json:"device_id"`
    Name            string  `gorm:"type:varchar(100);not null" json:"name"`
    Type            string  `gorm:"type:varchar(30)" json:"type"`
    PowerNominal    float64 `gorm:"type:decimal(10,2);not null" json:"power_nominal"`
    Resistance      float64 `gorm:"type:decimal(10,2)" json:"resistance"`
    VoltageNominal  float64 `gorm:"type:decimal(10,2);default:12.0" json:"voltage_nominal"`
    
    CoeffReserve    float64 `gorm:"type:decimal(10,2);default:2.0" json:"coeff_reserve"`
    CoeffEfficiency float64 `gorm:"type:decimal(10,2);default:0.85" json:"coeff_efficiency"`
    CurrentRequired float64 `gorm:"type:decimal(10,2)" json:"current_required"`
    Description     string  `gorm:"type:text" json:"description"`
    Image           string  `gorm:"type:varchar(100)" json:"image"`
    InStock         bool    `gorm:"type:boolean;default:true" json:"in_stock"`
    IsDelete        bool    `gorm:"type:boolean;default:false" json:"is_delete"`
}