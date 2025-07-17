package swagger_test

import (
	"time"

	"github.com/tinh-tinh/swagger/v2"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"github.com/tinh-tinh/tinhtinh/v2/middleware/storage"
)

type SignUpUser struct {
	Name     string    `validate:"isAlpha" example:"John"`
	Email    string    `validate:"required,isEmail" example:"john@gmail.com"`
	Password string    `validate:"required,isStrongPassword" example:"12345678@Tc"`
	Birth    time.Time `validate:"required" example:"2024-12-12"`
	Age      uint      `example:"3"`
	Score    float32   `example:"9.4"`
	IsAdmin  bool      `example:"false"`
}

type PostItem struct {
	Id   int    `example:"1"`
	Name string `example:"Title"`
}

type Post struct {
	Item     *PostItem `validate:"nested"`
	Managers []string  `example:"abc,def"`
}

type Param struct {
	ID string `path:"id"`
}

type FindUser struct {
	Name string `validate:"required,isAlpha" query:"name" example:"ac"`
	Age  int    `validate:"required,isInt" query:"age"`
}

func authController(module core.Module) core.Controller {
	authCtrl := module.NewController("Auth").Metadata(swagger.ApiTag("Auth")).Registry()

	authCtrl.Pipe(
		core.BodyParser[SignUpUser]{},
	).Post("", func(ctx core.Ctx) error {
		payload := ctx.Body().(*SignUpUser)

		return ctx.JSON(core.Map{
			"data": payload,
		})
	})

	return authCtrl
}

func managerController(module core.Module) core.Controller {
	ctrl := module.NewController("Users").Version("1").Metadata(
		swagger.ApiTag("User"),
		swagger.ApiSecurity("bearerAuth"),
	).Registry()

	ctrl.Pipe(
		core.QueryParser[FindUser]{},
	).Get("", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Pipe(
		core.BodyParser[SignUpUser]{},
	).Post("", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}

func UserModule(module core.Module) core.Module {
	userModule := module.New(core.NewModuleOptions{
		Controllers: []core.Controllers{managerController, authController},
	})

	return userModule
}

type UploadFile struct {
	File storage.File `example:"file"`
}

type Response struct {
	Title string `example:"Acrane"`
}

func postController(module core.Module) core.Controller {
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

	ctrl.Metadata(swagger.ApiOkResponse(&Response{})).Pipe(core.PathParser[Param]{}).Get("{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Pipe(core.PathParser[Param]{}).Patch("{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Pipe(core.PathParser[Param]{}, core.BodyParser[Post]{}).Put("{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Pipe(core.PathParser[Param]{}).Delete("{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}

func PostModule(module core.Module) core.Module {
	postModule := module.New(core.NewModuleOptions{
		Controllers: []core.Controllers{postController},
	})

	return postModule
}

func AppModule() core.Module {
	appModule := core.NewModule(core.NewModuleOptions{
		Imports: []core.Modules{
			UserModule,
			PostModule,
		},
	})

	return appModule
}
