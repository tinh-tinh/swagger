package swagger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	document.ParsePaths(server)
	assert.Equal(t, "2.0", document.Swagger)
	assert.Equal(t, "1.0", document.Info.Version)
	assert.Equal(t, "Swagger UI", document.Info.Title)
	assert.Equal(t, "This is a sample server.", document.Info.Description)
	assert.Equal(t, "http://swagger.io/terms/", document.Info.TermsOfService)
	assert.Equal(t, "API Support", document.Info.Contact.Name)
	assert.Equal(t, "http://www.swagger.io/support", document.Info.Contact.Url)
	assert.Equal(t, "support@swagger.io", document.Info.Contact.Email)
	assert.Equal(t, "Apache 2.0", document.Info.License.Name)
	assert.Equal(t, "http://www.apache.org/licenses/LICENSE-2.0.html", document.Info.License.Url)
	assert.Equal(t, []string{"http", "https"}, document.Schemes)
}
