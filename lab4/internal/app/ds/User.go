package ds

import "github.com/google/uuid"


type Users struct {
	User_ID     uuid.UUID   `gorm:"primary_key;autoIncrement"`
	Login       string 		`gorm:"type:varchar(20);unique;not null" json:"login"`
	Password    string 		`gorm:"type:varchar(20);not null" json:"-"`
	IsModerator bool   		`gorm:"type:boolean;default:false" json:"is_moderator"`
}