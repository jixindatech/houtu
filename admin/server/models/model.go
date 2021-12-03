package models

import (
	"admin/config"
	"admin/server/service"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var db *gorm.DB

type Model struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updateAt"`
}

const (
	admin = "admin"
	role  = "admin"
)

func Setup(cfg *config.DataBase, password string) error {
	var err error
	if cfg.Type == "mysql" {
		db, err = gorm.Open(cfg.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Name))
		if err != nil {
			return err
		}
	}

	if len(cfg.TablePrefix) > 0 {
		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return cfg.TablePrefix + defaultTableName
		}
	}

	db.SingularTable(true)
	db.LogMode(false)

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.AutoMigrate(
		User{},
	)

	// save admin record or update admin password
	if len(password) > 0 {
		user := service.User{}
		err := user.SaveAdmin(admin, role, password)
		if err != nil {
			return err
		}
	}
	return nil
}

func CloseDB() {
	_ = db.Close()
}
