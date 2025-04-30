package swagger_test

import (
	"net/http/httptest"
	"testing"

	"github.com/tinh-tinh/swagger/v2"
	"github.com/tinh-tinh/tinhtinh/v2/core"
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

	swagger.SetUp("/swagger", server, document)
	testServer := httptest.NewServer(server.PrepareBeforeListen())
	defer testServer.Close()
}
