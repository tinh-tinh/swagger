package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsNil(t *testing.T) {
	asrt := assert.New(t)
	asrt.True(IsNil(""))
	asrt.True(IsNil([]string{}))
	asrt.True(IsNil([]*interface{}{}))
	abc := make(map[string]interface{})
	asrt.True(IsNil(abc))
}
