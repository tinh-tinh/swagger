package swagger

// -------- Info Object --------
type InfoObject struct {
	Title          string             `json:"title"` // required
	Description    string             `json:"description,omitempty"`
	Version        string             `json:"version"`                  // required
	TermsOfService string             `json:"termsOfService,omitempty"` // typo fixed
	Contact        *ContactInfoObject `json:"contact,omitempty"`
	License        *LicenseInfoObject `json:"license,omitempty"`
}

type ContactInfoObject struct {
	Name  string `json:"name,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type LicenseInfoObject struct {
	Name string `json:"name"` // required by OpenAPI
	Url  string `json:"url,omitempty"`
}

// -------- Path Object --------
type PathObject map[string]*PathItemObject

// -------- Path Item Object --------
type PathItemObject struct {
	Ref    string           `json:"$ref,omitempty"` // OpenAPI uses "$ref"
	Post   *OperationObject `json:"post,omitempty"`
	Get    *OperationObject `json:"get,omitempty"`
	Put    *OperationObject `json:"put,omitempty"`
	Delete *OperationObject `json:"delete,omitempty"`
}

// -------- Operation Object --------
type OperationObject struct {
	Tags        []string                   `json:"tags,omitempty"`
	Summary     string                     `json:"summary,omitempty"`
	Description string                     `json:"description,omitempty"`
	OperationID string                     `json:"operationId,omitempty"` // camelCase fix
	Consumes    []string                   `json:"consumes,omitempty"`    // technically OpenAPI 2.0, not 3.0
	Produces    []string                   `json:"produces,omitempty"`    // technically OpenAPI 2.0
	Parameters  []*ParameterObject         `json:"parameters,omitempty"`
	RequestBody *RequestBodyObject         `json:"requestBody,omitempty"`
	Schemes     []string                   `json:"schemes,omitempty"` // OpenAPI 2.0 field, not 3.0
	Deprecated  bool                       `json:"deprecated,omitempty"`
	Security    []map[string][]string      `json:"security,omitempty"`
	Responses   map[string]*ResponseObject `json:"responses"` // required
}

// -------- Parameter Object --------
type ParameterObject struct {
	Name        string        `json:"name"` // required by OpenAPI
	In          string        `json:"in"`   // required by OpenAPI ("query", "header", "path", or "cookie")
	Description string        `json:"description,omitempty"`
	Default     string        `json:"default,omitempty"`
	Required    bool          `json:"required,omitempty"`
	Format      string        `json:"format,omitempty"`
	Schema      *SchemaObject `json:"schema,omitempty"`
}

type RequestBodyObject struct {
	Required    bool                        `json:"required,omitempty"`
	Description string                      `json:"description,omitempty"`
	Content     map[string]*MediaTypeObject `json:"content,omitempty"`
}

type MediaTypeObject struct {
	Schema  *SchemaObject `json:"schema,omitempty"`
	Example any           `json:"example,omitempty"`
}

type SchemasObject struct {
	Type       string                   `json:"type,omitempty"`
	Required   []string                 `json:"required,omitempty"`
	Properties map[string]*SchemaObject `json:"properties,omitempty"`
}

type SchemaObject struct {
	Type       string                   `json:"type,omitempty"`
	Required   []string                 `json:"required,omitempty"`
	Ref        string                   `json:"$ref,omitempty"` // Use $ref per OpenAPI spec
	Example    any                      `json:"example,omitempty"`
	Format     string                   `json:"format,omitempty"`
	Enum       []string                 `json:"enum,omitempty"`
	Items      *ItemsObject             `json:"items,omitempty"`
	Properties map[string]*SchemaObject `json:"properties,omitempty"`
}

type ResponseObject struct {
	Description string        `json:"description,omitempty"`
	Schema      *SchemaObject `json:"schema,omitempty"`
}

type ItemsObject struct {
	Type     string   `json:"type,omitempty"`
	Format   string   `json:"format,omitempty"`
	Required string   `json:"required,omitempty"` // consider using bool or []string?
	Enum     []string `json:"enum,omitempty"`
}

// -------- Security Scheme Object --------
type SecuritySchemeObject struct {
	Type         string `json:"type,omitempty"`
	Description  string `json:"description,omitempty"`
	Name         string `json:"name,omitempty"`
	In           string `json:"in,omitempty"`
	Scheme       string `json:"scheme,omitempty"`
	BearerFormat string `json:"bearerFormat,omitempty"`
	Flow         string `json:"flow,omitempty"` // Possibly incomplete; OAuth2 may need flows object
}

type HeaderObject struct {
	Description string   `json:"description,omitempty"`
	Type        string   `json:"type,omitempty"`
	Format      string   `json:"format,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

type ComponentObject struct {
	Schemas         map[string]*SchemaObject         `json:"schemas,omitempty"`
	SecuritySchemes map[string]*SecuritySchemeObject `json:"securitySchemes,omitempty"`
}

type ServerVariableObject struct {
	Enum        []string `json:"enum,omitempty"` // Should be an array
	Default     string   `json:"default"`
	Description string   `json:"description,omitempty"`
}

type ServerObject struct {
	Url         string                           `json:"url"`
	Description string                           `json:"description,omitempty"`
	Variables   map[string]*ServerVariableObject `json:"variables,omitempty"`
}

type SpecBuilder struct {
	Openapi    string           `json:"openapi"`
	Info       *InfoObject      `json:"info"`
	Schemes    []string         `json:"schemes,omitempty"`
	Produces   []string         `json:"produces,omitempty"`
	Consumes   []string         `json:"consumes,omitempty"`
	Servers    []*ServerObject  `json:"servers,omitempty"`
	Paths      PathObject       `json:"paths"`
	Components *ComponentObject `json:"components,omitempty"`
}

type Config struct {
	PersistAuthorization bool
}
