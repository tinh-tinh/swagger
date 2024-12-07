package swagger_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
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

	swagger.SetUp("/swagger", server, document)
	testServer := httptest.NewServer(server.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()
	resp, err := testClient.Get(testServer.URL + "/api/swagger/doc.json")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
