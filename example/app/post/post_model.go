package post

type Post struct {
	Title   string `gorm:"type:varchar(255)"`
	Content string `gorm:"type:text"`
}
