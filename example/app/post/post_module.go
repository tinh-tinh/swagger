package post

import (
	"github.com/tinh-tinh/tinhtinh/core"
)

func Module(module *core.DynamicModule) *core.DynamicModule {
	postModule := module.New(core.NewModuleOptions{
		Controllers: []core.Controller{controller},
		Providers:   []core.Provider{service},
	})

	return postModule
}
