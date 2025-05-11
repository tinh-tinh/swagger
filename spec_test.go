package swagger_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/swagger/v2"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_Spec(t *testing.T) {
	server := core.CreateFactory(AppModule)
	server.SetGlobalPrefix("api")

	document := swagger.NewSpecBuilder().
		SetTitle("Swagger Document UI").
		SetServer(&swagger.ServerObject{
			Url: "http://localhost:3000/api",
		}).
		SetDescription("This is a document for apis").
		SetVersion("1.0.0").
		AddSecurity(&swagger.SecuritySchemeObject{
			Type:         "http",
			Scheme:       "Bearer",
			BearerFormat: "JWT",
			Name:         "bearerAuth",
		}).Build()

	document.ParsePaths(server)
	// Default
	assert.Equal(t, "3.0.0", document.Openapi)
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
	assert.Equal(t, []*swagger.ServerObject{{Url: "http://localhost:3000/api"}}, document.Servers)

	assert.NotNil(t, document.Components.Schemas)
	assert.NotNil(t, document.Components.Schemas["SignUpUser"])
	assert.Equal(t, "object", document.Components.Schemas["SignUpUser"].Type)
	assert.NotNil(t, document.Components.Schemas["SignUpUser"].Properties)
	assert.Equal(t, "string", document.Components.Schemas["SignUpUser"].Properties["Email"].Type)
	assert.Equal(t, "john@gmail.com", document.Components.Schemas["SignUpUser"].Properties["Email"].Example)

	assert.NotNil(t, document.Components.SecuritySchemes["bearerAuth"])
	assert.Equal(t, "Bearer", document.Components.SecuritySchemes["bearerAuth"].Scheme)
	assert.Equal(t, "bearerAuth", document.Components.SecuritySchemes["bearerAuth"].Name)
	assert.Equal(t, "JWT", document.Components.SecuritySchemes["bearerAuth"].BearerFormat)

	paths := document.Paths
	assert.NotNil(t, paths)
	assert.Len(t, paths, 4)
	assert.NotNil(t, paths["/api/auth"])
	assert.NotNil(t, paths["/api/auth"].Post)
	assert.Equal(t, []string{"Auth"}, paths["/api/auth"].Post.Tags)
	assert.Empty(t, paths["/api/auth"].Post.Summary)
	assert.Empty(t, paths["/api/auth"].Post.Description)
	assert.Empty(t, paths["/api/auth"].Post.OperationID)
	assert.Empty(t, paths["/api/auth"].Post.Consumes)
	assert.Empty(t, paths["/api/auth"].Post.Produces)
	assert.NotNil(t, paths["/api/auth"].Post.RequestBody.Content["SignUpUser"])
	assert.Equal(t, "#/components/schemas/signUpUser", paths["/api/auth"].Post.RequestBody.Content["SignUpUser"].Schema.Ref)
	assert.Empty(t, paths["/api/auth"].Post.Schemes)
	assert.False(t, paths["/api/auth"].Post.Deprecated)
	assert.Empty(t, paths["/api/auth"].Post.Security)
	assert.Equal(t, "Ok", paths["/api/auth"].Post.Responses["200"].Description)

	assert.NotNil(t, paths["/api/users"])
	assert.NotNil(t, paths["/api/users"].Get)
	assert.Equal(t, []string{"User"}, paths["/api/users"].Get.Tags)
	assert.NotNil(t, paths["/api/users"].Post)
	assert.Equal(t, []string{"User"}, paths["/api/users"].Post.Tags)
	assert.NotNil(t, document.Paths["/api/users"].Post.Security[0])
	assert.Empty(t, document.Paths["/api/users"].Post.Security[0]["bearerAuth"])

	assert.Equal(t, "query", paths["/api/users"].Get.Parameters[0].In)
	assert.Equal(t, "name", paths["/api/users"].Get.Parameters[0].Name)
	assert.Equal(t, "string", paths["/api/users"].Get.Parameters[0].Schema.Type)
	assert.Equal(t, "ac", paths["/api/users"].Get.Parameters[0].Default)
	assert.Equal(t, "age", paths["/api/users"].Get.Parameters[1].Name)
	assert.Equal(t, "integer", paths["/api/users"].Get.Parameters[1].Schema.Type)

	assert.Equal(t, []string{"Post"}, paths["/api/posts"].Post.Tags)
	assert.Equal(t, []string{"multipart/form-data"}, paths["/api/posts"].Post.Consumes)
	assert.Equal(t, "file", paths["/api/posts"].Post.Parameters[0].Name)
	assert.Equal(t, "formData", paths["/api/posts"].Post.Parameters[0].In)
	assert.Equal(t, "file upload", paths["/api/posts"].Post.Parameters[0].Description)
	assert.True(t, paths["/api/posts"].Post.Parameters[0].Required)

	assert.Equal(t, "id", paths["/api/posts/{id}"].Get.Parameters[0].Name)
	assert.Equal(t, "path", paths["/api/posts/{id}"].Get.Parameters[0].In)
	assert.Equal(t, "#/components/schemas/response", paths["/api/posts/{id}"].Get.Responses["200"].Schema.Ref)

}
