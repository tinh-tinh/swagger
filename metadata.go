package swagger

import "github.com/tinh-tinh/tinhtinh/v2/core"

const TAG = "openapi_tag"

func ApiTag(names ...string) *core.Metadata {
	return core.SetMetadata(TAG, names)
}

const DESCRIPTION = "openapi_description"

func ApiDescription(description string) *core.Metadata {
	return core.SetMetadata(DESCRIPTION, description)
}

const SUMMARY = "openapi_summary"

func ApiSummary(summary string) *core.Metadata {
	return core.SetMetadata(SUMMARY, summary)
}

const SECURITY = "openapi_security"

func ApiSecurity(names ...string) *core.Metadata {
	return core.SetMetadata(SECURITY, names)
}

const CONSUMER = "openapi_consumer"

func ApiConsumer(names ...string) *core.Metadata {
	return core.SetMetadata(CONSUMER, names)
}

const FILE = "openapi_file"

type FileOptions struct {
	Name        string
	Required    bool
	Description string
}

func ApiFile(opts ...FileOptions) *core.Metadata {
	return core.SetMetadata(FILE, opts)
}
