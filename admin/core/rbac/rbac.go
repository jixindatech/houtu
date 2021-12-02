package rbac

import (
	"admin/config"
	"github.com/casbin/casbin/v2"
)

var enforcer *casbin.Enforcer

func Setup(cfg *config.Rbac) error {
	var err error
	enforcer, err = casbin.NewEnforcer(cfg.Model, cfg.Policy)
	if err != nil {
		return err
	}

	enforcer.EnableLog(true)

	return nil
}

func GetEnforcer() *casbin.Enforcer {
	return enforcer
}
