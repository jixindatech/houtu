package service

import "admin/server/models"

type Email struct {
	ID uint

	Host     string
	Port     int
	Sender   string
	Password string
}

var emailCache *Email

func (e *Email) Save() error {
	data := make(map[string]interface{})
	data["host"] = e.Host
	data["port"] = e.Port
	data["sender"] = e.Sender
	data["password"] = e.Password

	if e.ID > 0 {
		emailCache = nil
		return models.UpdateEmail(e.ID, data)
	}

	return models.AddEmail(data)
}

func (e *Email) Get() (*models.Email, error) {
	return models.GetEmail()
}

func SendMail() error {
	if emailCache == nil {
		email, err := models.GetEmail()
		if err != nil {
			return err
		}
		emailCache = new(Email)
		emailCache.Host = email.Host
		emailCache.Port = email.Port
		emailCache.Sender = email.Sender
		emailCache.Password = email.Password
	}

	return nil
}
