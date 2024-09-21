package user

type User struct {
	Name     string `gorm:"type:varchar(100)"`
	Email    string `gorm:"type:varchar(100);uniqueIndex:idx_email"`
	Password string `gorm:"type:varchar(255)"`
	Role     string `gorm:"type:varchar(100)"`
	Active   bool   `gorm:"type:bool"`
}
