package service

import (
	"admin/server/models"
	"admin/server/util"
	"fmt"
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

	Page     int
	PageSize int
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

func (u *User) Delete() error {
	if u.ID == 1 {
		return fmt.Errorf("%s", "invalid user id")
	}
	return models.DeleteUser(u.ID)
}

func (u *User) Get() (*models.User, error) {
	return models.GetUser(u.ID)
}

func (u *User) GetList() ([]*models.User, uint, error) {
	var query = make(map[string]interface{})
	if len(u.Username) > 0 {
		query["username"] = u.Username
	}

	return models.GetUsers(query, u.Page, u.PageSize)
}

func (u *User) GetLoginUser() (*models.User, error) {
	user, err := models.GetUserByUsername(u.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("invalid username")
	}

	if util.VerifyRawPassword(u.Password, user.Password, user.Salt) {
		return user, nil
	}

	return nil, fmt.Errorf("invalid password")
}

func SaveAdmin(id uint, username, role, password string) error {
	admin, err := models.GetUser(id)
	if err != nil {
		return err
	}

	if admin.ID > 0 {
		salt, encodedPassword := util.GetSaltAndEncodedPassword(password)
		data := make(map[string]interface{})
		data["salt"] = salt
		data["password"] = encodedPassword

		return models.UpdateUser(admin.ID, data)
	}

	data := map[string]interface{}{
		"username":    username,
		"displayName": "admin",
		"loginType":   "standard",
		"email":       "admin@admin.com",
		"phone":       "13200000000",
		"status":      1,
		"role":        role,
		"remark":      "administrator",
	}

	salt, encodedPassword := util.GetSaltAndEncodedPassword(password)
	data["salt"] = salt
	data["password"] = encodedPassword

	return models.AddUser(data)
}
