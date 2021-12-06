package service

import (
	"admin/core/log"
	"admin/server/models"
	"admin/server/util"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"go.uber.org/zap"
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

func (u *User) UpdatePassword() error {
	data := make(map[string]interface{})

	user, err := models.GetUser(u.ID)
	if err != nil {
		return err
	}
	if user.ID == 0 {
		return fmt.Errorf("%s", "invalid user id")
	}

	passwordStr := util.GeneratePassword(16)
	fmt.Println(string(passwordStr))
	salt, password := util.GetSaltAndEncodedPassword(string(passwordStr))
	data["salt"] = salt
	data["password"] = password

	err = models.UpdateUser(u.ID, data)
	if err != nil {
		return err
	}

	return SendMail(user.Email, "重置密码", passwordStr)
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

	if user.Status == 0 {
		return nil, fmt.Errorf("invalid user status")
	}

	if user.LoginType == "standard" {
		if util.VerifyRawPassword(u.Password, user.Password, user.Salt) {
			return user, nil
		}
	} else if user.LoginType == "ldap" {
		ldapConfig, err := getLdapConfig()
		if err != nil {
			log.Logger.Error("ldap", zap.String("err", err.Error()))
		}
		if ldapConfig != nil {
			conn, err := ldap.DialURL("ldap://" + ldapConfig.Host + fmt.Sprint(":%d", ldapConfig.Port))
			if err != nil {
				log.Logger.Error("ldap", zap.String("err", err.Error()))
			}
			if conn != nil {
				err = conn.Bind(ldapConfig.DN, ldapConfig.Password)
				if err != nil {
					return nil, err
				}
				//searchRequest := ldap.NewSearchRequest()

				defer conn.Close()
			}
		}
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
