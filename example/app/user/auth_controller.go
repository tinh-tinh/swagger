package user

import (
	"github.com/tinh-tinh/swagger/example/app/user/dto"
	"github.com/tinh-tinh/tinhtinh/core"
)

func authController(module *core.DynamicModule) *core.DynamicController {
	authCtrl := module.NewController("Auth").Tag("Global")

	authCtrl.Pipe(
		core.Body(&dto.SignUpUser{}),
	).Post("/", func(ctx core.Ctx) {
		payload := ctx.Body().(*dto.SignUpUser)

		userService := authCtrl.Inject(USER_SERVICE).(CrudService)
		user := userService.Create(*payload)
		ctx.JSON(core.Map{
			"status": user,
		})
	})

	return authCtrl
}
