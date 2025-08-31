package cron

import (
	"hertz_template/biz/dal"
	"hertz_template/bootstrao"
	"hertz_template/utils/config"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// CleanupTask 数据库清理任务
func CleanupTask() {
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(config.Cfg.Server.DeleteDataCron, func() {
		performCleanup()
		if err := bootstrao.Migrate(dal.DB); err != nil {
			hlog.Errorf("初始化数据失败: %v", err)
		}
	})
	if err != nil {
		hlog.Errorf("添加定时任务失败: %v", err)
		return
	}

	c.Start()
	hlog.Info("CleanupTask 定时任务已启动")

	// 阻塞主线程，防止退出
	select {}
}

// performCleanup 执行数据库清理
func performCleanup() {
	start := time.Now()

	// 使用事务确保数据一致性
	err := dal.DB.Transaction(func(tx *gorm.DB) error {
		// 获取所有表名
		tables, err := dal.DB.Migrator().GetTables()
		if err != nil {
			return err
		}

		for _, table := range tables {
			err := tx.Migrator().DropTable(table)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		hlog.Errorf("数据库清理失败: %v", err)
	} else {
		elapsed := time.Since(start)
		hlog.Infof("数据库清理完成，耗时: %v", elapsed)
	}
}
