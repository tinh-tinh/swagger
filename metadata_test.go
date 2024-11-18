package swagger_test

import (
	"testing"

	"github.com/tinh-tinh/swagger"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Tag(t *testing.T) {
	server := core.CreateFactory(AppModule)
	server.SetGlobalPrefix("api")

	document := swagger.NewSpecBuilder()
	document.SetHost("localhost:3000").SetBasePath("/api").AddSecurity(&swagger.SecuritySchemeObject{
		Type: "apiKey",
		In:   "header",
		Name: "Authorization",
	})

	swagger.SetUp("docs", server, document)
}
