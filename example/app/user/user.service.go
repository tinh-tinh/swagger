package user

import (
	"github.com/tinh-tinh/swagger/example/app/user/dto"
	"github.com/tinh-tinh/tinhtinh/core"
)

type CrudService struct {
}

const USER_SERVICE core.Provide = "UserService"

func service(module *core.DynamicModule) *core.DynamicProvider {
	userSv := module.NewProvider(CrudService{}, USER_SERVICE)

	return userSv
}

func (s *CrudService) GetAll() []User {
	var user []User

	return user
}

func (s *CrudService) Create(input dto.SignUpUser) *User {
	newUser := User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     "user",
		Active:   true,
	}

	return &newUser
}
