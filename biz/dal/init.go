package dal

import (
	"hertz_template/biz/dal/mysql"
	"hertz_template/biz/dal/postgres"
	"hertz_template/biz/dal/sqlite"
	"hertz_template/bootstrao"
	"hertz_template/utils/config"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	dbType := config.Cfg.Db.Type

	hlog.Infof("当前数据库为%s", dbType)

	var gormLogger logger.Interface
	if config.Cfg.Server.LogLevel != "debug" {
		gormLogger = logger.Default.LogMode(logger.Error) // 只有错误日志
	} else {
		gormLogger = logger.Default.LogMode(logger.Info) // 输出信息级别的日志
	}

	switch dbType {
	case "mysql":
		DB = mysql.Init(config.Cfg.Db.User, config.Cfg.Db.Password, config.Cfg.Db.Host, config.Cfg.Db.Port, config.Cfg.Db.Database, gormLogger)
		err := bootstrao.Migrate(DB)
		if err != nil {
			return
		}
	case "postgres":
		DB = postgres.Init(config.Cfg.Db.User, config.Cfg.Db.Password, config.Cfg.Db.Host, config.Cfg.Db.Port, config.Cfg.Db.Database, gormLogger)
		err := bootstrao.Migrate(DB)
		if err != nil {
			return
		}
	case "sqlite3":
		DB = sqlite.Init(config.Cfg.Db.Database, gormLogger)
		err := bootstrao.Migrate(DB)
		if err != nil {
			return
		}
	}

}
