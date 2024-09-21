package main

import (
	"github.com/tinh-tinh/swagger"
	"github.com/tinh-tinh/swagger/example/app"
	"github.com/tinh-tinh/tinhtinh/core"
)

func main() {
	server := core.CreateFactory(app.NewModule, "api")

	document := swagger.NewSpecBuilder()
	document.SetHost("localhost:3000").SetBasePath("/api").AddSecurity(&swagger.SecuritySchemeObject{
		Type: "apiKey",
		In:   "header",
		Name: "Authorization",
	})

	swagger.SetUp("docs", server, document)
	server.Listen(3000)
}
