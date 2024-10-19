package user

import (
	"github.com/tinh-tinh/swagger"
	"github.com/tinh-tinh/swagger/example/app/user/dto"
	"github.com/tinh-tinh/tinhtinh/core"
)

func managerController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("Users").Version("1").Metadata(
		swagger.ApiTag("User"),
		swagger.ApiSecurity("authorization"),
	).Registry()

	ctrl.Pipe(
		core.Query(&dto.FindUser{}),
	).Get("/", func(ctx core.Ctx) error {
		userService := ctrl.Inject(USER_SERVICE).(CrudService)
		data := userService.GetAll()
		return ctx.JSON(core.Map{"data": data})
	})

	ctrl.Pipe(
		core.Body(&dto.SignUpUser{}),
	).Post("/", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}
