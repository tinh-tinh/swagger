package swagger

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
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
		Components: &ComponentObject{
			Schemas:         make(map[string]*SchemaObject),
			SecuritySchemes: make(map[string]*SecuritySchemeObject),
		},
		Openapi: "3.0.0",
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

func (spec *SpecBuilder) SetServer(server *ServerObject) *SpecBuilder {
	spec.Servers = append(spec.Servers, server)
	return spec
}

func (spec *SpecBuilder) AddSecurity(security ...*SecuritySchemeObject) *SpecBuilder {
	for _, v := range security {
		spec.Components.SecuritySchemes[v.Name] = v
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
func SetUp(path string, app *core.App, spec *SpecBuilder, configs ...Config) {
	spec.ParsePaths(app)
	jsonBytes, _ := json.Marshal(spec)

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromData(jsonBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Validate document
	_ = doc.Validate(ctx)

	// Serve the OpenAPI document as JSON
	app.Mux.Handle("/openapi.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(doc); err != nil {
			http.Error(w, "Failed to encode OpenAPI document", http.StatusInternalServerError)
		}
	}))

	route := fmt.Sprintf("%s%s", core.IfSlashPrefixString(app.Prefix), core.IfSlashPrefixString(path))
	// Serve Swagger UI HTML from CDN
	app.Mux.Handle(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var persistAuth string
		if len(configs) > 0 {
			config := configs[0]
			if config.PersistAuthorization {
				persistAuth += "persistAuthorization: true,\n"
			}
		}
		htmlParser := fmt.Sprintf(`
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>Swagger UI</title>
            <!-- Load Swagger UI from CDN -->
            <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4.15.5/swagger-ui.css" />
            <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
            <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@4.15.5/swagger-ui-standalone-preset.js"></script>
        </head>
        <body>
            <div id="swagger-ui"></div>
            <script>
                const ui = SwaggerUIBundle({
                    url: "/openapi.json",  // URL for your OpenAPI spec
                    dom_id: '#swagger-ui',
                    deepLinking: true,
                    presets: [
                        SwaggerUIBundle.presets.apis,
                        SwaggerUIBundle.SwaggerUIStandalonePreset
                    ],
                    layout: "BaseLayout",
					%s
                });
            </script>
        </body>
        </html>
        `, persistAuth)
		w.Header().Set("Content-Type", "text/html")
		if _, err := w.Write([]byte(htmlParser)); err != nil {
			fmt.Println(err)
		}
	}))
}
