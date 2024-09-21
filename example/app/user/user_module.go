package user

import (
	"github.com/tinh-tinh/tinhtinh/core"
)

func Module(module *core.DynamicModule) *core.DynamicModule {
	userModule := module.New(core.NewModuleOptions{
		Controllers: []core.Controller{managerController, authController},
		Providers:   []core.Provider{service},
	})

	return userModule
}
