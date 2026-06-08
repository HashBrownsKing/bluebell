package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 使用结构体变量保存配置信息
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Port      int    `mapstructure:"port"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxOpenConnects int    `mapstructure:"max_open_connects"`
	MaxIdleConnects int    `mapstructure:"max_idle_connects"`
	MaxLifeTime     int    `mapstructure:"max_life_time"`
	MaxIdleTime     int    `mapstructure:"max_idle_time"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
	Password string `mapstructure:"password"`
}

var Config = new(AppConfig)

func Init() (err error) {
	// 指定配置文件方式1
	viper.SetConfigFile("config.yaml") // 直接指定配置文件为config.yaml

	/*
		// 指定配置文件方式2
		viper.SetConfigName("config")      // 仅指定文件名（不带后缀），必须配合 AddConfigPath() 使用，让 Viper 自动查找 config.yaml/config.yml/config.json。
		// SetConfigType 指定配置文件的解析类型（yaml/json/ini等）
		// 作用：告诉Viper使用哪种格式解析配置
		// 适用场景：1. 配置文件无后缀名 2. 配合SetConfigName自动查找配置 3. 加载远程配置
		// 注意：使用SetConfigFile指定带后缀的本地文件时，无需调用此方法
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	*/
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed, err: %v\n", err)
		return
	}

	if err = viper.Unmarshal(Config); err != nil { // 把参数反序列化配置到结构体
		fmt.Printf("viper.Unmarshal() failed, err: %v\n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config.yaml changed...")
		if err = viper.Unmarshal(Config); err != nil { // 把参数反序列化配置到结构体
			fmt.Printf("viper.Unmarshal() failed, err: %v\n", err)
		}
	})
	return
}
