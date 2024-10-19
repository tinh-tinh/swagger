package post

import (
	"github.com/tinh-tinh/swagger"
	"github.com/tinh-tinh/tinhtinh/core"
)

func controller(module *core.DynamicModule) *core.DynamicController {
	ctrl := module.NewController("Posts").Metadata(swagger.Tag("Post")).Registry()

	ctrl.Post("/", func(ctx core.Ctx) error {
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
