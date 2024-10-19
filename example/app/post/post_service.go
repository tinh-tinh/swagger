package post

import "github.com/tinh-tinh/tinhtinh/core"

const POTS_SERVICE core.Provide = "PostService"

type PostService struct {
}

func service(module *core.DynamicModule) *core.DynamicProvider {
	postSv := module.NewProvider(core.ProviderOptions{
		Name:  POTS_SERVICE,
		Value: &PostService{},
	})

	return postSv
}
