package swagger

import (
	"encoding/json"
	"fmt"

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func NewSpecBuilder() *SpecBuilder {
	return &SpecBuilder{
		Info: &InfoObject{
			Version:        "1.0",
			Title:          "Swagger UI",
			Description:    "This is a sample server.",
			TermsOfService: "http://swagger.io/terms/",
			Contact: &ContactInfoObject{
				Name:  "API Support",
				Url:   "http://www.swagger.io/support",
				Email: "support@swagger.io",
			},
			License: &LicenseInfoObject{
				Name: "Apache 2.0",
				Url:  "http://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		Swagger: "2.0",
		Schemes: []string{"http", "https"},
	}
}

func (spec *SpecBuilder) SetTitle(title string) *SpecBuilder {
	spec.Info.Title = title
	return spec
}

func (spec *SpecBuilder) SetDescription(description string) *SpecBuilder {
	spec.Info.Description = description
	return spec
}

func (spec *SpecBuilder) SetVersion(version string) *SpecBuilder {
	spec.Info.Version = version
	return spec
}

func (spec *SpecBuilder) SetHost(host string) *SpecBuilder {
	spec.Host = host
	return spec
}

func (spec *SpecBuilder) SetBasePath(basePath string) *SpecBuilder {
	spec.BasePath = basePath
	return spec
}

func (spec *SpecBuilder) AddSecurity(security ...*SecuritySchemeObject) *SpecBuilder {
	if spec.SecurityDefinitions == nil {
		spec.SecurityDefinitions = make(map[string]*SecuritySchemeObject)
	}
	for _, v := range security {
		spec.SecurityDefinitions[v.Name] = v
	}
	return spec
}

// Build builds the swagger spec.
//
// It takes the SpecBuilder instance and returns the same instance
// with the swagger spec built.
func (spec *SpecBuilder) Build() *SpecBuilder {
	return spec
}

// SetUp sets up the swagger UI and API endpoint.
//
// It takes a prefix to mount the swagger UI and API endpoint, an app
// instance, and a SpecBuilder instance. It will parse the app's routes
// using the SpecBuilder and generate the swagger spec. It will then
// register a swagger handler with the app. The swagger handler will
// serve the swagger UI and API endpoint.
//
// The swagger UI will be available at the path <prefix>/doc.json. The
// swagger API endpoint will be available at the path <prefix>/doc.json.
//
// For example, if you call SetUp("/swagger", app, spec), you can access
// the swagger UI at http://localhost:8080/swagger/doc.json and the
// swagger API endpoint at http://localhost:8080/swagger/doc.json.
func SetUp(path string, app *core.App, spec *SpecBuilder) {
	spec.ParsePaths(app)
	mapper := RecursiveParseStandardSwagger(spec)
	jsonBytes, _ := json.Marshal(mapper)

	swaggerInfo := &swag.Spec{
		Version:          spec.Info.Version,
		Host:             spec.Host,
		BasePath:         spec.BasePath,
		Schemes:          spec.Schemes,
		Title:            spec.Info.Title,
		Description:      spec.Info.Description,
		InfoInstanceName: "swagger",
		SwaggerTemplate:  string(jsonBytes),
		LeftDelim:        "{{",
		RightDelim:       "}}",
	}

	route := fmt.Sprintf("%s%s", core.IfSlashPrefixString(app.Prefix), core.IfSlashPrefixString(path))

	swag.Register(swaggerInfo.InstanceName(), swaggerInfo)
	app.Mux.Handle("GET "+route+"/*", httpSwagger.Handler(
		httpSwagger.URL("http://"+spec.Host+route+"/doc.json"),
	))
}
