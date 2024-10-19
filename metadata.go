package swagger

import "github.com/tinh-tinh/tinhtinh/core"

const TAG = "tag"

func Tag(names ...string) *core.Metadata {
	return core.SetMetadata(TAG, names)
}

const SECURITY = "security"

func Security(names ...string) *core.Metadata {
	return core.SetMetadata(SECURITY, names)
}
