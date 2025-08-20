package swagger

import (
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

// ParsePaths parse all routes in the app and create a swagger spec.
//
// This method will loop through all routes in the app and parse the route
// path, method, and dtos. It will then create a swagger spec in the
// spec.Paths and spec.Definitions fields.
//
// The rules for parsing the route path are as follows:
//   - If the route path contains a path parameter, it will be replaced with
//     {parameter_name}.
//   - If the route path contains multiple path parameters, they will be
//     replaced with {parameter_name1}/{parameter_name2}/.../{parameter_nameN}
//
// The rules for parsing the dtos are as follows:
//   - If the dto is in the body, it will be replaced with the name of the
//     dto in the definitions section.
//   - If the dto is in the query or path, it will be replaced with the name
//     of the dto in the parameters section.
func (spec *SpecBuilder) ParsePaths(app *core.App) {
	// mapperDoc := app.Module.MapperDoc
	routes := app.Module.GetRouters()

	pathObject := make(PathObject)
	schemas := make(map[string]*SchemaObject)

	// Parse routes
	for _, route := range routes {
		parseRoute := core.ParseRoute(route.Method + " " + route.Path)
		parseRoute.SetPrefix(route.Name)
		if app.Prefix != "" {
			parseRoute.SetPrefix(app.Prefix)
		}
		parameters := []*ParameterObject{}
		mediaTypes := make(map[string]*MediaTypeObject)
		dtos := route.Dtos
		// Parse dto from pipe
		for _, dto := range dtos {
			val := dto.GetValue()
			switch dto.GetLocation() {
			case core.InBody:
				schemas[common.GetStructName(val)] = ParseSchema(val)
				mediaTypes[common.GetStructName(val)] = &MediaTypeObject{
					Schema: &SchemaObject{
						Ref: "#/components/schemas/" + common.GetStructName(val),
					},
				}
			case core.InQuery:
				parameters = append(parameters, ScanQuery(val, dto.GetLocation())...)
			case core.InPath:
				parameters = append(parameters, ScanQuery(val, dto.GetLocation())...)
			}
		}

		fileIdx := slices.IndexFunc(route.Metadata, func(v *core.Metadata) bool {
			return v.Key == FILE
		})
		if fileIdx != -1 {
			files, ok := route.Metadata[fileIdx].Value.([]FileOptions)
			if ok {
				for _, file := range files {
					parameters = append(parameters, &ParameterObject{
						Name: file.Name,
						In:   "formData",
						// Type:        "file",
						Required:    file.Required,
						Description: file.Description,
						Schema: &SchemaObject{
							Type: "file",
						},
					})
				}
			}
		}

		if pathObject[parseRoute.Path] == nil {
			pathObject[parseRoute.Path] = &PathItemObject{}
		}
		itemObject := pathObject[parseRoute.Path]
		response := &ResponseObject{
			Description: "Ok",
		}

		findOkIdx := slices.IndexFunc(route.Metadata, func(v *core.Metadata) bool {
			return v.Key == OK_RESPONSE
		})

		if findOkIdx != -1 {
			res := route.Metadata[findOkIdx].Value
			schemas[common.GetStructName(res)] = ParseSchema(res)

			response.Schema = &SchemaObject{
				Ref: "#/components/schemas/" + common.GetStructName(res),
			}
		}

		res := map[string]*ResponseObject{"200": response}
		operation := &OperationObject{
			Tags:       []string{},
			Consumes:   []string{},
			Parameters: parameters,
			Responses:  res,
			Security:   []map[string][]string{},
		}

		if len(mediaTypes) > 0 {
			operation.RequestBody = &RequestBodyObject{
				Content:  mediaTypes,
				Required: true,
			}
		}

		// Api Tag
		tagIndex := slices.IndexFunc(route.Metadata, func(v *core.Metadata) bool { return v.Key == TAG })
		if tagIndex != -1 {
			tags, ok := route.Metadata[tagIndex].Value.([]string)
			if ok {
				operation.Tags = tags
			}
		}

		// Api Security
		secureIndex := slices.IndexFunc(route.Metadata, func(v *core.Metadata) bool { return v.Key == SECURITY })
		if secureIndex != -1 {
			securities, ok := route.Metadata[secureIndex].Value.([]string)
			if ok {
				security := map[string][]string{}
				for _, s := range securities {
					security[s] = []string{}
				}
				operation.Security = append(operation.Security, security)
			}
		}

		// Api Consumer
		consumerIndex := slices.IndexFunc(route.Metadata, func(v *core.Metadata) bool { return v.Key == CONSUMER })
		if consumerIndex != -1 {
			consumers, ok := route.Metadata[consumerIndex].Value.([]string)
			if ok {
				operation.Consumes = consumers
			}
		}

		// Matching method
		switch parseRoute.Method {
		case "GET":
			itemObject.Get = operation
		case "POST":
			itemObject.Post = operation
		case "PUT":
			itemObject.Put = operation
		case "PATCH":
			itemObject.Patch = operation
		case "DELETE":
			itemObject.Delete = operation
		}
	}

	// spec.Definitions = definitions
	spec.Components.Schemas = schemas
	spec.Paths = pathObject
}

type Mapper map[string]interface{}

func ParseSchema(dto any) *SchemaObject {
	properties := make(map[string]*SchemaObject)
	requiredFields := []string{}

	v := reflect.ValueOf(dto)
	if v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		fieldJsons := fieldType.Tag.Get("json")
		fieldName := ""
		if fieldJsons != "" {
			fieldName = strings.Split(fieldJsons, ",")[0]
		}
		if fieldName == "" {
			fieldName = strings.ToLower(fieldType.Name)
		}

		hiddenField := fieldType.Tag.Get("hidden")
		if hiddenField != "" {
			continue
		}

		// Skip unexported fields
		if !field.CanSet() || !field.CanInterface() {
			continue
		}

		schema := &SchemaObject{
			Type: mappingType(field.Type()),
		}

		if reflect.TypeOf(field.Interface()) == reflect.TypeOf(time.Time{}) {
			schema.Format = "date-time"
		}

		// Parse validation tag
		validator := fieldType.Tag.Get("validate")
		requiredIdx := slices.IndexFunc(strings.Split(validator, ","), func(v string) bool {
			return v == "required"
		})
		if requiredIdx != -1 {
			requiredFields = append(requiredFields, fieldName)
		}

		if schema.Type == "array" {
			elemType := field.Type().Elem()
			schema.Items = &ItemsObject{
				Type: mappingType(elemType),
			}
		}

		nestedIdx := slices.IndexFunc(strings.Split(validator, ","), func(v string) bool {
			return v == "nested"
		})
		if nestedIdx != -1 {
			if field.Kind() == reflect.Ptr && field.IsNil() {
				// Create a new instance of the pointed-to struct
				newValue := reflect.New(field.Type().Elem())
				field.Set(newValue) // Set the new instance to the field
				if schema.Properties == nil {
					schema.Properties = make(map[string]*SchemaObject)
				}
				result := ParseSchema(field.Interface())
				properties[fieldName] = result
			}
		} else {
			// Parse example tag
			example := fieldType.Tag.Get("example")
			if example != "" {
				if schema.Type == "array" {
					schema.Example = strings.Split(example, ",")
				} else {
					schema.Example = example
				}
			}

			properties[fieldName] = schema
		}
	}

	return &SchemaObject{
		Type:       "object",
		Properties: properties,
		Required:   requiredFields,
	}
}

