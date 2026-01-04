package config

/**
 * @Description: 加载配置文件
 */
import (
	"homework4/pkg/logger"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `yaml:"app" mapstructure:"app"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
	Redis    RedisConfig    `yaml:"redis" mapstructure:"redis"`
	JWT      JWTConfig      `yaml:"jwt" mapstructure:"jwt"`
}

type AppConfig struct {
	Port int `yaml:"port" mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	User     string `yaml:"user" mapstructure:"user"`
	Password string `yaml:"password" mapstructure:"password"`
	DBName   string `yaml:"dbname" mapstructure:"dbname"`
}

type RedisConfig struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     int    `yaml:"port" mapstructure:"port"`
	Password string `yaml:"password" mapstructure:"password"`
	DB       int    `yaml:"db" mapstructure:"db"`
}

type JWTConfig struct {
	Secret  string `yaml:"secret" mapstructure:"secret"`
	Expires uint8  `yaml:"expires" mapstructure:"expires"`
}

var Cfg *Config

// 加载配置文件
func LoadConfig() {
	_, filename, _, _ := runtime.Caller(0)
	configPath := filepath.Join(filepath.Dir(filename), "config.yaml")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		logger.AppLog.Fatal("加载配置文件失败", logger.WrapMeta(err)...)
	}
	Cfg = &Config{}
	if err := viper.Unmarshal(&Cfg); err != nil {
		logger.AppLog.Fatal("解析配置文件失败", logger.WrapMeta(err)...)
	}
	if Cfg.App.Port == 0 {
		logger.AppLog.Fatal("配置信息应用端口为空，请检查配置文件")
	}
	if Cfg.Database.Host == "" {
		logger.AppLog.Fatal("配置信息数据库主机为空，请检查配置文件")
	}
	if Cfg.Database.Port == 0 {
		logger.AppLog.Fatal("配置信息数据库端口为空，请检查配置文件")
	}
	if Cfg.Database.User == "" {
		logger.AppLog.Fatal("配置信息数据库用户名为空，请检查配置文件")
	}
	if Cfg.Database.Password == "" {
		logger.AppLog.Fatal("配置信息数据库密码为空，请检查配置文件")
	}
	if Cfg.Database.DBName == "" {
		logger.AppLog.Fatal("配置信息数据库名为空，请检查配置文件")
	}
	if Cfg.Redis.Host == "" {
		logger.AppLog.Fatal("配置信息Redis主机为空，请检查配置文件")
	}
	if Cfg.Redis.Port == 0 {
		logger.AppLog.Fatal("配置信息Redis端口为空，请检查配置文件")
	}
	if Cfg.Redis.Password == "" {
		logger.AppLog.Fatal("配置信息Redis密码为空，请检查配置文件")
	}
	logger.AppLog.Info("配置文件加载成功")
}
