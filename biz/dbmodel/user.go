package dbmodel

type User struct {
	ID       uint    `gorm:"primarykey"`
	Username string  `gorm:"unique;column:username;type:varchar(255)"`
	Password string  `gorm:"column:password;type:varchar(255)"`
	Email    *string `gorm:"column:email;type:varchar(255)"`
}

func (u *User) TableName() string {
	return "user"
}
