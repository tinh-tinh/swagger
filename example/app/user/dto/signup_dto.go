package dto

import "time"

type SignUpUser struct {
	Name     string    `validate:"isAlpha" example:"John"`
	Email    string    `validate:"required,isEmail" example:"john@gmail.com"`
	Password string    `validate:"required,isStrongPassword" example:"12345678@Tc"`
	Birth    time.Time `validate:"required" example:"2024-12-12"`
}

type FindUser struct {
	Name string `validate:"required,isAlpha" query:"name" example:"ac"`
	Age  uint   `validate:"required,isInt" query:"age"`
}
