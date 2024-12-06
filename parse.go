package swagger

import (
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/tinh-tinh/tinhtinh/core"
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
	routes := app.Module.Routers

	pathObject := make(PathObject)
	definitions := make(map[string]*DefinitionObject)

	// Parse routes
	for _, route := range routes {
		parseRoute := core.ParseRoute(route.Method + " " + route.Path)
		parseRoute.SetPrefix(route.Name)
		parameters := []*ParameterObject{}
		dtos := route.Dtos
		// Parse dto from pipe
		for _, dto := range dtos {
			val := dto.GetValue()
			switch dto.GetLocation() {
			case core.InBody:
				definitions[GetNameStruct(val)] = ParseDefinition(val)
				parameters = append(parameters, &ParameterObject{
					Name: GetNameStruct(val),
					In:   string(dto.GetLocation()),
					Schema: &SchemaObject{
						Ref: "#/definitions/" + firstLetterToLower(GetNameStruct(val)),
					},
				})
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
						Name:        file.Name,
						In:          "formData",
						Type:        "file",
						Required:    file.Required,
						Description: file.Description,
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
		res := map[string]*ResponseObject{"200": response}
		operation := &OperationObject{
			Tags:       []string{},
			Consumes:   []string{},
			Parameters: parameters,
			Responses:  res,
			Security:   []map[string][]string{},
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
		case "DELETE":
			itemObject.Delete = operation
		}
	}

	spec.Definitions = definitions
	spec.Paths = pathObject
}

type Mapper map[string]interface{}

// RecursiveParseStandardSwagger takes a struct and recursively parses its fields
// to create a swagger-style mapper. The mapper is a map[string]interface{}
// where the keys are the field names (lowercased) and the values are the
// field values. The rules for parsing the fields are as follows:
//
// - If the field is a pointer, it is recursively parsed.
// - If the field is a map, its values are recursively parsed.
// - If the field is a slice, its elements are recursively parsed.
// - If the field is a primitive type, its value is used as is.
//
// The function returns a Mapper or nil if the input is nil.
func RecursiveParseStandardSwagger(val interface{}) Mapper {
	mapper := make(Mapper)

	if reflect.ValueOf(val).IsNil() {
		return nil
	}
	ct := reflect.ValueOf(val).Elem()
	for i := 0; i < ct.NumField(); i++ {
		field := ct.Type().Field(i)
		key := firstLetterToLower(field.Name)
		if key == "ref" {
			key = "$ref"
		}
		if ct.Field(i).Interface() == nil {
			continue
		}
		if field.Type.Kind() == reflect.Pointer {
			ptrVal := RecursiveParseStandardSwagger(ct.Field(i).Interface())
			if len(ptrVal) == 0 {
				continue
			}
			mapper[key] = ptrVal
		} else if field.Type.Kind() == reflect.Map {
			val := ct.Field(i).Interface()
			mapVal := reflect.ValueOf(val)
			subMapper := make(Mapper)
			for _, v := range mapVal.MapKeys() {
				subVal := RecursiveParseStandardSwagger(mapVal.MapIndex(v).Interface())
				if IsNil(subVal) {
					continue
				}
				subKey := firstLetterToLower(v.String())
				subMapper[subKey] = subVal
			}
			mapper[key] = subMapper
		} else if field.Type.Kind() == reflect.Slice {
			arrVal := reflect.ValueOf(ct.Field(i).Interface())
			if arrVal.IsValid() {
				arr := []interface{}{}
				for i := 0; i < arrVal.Len(); i++ {
					item := arrVal.Index(i)
					if item.Kind() == reflect.Pointer {
						arr = append(arr, RecursiveParseStandardSwagger(item.Interface()))
					} else {
						arr = append(arr, item.Interface())
					}
				}
				if IsNil(arr) {
					continue
				}
				mapper[key] = arr
			}
		} else {
			val := ct.Field(i).Interface()
			if IsNil(val) {
				continue
			}
			mapper[key] = ct.Field(i).Interface()
		}
	}

	if len(mapper) == 0 {
		return nil
	}

	return mapper
}

// ParseDefinition takes a struct and recursively parses its fields
// to create a swagger-style DefinitionObject. The rules for parsing the fields are as follows:
//
// - If the field is a pointer, it is recursively parsed.
// - If the field is a map, its values are recursively parsed.
// - If the field is a slice, its elements are recursively parsed.
// - If the field is a primitive type, its value is used as is.
//
// The function returns a DefinitionObject or nil if the input is nil.
func ParseDefinition(dto interface{}) *DefinitionObject {
	properties := make(map[string]*SchemaObject)
	ct := reflect.ValueOf(dto).Elem()
	for i := 0; i < ct.NumField(); i++ {
		schema := &SchemaObject{
			Type: mappingType(ct.Field(i)),
		}
		if reflect.TypeOf(ct.Field(i).Interface()) == reflect.TypeOf(time.Time{}) {
			schema.Format = "date-time"
		}

		field := ct.Type().Field(i)
		validator := field.Tag.Get("validate")
		isRequired := slices.IndexFunc(strings.Split(validator, ","), func(v string) bool { return v == "required" })
		if isRequired == -1 {
			schema.Required = false
		} else {
			schema.Required = true
		}
		example := field.Tag.Get("example")
		if example != "" {
			schema.Example = example
		}

		properties[ct.Type().Field(i).Name] = schema
	}

	return &DefinitionObject{
		Type:       "object",
		Properties: properties,
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
func ScanQuery(val interface{}, in core.InDto) []*ParameterObject {
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
			Type: mappingType(ct.Field(i)),
			In:   string(in),
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
