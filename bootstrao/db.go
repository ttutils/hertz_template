package bootstrao

import (
	"hertz_template/biz/dbmodel"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// 自动迁移表结构
	if err := db.AutoMigrate(
		&dbmodel.User{},
		&dbmodel.Book{},
	); err != nil {
		return err
	}

	err := InitData(db)
	if err != nil {
		return err
	}

	return nil
}