// ScanQuery takes a struct and recursively parses its fields to create a swagger-style mapper.
// The mapper is a slice of ParameterObject where the keys are the field names (lowercased) and the values are the
// field values. The rules for parsing the fields are as follows:
//
// - If the field is a pointer, it is recursively parsed.
// - If the field is a map, its values are recursively parsed.
// - If the field is a slice, its elements are recursively parsed.
// - If the field is a primitive type, its value is used as is.
//
// The function returns a slice of ParameterObject or nil if the input is nil.
func ScanQuery(val interface{}, in core.CtxKey) []*ParameterObject {
	ct := reflect.ValueOf(val).Elem()

	params := []*ParameterObject{}
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)

		name := ""
		if in == core.InQuery {
			name = field.Tag.Get("query")
		} else {
			name = field.Tag.Get("path")
		}
		param := &ParameterObject{
			Name: name,
			// Type: mappingType(ct.Field(i)),
			Schema: &SchemaObject{
				Type: mappingType(ct.Field(i).Type()),
			},
			In: string(in),
		}
		validator := field.Tag.Get("validate")
		isRequired := slices.IndexFunc(strings.Split(validator, ","), func(v string) bool { return v == "required" })
		if isRequired == -1 {
			param.Required = false
		} else {
			param.Required = true
		}
		example := field.Tag.Get("example")
		if example != "" {
			param.Default = example
		}

		params = append(params, param)
	}

	return params
}
