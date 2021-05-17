package lib

import (
	"fmt"

	"github.com/RealLiuSha/echo-admin/errors"
	"github.com/RealLiuSha/echo-admin/pkg/file"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var configPath = "./config.yml"
var casbinModelPath = "./casbin_model.conf"

var defaultConfig = Config{
	Name: "app",
	Http: &HttpConfig{
		Host: "0.0.0.0",
		Port: 9999,
	},
	Log: &LogConfig{
		Level:       "debug",
		Directory:   "/tmp/app",
		Development: true,
	},
	SuperAdmin: &SuperAdminConfig{},
	Auth:       &AuthConfig{},
	Casbin:     &CasbinConfig{Enable: false},
	Redis:      &RedisConfig{Host: "127.0.0.1", Port: 6379},
	Database: &DatabaseConfig{
		Parameters:   "charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true&timeout=5s",
		MaxLifetime:  7200,
		MaxOpenConns: 150,
		MaxIdleConns: 50,
	},
}

func NewConfig() Config {
	config := defaultConfig

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "failed to read config"))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(errors.Wrap(err, "failed to marshal config"))
	}

	config.Casbin.Model = casbinModelPath
	return config
}

func SetConfigPath(path string) {
	if !file.IsFile(path) {
		panic("config filepath does not exist")
	}

	configPath = path
}

func SetConfigCasbinModelPath(path string) {
	if !file.IsFile(path) {
		panic("casbin model filepath does not exist")
	}

	casbinModelPath = path
}

// Configuration are the available config values
type Config struct {
	Name       string            `mapstructure:"Name"`
	Http       *HttpConfig       `mapstructure:"Http"`
	Log        *LogConfig        `mapstructure:"Log"`
	SuperAdmin *SuperAdminConfig `mapstructure:"SuperAdmin"`
	Auth       *AuthConfig       `mapstructure:"Auth"`
	Casbin     *CasbinConfig     `mapstructure:"Casbin"`
	Redis      *RedisConfig      `mapstructure:"Redis"`
	Database   *DatabaseConfig   `mapstructure:"Database"`
}

type HttpConfig struct {
	Host string `mapstructure:"Host" validate:"ipv4"`
	Port int    `mapstructure:"Port" validate:"gte=1,lte=65535"`
}

// LogLevel     : debug,info,warn,error,dpanic,panic,fatal
//                default info
// Format       : json, console
//                default json
// Directory    : Log storage path
//                default "./"
type LogConfig struct {
	Level       string `mapstructure:"Level"`
	Format      string `mapstructure:"Format"`
	Directory   string `mapstructure:"Directory"`
	Development bool   `mapstructure:"Development"`
}

type SuperAdminConfig struct {
	Username string `mapstructure:"Username"`
	Realname string `mapstructure:"Realname"`
	Password string `mapstructure:"Password"`
}

type AuthConfig struct {
	Enable             bool     `mapstructure:"Enable"`
	TokenExpired       int      `mapstructure:"TokenExpired"`
	IgnorePathPrefixes []string `mapstructure:"IgnorePathPrefixes"`
}

type CasbinConfig struct {
	Enable             bool     `mapstructure:"Enable"`
	Debug              bool     `mapstructure:"Debug"`
	Model              string   `mapstructure:"Model"`
	AutoLoad           bool     `mapstructure:"AutoLoad"`
	AutoLoadInternal   int      `mapstructure:"AutoLoadInternal"`
	IgnorePathPrefixes []string `mapstructure:"IgnorePathPrefixes"`
}

type DatabaseConfig struct {
	Engine      string `mapstructure:"Engine"`
	Name        string `mapstructure:"Name"`
	Host        string `mapstructure:"Host"`
	Port        int    `mapstructure:"Port"`
	Username    string `mapstructure:"Username"`
	Password    string `mapstructure:"Password"`
	TablePrefix string `mapstructure:"TablePrefix"`
	Parameters  string `mapstructure:"Parameters"`

	MaxLifetime  int `mapstructure:"MaxLifetime"`
	MaxOpenConns int `mapstructure:"MaxOpenConns"`
	MaxIdleConns int `mapstructure:"MaxIdleConns"`
}

type RedisConfig struct {
	Host      string `mapstructure:"Host"`
	Port      int    `mapstructure:"Port"`
	Password  string `mapstructure:"Password"`
	KeyPrefix string `mapstructure:"KeyPrefix"`
}

func (a *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", a.Username, a.Password, a.Host, a.Port, a.Name, a.Parameters)
}

func (a *HttpConfig) ListenAddr() string {
	if err := validator.New().Struct(a); err != nil {
		return "0.0.0.0:5100"
	}

	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (a *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
