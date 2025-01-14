package swagger

import "github.com/tinh-tinh/tinhtinh/v2/core"

const OK_RESPONSE = "ok_response"

func ApiOkResponse(val interface{}) *core.Metadata {
	return core.SetMetadata(OK_RESPONSE, val)
}
