package ds

import (
	"database/sql"
	"time"
	"github.com/google/uuid"
)

type Current struct {
	Current_ID   			uint      		`gorm:"primaryKey; not null"`
	Status       			string   		`gorm:"type:varchar(15); not null"`
	Created_At   			time.Time 		`gorm:"not null"`
	Creator_ID	 			uuid.UUID		`gorm:"type:integer(15); not null"`
	Moderator_ID  			uuid.NullUUID	`gorm:"type:integer(15); default: null"`
	Forming_Date 			sql.NullTime 	`gorm:"default:null"`
	Finish_Date  			sql.NullTime 	`gorm:"default:null"`
	//Amperage     float64      `gorm:"type:numeric(3,1)"`
	

	//Amount     int     `json:"amount"`
	VoltageBord     		float64 		`gorm:"type:decimal(10,2);default:11.5" json:"voltage_bord"`
	

	Creator   				Users 			`gorm:"foreignKey:Creator_ID"`
	Moderator 				Users 			`gorm:"foreignKey:Moderator_ID"`

	TotalAmperage  float64        `json:"total_amperage,omitempty"`
}
