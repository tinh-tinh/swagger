package swagger

import "github.com/tinh-tinh/tinhtinh/core"

const TAG = "tag"

func ApiTag(names ...string) *core.Metadata {
	return core.SetMetadata(TAG, names)
}

const SECURITY = "security"

func ApiSecurity(names ...string) *core.Metadata {
	return core.SetMetadata(SECURITY, names)
}

const CONSUMER = "consumer"

func ApiConsumer(names ...string) *core.Metadata {
	return core.SetMetadata(CONSUMER, names)
}

const FILE = "file"

type FileOptions struct {
	Name        string
	Required    bool
	Description string
}

func ApiFile(opts ...FileOptions) *core.Metadata {
	return core.SetMetadata(FILE, opts)
}
