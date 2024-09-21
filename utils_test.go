package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_firstLetterToLower(t *testing.T) {
	asrt := assert.New(t)

	asrt.Equal("aBCD", firstLetterToLower("ABCD"))
	asrt.Equal("a123dFVV", firstLetterToLower("A123dFVV"))
	asrt.Equal("l#$V#vdfDVG", firstLetterToLower("L#$V#vdfDVG"))
	asrt.Equal("c VFfvfv", firstLetterToLower("C VFfvfv"))
}

func Test_IsNil(t *testing.T) {
	asrt := assert.New(t)
	asrt.True(IsNil(""))
	asrt.True(IsNil([]string{}))
	asrt.True(IsNil([]*interface{}{}))
	abc := make(map[string]interface{})
	asrt.True(IsNil(abc))
}
