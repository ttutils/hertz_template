package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port          int    `mapstructure:"port"`
	Name          string `mapstructure:"name"`
	LogLevel      string `mapstructure:"log_level"`
	EnableSwagger bool   `mapstructure:"swagger"`
}

type DbConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type JwtConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"`
}

type AuthConfig struct {
	ExcludedPaths []string `mapstructure:"excludedPaths"`
}

type AdminConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
	Db     DbConfig     `mapstructure:"db"`
	Jwt    JwtConfig    `mapstructure:"jwt"`
	Auth   AuthConfig   `mapstructure:"auth"`
	Admin  AdminConfig  `mapstructure:"admin"`
}

var Cfg AppConfig

func InitConfig(defaultConfigContent []byte) {
	// 1. 处理命令行参数
	var configFile string
	pflag.StringVar(&configFile, "config", "", "Path to custom config file")
	pflag.IntVar(&Cfg.Server.Port, "server.port", 8888, "Server port")
	pflag.Parse()

	v := viper.New()

	// 2. 加载嵌入的默认配置文件
	if len(defaultConfigContent) > 0 {
		v.SetConfigType("yaml")
		if err := v.ReadConfig(bytes.NewBuffer(defaultConfigContent)); err != nil {
			fmt.Printf("加载默认配置失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 3. 加载外部配置文件（如果存在）
	if configFile != "" {
		if _, err := os.Stat(configFile); err == nil {
			v.SetConfigFile(configFile)
			if err := v.MergeInConfig(); err != nil {
				fmt.Printf("加载外部配置失败: %v (路径: %s)\n", err, configFile)
				os.Exit(1)
			}
		} else {
			fmt.Printf("警告: 外部配置文件不存在，使用默认配置 (路径: %s)\n", configFile)
		}
	}

	// 4. 环境变量覆盖（支持 HERTZ_SERVER_PORT 这类变量）
	v.SetEnvPrefix("HERTZ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 5. 合并命令行参数
	_ = v.BindPFlags(pflag.CommandLine)

	// 6. 映射到结构体
	if err := v.Unmarshal(&Cfg); err != nil {
		fmt.Println("解析配置失败:", err)
		os.Exit(1)
	}
}
