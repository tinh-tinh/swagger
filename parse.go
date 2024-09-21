package swagger

import (
	"reflect"
	"slices"
	"strings"

	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/utils"
)

func (spec *SpecBuilder) ParsePaths(app *core.App) {
	// mapperDoc := app.Module.MapperDoc
	routes := app.Module.Routers

	pathObject := make(PathObject)
	definitions := make(map[string]*DefinitionObject)

	for _, route := range routes {
		parseRoute := core.ParseRoute(route.Path)
		parameters := []*ParameterObject{}
		dtos := route.Dtos
		for _, dto := range dtos {
			switch dto.In {
			case core.InBody:
				definitions[utils.GetNameStruct(dto.Dto)] = ParseDefinition(dto.Dto)
				parameters = append(parameters, &ParameterObject{
					Name: utils.GetNameStruct(dto.Dto),
					In:   string(dto.In),
					Schema: &SchemaObject{
						Ref: "#/definitions/" + firstLetterToLower(utils.GetNameStruct(dto.Dto)),
					},
				})
			case core.InQuery:
				parameters = append(parameters, ScanQuery(dto.Dto, dto.In)...)
			case core.InPath:
				parameters = append(parameters, ScanQuery(dto.Dto, dto.In)...)
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
			Tags:       []string{route.Tag},
			Parameters: parameters,
			Responses:  res,
			Security:   []map[string][]string{},
		}

		if len(route.Security) > 0 {
			security := map[string][]string{}
			for _, s := range route.Security {
				security[s] = []string{}
			}
			operation.Security = append(operation.Security, security)
		}
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

func recursiveParseStandardSwagger(val interface{}) Mapper {
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
			ptrVal := recursiveParseStandardSwagger(ct.Field(i).Interface())
			if len(ptrVal) == 0 {
				continue
			}
			mapper[key] = ptrVal
		} else if field.Type.Kind() == reflect.Map {
			val := ct.Field(i).Interface()
			mapVal := reflect.ValueOf(val)
			subMapper := make(Mapper)
			for _, v := range mapVal.MapKeys() {
				subVal := recursiveParseStandardSwagger(mapVal.MapIndex(v).Interface())
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
						arr = append(arr, recursiveParseStandardSwagger(item.Interface()))
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

func ParseDefinition(dto interface{}) *DefinitionObject {
	properties := make(map[string]*SchemaObject)
	ct := reflect.ValueOf(dto).Elem()
	for i := 0; i < ct.NumField(); i++ {
		schema := &SchemaObject{
			Type: ct.Field(i).Kind().String(),
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
			Type: field.Type.Name(),
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
