package swagger

// -------- Info Object --------
type InfoObject struct {
	Title          string
	Description    string
	Version        string
	TermsOfService string
	Contact        *ContactInfoObject
	License        *LicenseInfoObject
}

type ContactInfoObject struct {
	Name  string
	Url   string
	Email string
}

type LicenseInfoObject struct {
	Name string
	Url  string
}

// -------- Path Object --------
type PathObject map[string]*PathItemObject

// -------- Path Item Object --------
type PathItemObject struct {
	Ref    string
	Post   *OperationObject
	Get    *OperationObject
	Put    *OperationObject
	Delete *OperationObject
}

// -------- Operation Object --------
type OperationObject struct {
	Tags        []string
	Summary     string
	Description string
	OperationID string
	Consumes    []string
	Produces    []string
	Parameters  []*ParameterObject
	RequestBody *RequestBodyObject
	Schemes     []string
	Deprecated  bool
	Security    []map[string][]string
	Responses   map[string]*ResponseObject
}

// -------- Parameter Object --------
type ParameterObject struct {
	Name        string
	In          string
	Description string
	Default     string
	Required    bool
	Format      string
	Schema      *SchemaObject
}

type RequestBodyObject struct {
	Required    bool
	Description string
	Content     map[string]*MediaTypeObject
}

type MediaTypeObject struct {
	Schema  *SchemaObject
	Example any
}

// -------- Component Object --------
type SchemasObject struct {
	Type       string
	Required   []string
	Properties map[string]*SchemaObject
}

// -------- Schema Object --------
type SchemaObject struct {
	Type       string
	Required   []string
	Ref        string
	Example    string
	Format     string
	Enum       []string
	Items      *ItemsObject
	Properties map[string]*SchemaObject
}

// -------- Response Object --------
type ResponseObject struct {
	Description string
	Schema      *SchemaObject
}

// -------- Items Object --------
type ItemsObject struct {
	Type     string
	Format   string
	Required string
	Enum     []string
}

// -------- Security Scheme Object --------
type SecuritySchemeObject struct {
	Type         string
	Description  string
	Name         string
	In           string
	Scheme       string
	BearerFormat string
	Flow         string
}

// -------- Header Object --------
type HeaderObject struct {
	Description string
	Type        string
	Format      string
	Enum        []string
}

type ComponentObject struct {
	Schemas         map[string]*SchemaObject
	SecuritySchemes map[string]*SecuritySchemeObject
}

type ServerVariableObject struct {
	Enum        string
	Default     string
	Description string
}

type ServerObject struct {
	Url         string
	Description string
	Variables   map[string]*ServerVariableObject
}

type SpecBuilder struct {
	Openapi    string
	Info       *InfoObject
	Schemes    []string
	Produces   []string
	Consumes   []string
	Servers    []*ServerObject
	Paths      PathObject
	Components *ComponentObject
}
