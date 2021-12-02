package service

import (
	"admin/server/models"
	"admin/server/util"
)

type User struct {
	ID uint

	Username    string
	DisplayName string
	LoginType   string
	Password    string
	Salt        string
	Email       string
	Phone       string
	Status      int
	Role        string

	Remark string
}

func (u *User) Save() error {
	data := map[string]interface{}{
		"username":    u.Username,
		"displayName": u.DisplayName,
		"loginType":   u.LoginType,
		"email":       u.Email,
		"phone":       u.Phone,
		"status":      u.Status,
		"role":        u.Role,
		"remark":      u.Remark,
	}
	if u.ID > 0 {
		if len(u.Password) > 0 {
			salt, password := util.GetSaltAndEncodedPassword(u.Password)
			data["salt"] = salt
			data["password"] = password
		}
		return models.UpdateUser(u.ID, data)
	}

	return models.AddUser(data)
}
