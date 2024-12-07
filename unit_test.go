package swagger_test

import (
	"time"

	"github.com/tinh-tinh/swagger"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware/storage"
)

type SignUpUser struct {
	Name     string    `validate:"isAlpha" example:"John"`
	Email    string    `validate:"required,isEmail" example:"john@gmail.com"`
	Password string    `validate:"required,isStrongPassword" example:"12345678@Tc"`
	Birth    time.Time `validate:"required" example:"2024-12-12"`
}

type FindUser struct {
	Name string `validate:"required,isAlpha" query:"name" example:"ac"`
	Age  uint   `validate:"required,isInt" query:"age"`
}

func authController(module *core.DynamicModule) *core.DynamicController {
	authCtrl := module.NewController("Auth").Metadata(swagger.ApiTag("Auth")).Registry()

	authCtrl.Pipe(
		core.Body(SignUpUser{}),
	).Post("", func(ctx core.Ctx) error {
		payload := ctx.Body().(*SignUpUser)

		return ctx.JSON(core.Map{
			"data": payload,
		})
	})

	return authCtrl
}

func managerController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("Users").Version("1").Metadata(
		swagger.ApiTag("User"),
		swagger.ApiSecurity("authorization"),
	).Registry()

	ctrl.Pipe(
		core.Query(FindUser{}),
	).Get("", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Pipe(
		core.Body(SignUpUser{}),
	).Post("", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}

func UserModule(module *core.DynamicModule) *core.DynamicModule {
	userModule := module.New(core.NewModuleOptions{
		Controllers: []core.Controller{managerController, authController},
	})

	return userModule
}

type UploadFile struct {
	File storage.File `example:"file"`
}

func postController(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("Posts").Metadata(swagger.ApiTag("Post")).Registry()

	ctrl.Metadata(
		swagger.ApiConsumer("multipart/form-data"),
		swagger.ApiFile(swagger.FileOptions{
			Name:        "file",
			Description: "file upload",
			Required:    true,
		}),
	).Post("", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Get("", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Get("{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Put("{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Delete("{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}

func PostModule(module *core.DynamicModule) *core.DynamicModule {
	postModule := module.New(core.NewModuleOptions{
		Controllers: []core.Controller{postController},
	})

	return postModule
}

func AppModule() *core.DynamicModule {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Module{
			UserModule,
			PostModule,
		},
	})

	return appModule
}
