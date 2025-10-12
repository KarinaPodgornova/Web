package serializer

import "lab3/internal/app/ds"

// UserJSON представляет пользователя в формате JSON
type UserJSON struct {
	ID          uint   `json:"id"`          // Идентификатор пользователя
	Login       string `json:"login"`       // Логин пользователя
	Password    string `json:"password"`    // Пароль пользователя (не экспортируется в JSON)
	IsModerator bool   `json:"is_moderator"`// Флаг модератора
}

// UserToJSON преобразует ds.Users в UserJSON
func UserToJSON(user ds.Users) UserJSON {
	return UserJSON{
		ID:          user.User_ID,
		Login:       user.Login,
		Password:    user.Password,
		IsModerator: user.IsModerator,
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