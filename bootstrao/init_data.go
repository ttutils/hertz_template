package bootstrao

import (
	"hertz_template/biz/dbmodel"
	"hertz_template/utils"
	"hertz_template/utils/config"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

func InitData(db *gorm.DB) error {
	// 插入初始化账号
	var count int64
	if err := db.Model(&dbmodel.User{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		return err
	}

	// 如果不存在则创建
	if count == 0 {
		hlog.Infof("%s 用户不存在，密码为:%s", config.Cfg.Admin.Username, config.Cfg.Admin.Password)
		adminUser := &dbmodel.User{
			Username: config.Cfg.Admin.Username,
			Password: utils.MD5(config.Cfg.Admin.Password),
		}
		if err := db.Create(adminUser).Error; err != nil {
			return err
		}
	}

	return nil
}
