package rbac

import (
	"github.com/casbin/casbin/v2"
)

var enforcer *casbin.Enforcer

func setupCasbin(model, policy string) error {
	var err error
	enforcer, err = casbin.NewEnforcer(model, policy)
	if err != nil {
		return err
	}

	enforcer.EnableLog(true)

	return nil
}
