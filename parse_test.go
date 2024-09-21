package swagger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/tinhtinh/core"
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
	asrt.Equal("string", queries[0].Type)
	asrt.Equal(true, queries[0].Required)
	asrt.Equal("age", queries[1].Name)
	asrt.Equal("integer", queries[1].Type)
	asrt.Equal(false, queries[1].Required)
	asrt.Equal("email", queries[2].Name)
	asrt.Equal("string", queries[2].Type)
	asrt.Equal(false, queries[2].Required)
	asrt.Equal("search", queries[3].Name)
	asrt.Equal("string", queries[3].Type)
	asrt.Equal(false, queries[3].Required)

	type Params struct {
		ID string `path:"id" validate:"required"`
	}
	param := ScanQuery(&Params{}, core.InPath)
	asrt.NotNil(param)
	asrt.Len(param, 1)
	asrt.Equal("id", param[0].Name)
	asrt.Equal("string", param[0].Type)
	asrt.Equal(true, param[0].Required)
}

func Test_ParseDefinition(t *testing.T) {
	type User struct {
		Name string `validate:"required" example:"abc"`
		Age  int    `example:"12"`
	}
	dto := &User{}
	defintion := ParseDefinition(dto)

	asrt := assert.New(t)

	asrt.Equal("object", defintion.Type)
	asrt.Equal(2, len(defintion.Properties))
	asrt.Equal("string", defintion.Properties["Name"].Type)
	asrt.Equal("integer", defintion.Properties["Age"].Type)
	asrt.Equal(true, defintion.Properties["Name"].Required)
	asrt.Equal(false, defintion.Properties["Age"].Required)
	asrt.Equal("abc", defintion.Properties["Name"].Example)
	asrt.Equal("12", defintion.Properties["Age"].Example)
}

func Test_recursive(t *testing.T) {
	type Children struct {
		Name string
	}
	child := &Children{
		Name: "abc",
	}
	child2 := &Children{
		Name: "def",
	}
	type Parent struct {
		Items []*Children
	}
	parent := &Parent{
		Items: []*Children{child, child2},
	}
	parent2 := &Parent{}
	mapper1 := recursiveParseStandardSwagger(parent)
	mapper2 := recursiveParseStandardSwagger(parent2)

	asrt := assert.New(t)
	asrt.Equal(1, len(mapper1))
	asrt.Equal(0, len(mapper2))

	fmt.Println(mapper1)
	fmt.Println(mapper2)
}
