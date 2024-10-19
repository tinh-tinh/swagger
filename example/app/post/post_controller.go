package post

import (
	"github.com/tinh-tinh/swagger"
	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/middleware/storage"
)

type UploadFile struct {
	File storage.File `example:"file"`
}

func controller(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("Posts").Metadata(swagger.ApiTag("Post")).Registry()

	ctrl.Metadata(
		swagger.ApiConsumer("multipart/form-data"),
		swagger.ApiFile(swagger.FileOptions{
			Name:        "file",
			Description: "file upload",
			Required:    true,
		}),
	).Post("/", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Get("/", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Get("/{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Put("/{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	ctrl.Delete("/{id}", func(ctx core.Ctx) error {
		return ctx.JSON(core.Map{"data": "ok"})
	})

	return ctrl
}
