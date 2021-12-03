package app

import (
	"admin/core/rbac"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var validate *validator.Validate

var phoneReg *regexp.Regexp

func ValidatePhone(fl validator.FieldLevel) bool {
	var err error
	if phoneReg == nil {
		phoneReg, err = regexp.Compile("^1[3456789]\\d{9}$")
		if err != nil {
			return false
		}
	}

	phone := fl.Field().String()
	return phoneReg.Match([]byte(phone))
}

func ValidateRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	if rbac.ROLES[role] != false {
		return true
	}

	return false
}

func SetupValidate() error {
	var err error

	validate = validator.New()
	err = validate.RegisterValidation("phone", ValidatePhone)
	if err != nil {
		return err
	}

	err = validate.RegisterValidation("role", ValidateRole)
	if err != nil {
		return err
	}

	return nil
}

func BindAndValid(c *gin.Context, form interface{}) error {
	err := c.Bind(form)
	if err != nil {
		return err
	}

	return validate.Struct(form)
}