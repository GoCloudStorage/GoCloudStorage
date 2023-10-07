package opt

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type service struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type storage struct {
	Type            string `mapstructure:"type"`
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	BlockSize       int32  `mapstructure:"block_size"`
	UseSSL          bool   `mapstructure:"use_ssl"`
	BucketName      string `mapstructure:"bucket_name"`
}

type pg struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

type file struct {
	BlockSize int `mapstructure:"block_size"`
}

type redis struct {
	Addr        string `mapstructure:"addr"`
	Password    string `mapstructure:"password"`
	DB          int    `mapstructure:"db"`
	UniqueValue int    `mapstructure:"unique_value"`
}

type storageRPC struct {
	Domain    string   `mapstructure:"domain"`
	Endpoints []string `mapstructure:"endpoints"`
}

type fileRPC struct {
	Domain    string   `mapstructure:"domain"`
	Endpoints []string `mapstructure:"endpoints"`
}

type userRPC struct {
	Domain    string   `mapstructure:"domain"`
	Endpoints []string `mapstructure:"endpoints"`
}

type MQ struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type config struct {
	Mq             MQ         `mapstructure:"mq"`
	UserService    service    `mapstructure:"user_service"`
	FileService    service    `mapstructure:"file_service"`
	StorageService service    `mapstructure:"storage_service"`
	LogLevel       string     `mapstructure:"log_level"`
	Storage        storage    `mapstructure:"storage"`
	Pg             pg         `mapstructure:"pg"`
	File           file       `mapstructure:"file"`
	Redis          redis      `mapstructure:"redis"`
	StorageRPC     storageRPC `mapstructure:"storageRPC"`
	FileRPC        fileRPC    `mapstructure:"FileRPC"`
	UserRPC        userRPC    `mapstructure:"userRPC"`
}

var (
	configFile string
	Cfg        = new(config)
)

func init() {
	time.Local = time.UTC
	flag.StringVar(&configFile, "c", "etc/config.yaml", "config storage_engine path")
}

func InitConfig() {
	flag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panicf("failed to read config Cloud Storage, err: %v", err)
	}
	if err := viper.Unmarshal(Cfg); err != nil {
		logrus.Panicf("failed to unmarshal config, err: %v", err)
	}
	switch strings.ToLower(Cfg.LogLevel) {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
		logrus.SetReportCaller(true)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	}
	logrus.Infof("read config detail: %+v", Cfg)
}
