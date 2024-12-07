package swagger_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/swagger"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Spec(t *testing.T) {
	server := core.CreateFactory(AppModule)
	server.SetGlobalPrefix("api")

	document := swagger.NewSpecBuilder().
		SetTitle("Swagger Document UI").
		SetDescription("This is a document for apis").
		SetVersion("1.0.0").
		SetHost("localhost:3000").
		SetBasePath("/api").
		AddSecurity(&swagger.SecuritySchemeObject{
			Type: "apiKey",
			In:   "header",
			Name: "Authorization",
		}).Build()

	document.ParsePaths(server)
	// Default
	assert.Equal(t, "2.0", document.Swagger)
	assert.Equal(t, "1.0.0", document.Info.Version)
	assert.Equal(t, "Swagger Document UI", document.Info.Title)
	assert.Equal(t, "This is a document for apis", document.Info.Description)
	assert.Equal(t, "http://swagger.io/terms/", document.Info.TermsOfService)
	assert.Equal(t, "API Support", document.Info.Contact.Name)
	assert.Equal(t, "http://www.swagger.io/support", document.Info.Contact.Url)
	assert.Equal(t, "support@swagger.io", document.Info.Contact.Email)
	assert.Equal(t, "Apache 2.0", document.Info.License.Name)
	assert.Equal(t, "http://www.apache.org/licenses/LICENSE-2.0.html", document.Info.License.Url)
	assert.Equal(t, []string{"http", "https"}, document.Schemes)

	assert.Equal(t, "localhost:3000", document.Host)
	assert.Equal(t, "/api", document.BasePath)
	assert.Equal(t, 1, len(document.SecurityDefinitions))
	assert.NotNil(t, document.SecurityDefinitions["Authorization"])
	assert.Equal(t, "apiKey", document.SecurityDefinitions["Authorization"].Type)
	assert.Equal(t, "header", document.SecurityDefinitions["Authorization"].In)

	paths := document.Paths
	assert.NotNil(t, paths)

	fmt.Println(paths)
	assert.Len(t, paths, 4)
	assert.Equal(t, []string{"Auth"}, paths["/auth"].Post.Tags)
	assert.Empty(t, paths["/auth"].Post.Summary)
	assert.Empty(t, paths["/auth"].Post.Description)
	assert.Empty(t, paths["/auth"].Post.OperationID)
	assert.Empty(t, paths["/auth"].Post.Consumes)
	assert.Empty(t, paths["/auth"].Post.Produces)
	assert.Equal(t, "SignUpUser", paths["/auth"].Post.Parameters[0].Name)
	assert.Equal(t, "body", paths["/auth"].Post.Parameters[0].In)
	assert.False(t, paths["/auth"].Post.Parameters[0].Required)
	assert.Empty(t, paths["/auth"].Post.Parameters[0].Items)
	assert.Equal(t, "#/definitions/signUpUser", paths["/auth"].Post.Parameters[0].Schema.Ref)
	assert.False(t, paths["/auth"].Post.Parameters[0].Schema.Required)
	assert.Empty(t, paths["/auth"].Post.Parameters[0].Schema.Enum)
	assert.Empty(t, paths["/auth"].Post.Schemes)
	assert.False(t, paths["/auth"].Post.Deprecated)
	assert.Empty(t, paths["/auth"].Post.Security)
	assert.Equal(t, "Ok", paths["/auth"].Post.Responses["200"].Description)

	assert.Equal(t, []string{"User"}, paths["/users"].Get.Tags)
	assert.Equal(t, []string{"User"}, paths["/users"].Post.Tags)
	assert.NotNil(t, document.Paths["/users"].Post.Security[0])
	assert.Empty(t, document.Paths["/users"].Post.Security[0]["authorization"])

	assert.Equal(t, []string{"User"}, paths["/users"].Get.Tags)
	assert.Equal(t, []string{"User"}, paths["/users"].Post.Tags)
	assert.NotNil(t, document.Paths["/users"].Post.Security[0])
	assert.Empty(t, document.Paths["/users"].Post.Security[0]["authorization"])
	assert.Equal(t, "query", paths["/users"].Get.Parameters[0].In)
	assert.Equal(t, "name", paths["/users"].Get.Parameters[0].Name)
	assert.Equal(t, "string", paths["/users"].Get.Parameters[0].Type)
	assert.Equal(t, "ac", paths["/users"].Get.Parameters[0].Default)
	assert.Equal(t, "age", paths["/users"].Get.Parameters[1].Name)
	assert.Equal(t, "integer", paths["/users"].Get.Parameters[1].Type)

	assert.Equal(t, []string{"Post"}, paths["/posts"].Post.Tags)
	assert.Equal(t, []string{"multipart/form-data"}, paths["/posts"].Post.Consumes)
	assert.Equal(t, "file", paths["/posts"].Post.Parameters[0].Name)
	assert.Equal(t, "formData", paths["/posts"].Post.Parameters[0].In)
	assert.Equal(t, "file upload", paths["/posts"].Post.Parameters[0].Description)
	assert.True(t, paths["/posts"].Post.Parameters[0].Required)
	assert.Empty(t, paths["/posts"].Post.Parameters[0].Items)
}
