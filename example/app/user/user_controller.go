package user

import (
	"github.com/tinh-tinh/swagger/example/app/user/dto"
	"github.com/tinh-tinh/tinhtinh/core"
)

func managerController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("Users").Version("1").Guard()

	ctrl.AddSecurity("authorization").Pipe(
		core.Query(&dto.FindUser{}),
	).Get("/", func(ctx core.Ctx) {
		userService := ctrl.Inject(USER_SERVICE).(CrudService)
		data := userService.GetAll()
		ctx.JSON(core.Map{"data": data})
	})

	ctrl.Pipe(
		core.Body(&dto.SignUpUser{}),
	).Post("/", func(ctx core.Ctx) {
		ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}
