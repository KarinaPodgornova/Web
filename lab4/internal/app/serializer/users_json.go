package serializer

import (
	"lab4/internal/app/ds"

	"github.com/google/uuid"
)

// UserJSON представляет пользователя в формате JSON
type UserJSON struct {
	ID 			uuid.UUID	`json:"id"`   
	Login       string 		`json:"login"`      
	Password    string 		`json:"password"`  
	IsModerator bool   		`json:"is_moderator"`
}

// UserToJSON преобразует ds.Users в UserJSON
func UserToJSON(user ds.Users) UserJSON {
	return UserJSON{
		ID: 			uuid.UUID(user.User_ID),
		Login:       	user.Login,
		Password:    	user.Password,
		IsModerator: 	user.IsModerator,
	}
}

// UserFromJSON преобразует UserJSON в ds.Users
func UserFromJSON(userJSON UserJSON) ds.Users {
	return ds.Users{
		User_ID:     userJSON.ID,
		Login:       userJSON.Login,
		Password:    userJSON.Password,
		IsModerator: userJSON.IsModerator,
	}
}