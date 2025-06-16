package dbmodel

type Book struct {
	ID      uint    `gorm:"primarykey"`
	Title   string  `gorm:"type:varchar(255);not null"`
	Author  string  `gorm:"type:varchar(255);not null"`
	Year    string  `gorm:"type:varchar(7);not null"`
	Summary *string `gorm:"type:text"`
}

func (Book) TableName() string {
	return "book"
}
