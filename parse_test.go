package swagger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

func Test_ScanQuery(t *testing.T) {
	type FilterUser struct {
		Name   string `query:"name" validate:"required"`
		Age    int    `query:"age"`
		Email  string `query:"email"`
		Search string `query:"search"`
	}

	asrt := assert.New(t)

	queries := ScanQuery(&FilterUser{}, core.InQuery)
	asrt.NotNil(queries)
	asrt.Len(queries, 4)
	asrt.Equal("name", queries[0].Name)
	// asrt.Equal("string", queries[0].Type)
	asrt.Equal(true, queries[0].Required)
	asrt.Equal("age", queries[1].Name)
	// asrt.Equal("integer", queries[1].Type)
	asrt.Equal(false, queries[1].Required)
	asrt.Equal("email", queries[2].Name)
	// asrt.Equal("string", queries[2].Type)
	asrt.Equal(false, queries[2].Required)
	asrt.Equal("search", queries[3].Name)
	// asrt.Equal("string", queries[3].Type)
	asrt.Equal(false, queries[3].Required)

	type Params struct {
		ID string `path:"id" validate:"required"`
	}
	param := ScanQuery(&Params{}, core.InPath)
	asrt.NotNil(param)
	asrt.Len(param, 1)
	asrt.Equal("id", param[0].Name)
	// asrt.Equal("string", param[0].Type)
	asrt.Equal(true, param[0].Required)
}

func Test_ParseDefinition(t *testing.T) {
	type User struct {
		Name string `validate:"required" example:"abc"`
		Age  int    `example:"12"`
	}
	dto := &User{}
	defintion := ParseSchema(dto)

	asrt := assert.New(t)

	fmt.Printf("%+v\n", defintion.Properties)

	asrt.Equal("object", defintion.Type)
	asrt.Equal(2, len(defintion.Properties))
	asrt.NotNil(defintion.Properties["name"])
	asrt.Equal("string", defintion.Properties["name"].Type)
	asrt.Equal("abc", defintion.Properties["name"].Example)
	asrt.NotNil(defintion.Properties["age"])
	asrt.Equal("integer", defintion.Properties["age"].Type)
	asrt.Equal("12", defintion.Properties["age"].Example)
}
