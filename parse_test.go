package swagger

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Name   string `json:"name,omitempty" validate:"required" example:"abc"`
		Age    int    `json:"age,omitempty" example:"12"`
		Hidden int    `hidden:"true"`
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
	asrt.Nil(defintion.Properties["hidden"])
}

type CreateTimeoffTypeInput struct {
	Name         string             `json:"name" example:"Annual Leave"`
	Country      string             `json:"country" example:"US"`
	LocationID   string             `json:"locationId,omitempty" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Category     string             `json:"category" example:"paid-time-off"`
	Config       *UpsertConfigInput `json:"config" validate:"nested"`
	RequiredInfo *RequiredInfoInput `json:"requiredInfo,omitempty" validate:"nested"`
}

type RequiredInfoInput struct {
	Gender       string `json:"gender" example:"male"`
	EmployeeType string `json:"employeeType" example:"full-time"`
}

type UpsertConfigInput struct {
	AllowedApplyFuture     bool                         `json:"allowedApplyFuture" example:"true"`
	TimeUnit               string                       `json:"timeUnit" example:"d"`
	MaxLeaveAmount         float32                      `json:"maxLeaveAmount" example:"5"`
	MinLeaveAmount         float32                      `json:"minLeaveAmount" example:"0.5"`
	LeaveApplicationStart  uint                         `json:"leaveApplicationStart" example:"60"`
	AnnualResetPolicy      *AnnualResetPolicyInput      `json:"annualResetPolicy" validate:"nested"`
	AccrualPolicy          *AccrualPolicyInput          `json:"accrualPolicy" validate:"nested"`
	CarryForwardPolicy     *CarryForwardPolicyInput     `json:"carryForwardPolicy" validate:"nested"`
	NewHireProbationPolicy *NewHireProbationPolicyInput `json:"newHireProbationPolicy" validate:"nested"`
	AutoApproval           *AutoApprovalInput           `json:"autoApproval" validate:"nested"`
	EmailReminder          *EmailReminderInput          `json:"emailReminder" validate:"nested"`
}

type AnnualResetPolicyInput struct {
	Type string `json:"type" example:"calendarDate"`
	Date string `json:"date,omitempty" example:"2024-01-01"`
}

type AccrualPolicyInput struct {
	Method string           `json:"accrualMethod" example:"year"`
	Rates  []RateRangeInput `json:"accrualRates" validate:"nested"`
}

type CarryForwardPolicyInput struct {
	ExpireDuration int              `json:"expireDuration" example:"90"`
	Rates          []RateRangeInput `json:"carryForwardRates,omitempty" validate:"nested"`
}

type NewHireProbationPolicyInput struct {
	IsEnable bool             `json:"isEnable" example:"false"`
	Rules    []RateRangeInput `json:"rules,omitempty" validate:"nested"`
}

type RateRangeInput struct {
	From  uint `json:"from" example:"0"`
	To    uint `json:"to" example:"100"`
	Value uint `json:"value" example:"12"`
}

type AutoApprovalInput struct {
	IsEnable       bool    `json:"isEnable" example:"true"`
	ExpireDuration uint    `json:"expireDuration,omitempty" example:"72"`
	LeaveAmount    float32 `json:"leaveAmount,omitempty" example:"3"`
}

type EmailReminderInput struct {
	IsEnable       bool `json:"isEnable" example:"true"`
	ExpireDuration uint `json:"expireDuration,omitempty" example:"24"`
}

func Test_ComplexStruct(t *testing.T) {

	dto := &CreateTimeoffTypeInput{}
	defintion := ParseSchema(dto)

	text, err := json.Marshal(defintion)
	require.Nil(t, err)
	require.Equal(t, `{"type":"object","properties":{"category":{"type":"string","example":"paid-time-off"},"config":{"type":"object","properties":{"accrualPolicy":{"type":"object","properties":{"accrualMethod":{"type":"string","example":"year"},"accrualRates":{"type":"array","items":{"type":"object","properties":{"from":{"type":"integer","example":"0"},"to":{"type":"integer","example":"100"},"value":{"type":"integer","example":"12"}}}}}},"allowedApplyFuture":{"type":"boolean","example":"true"},"annualResetPolicy":{"type":"object","properties":{"date":{"type":"string","example":"2024-01-01"},"type":{"type":"string","example":"calendarDate"}}},"autoApproval":{"type":"object","properties":{"expireDuration":{"type":"integer","example":"72"},"isEnable":{"type":"boolean","example":"true"},"leaveAmount":{"type":"number","example":"3"}}},"carryForwardPolicy":{"type":"object","properties":{"carryForwardRates":{"type":"array","items":{"type":"object","properties":{"from":{"type":"integer","example":"0"},"to":{"type":"integer","example":"100"},"value":{"type":"integer","example":"12"}}}},"expireDuration":{"type":"integer","example":"90"}}},"emailReminder":{"type":"object","properties":{"expireDuration":{"type":"integer","example":"24"},"isEnable":{"type":"boolean","example":"true"}}},"leaveApplicationStart":{"type":"integer","example":"60"},"maxLeaveAmount":{"type":"number","example":"5"},"minLeaveAmount":{"type":"number","example":"0.5"},"newHireProbationPolicy":{"type":"object","properties":{"isEnable":{"type":"boolean","example":"false"},"rules":{"type":"array","items":{"type":"object","properties":{"from":{"type":"integer","example":"0"},"to":{"type":"integer","example":"100"},"value":{"type":"integer","example":"12"}}}}}},"timeUnit":{"type":"string","example":"d"}}},"country":{"type":"string","example":"US"},"locationId":{"type":"string","example":"3fa85f64-5717-4562-b3fc-2c963f66afa6"},"name":{"type":"string","example":"Annual Leave"},"requiredInfo":{"type":"object","properties":{"employeeType":{"type":"string","example":"full-time"},"gender":{"type":"string","example":"male"}}}}}`, string(text))
}
