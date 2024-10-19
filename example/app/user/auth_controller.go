package user

import (
	"github.com/tinh-tinh/swagger"
	"github.com/tinh-tinh/swagger/example/app/user/dto"
	"github.com/tinh-tinh/tinhtinh/core"
)

func authController(module *core.DynamicModule) *core.DynamicController {
	authCtrl := module.NewController("Auth").Metadata(swagger.Tag("Auth")).Registry()

	authCtrl.Pipe(
		core.Body(&dto.SignUpUser{}),
	).Post("/", func(ctx core.Ctx) error {
		payload := ctx.Body().(*dto.SignUpUser)

		userService := authCtrl.Inject(USER_SERVICE).(CrudService)
		user := userService.Create(*payload)
		return ctx.JSON(core.Map{
			"status": user,
		})
	})

	return authCtrl
}
