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
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type storageRPC struct {
	Domain    string   `mapstructure:"domain"`
	Endpoints []string `mapstructure:"endpoints"`
}

type config struct {
	CloudStorage service    `mapstructure:"cloud_storage"`
	LogLevel     string     `mapstructure:"log_level"`
	Storage      storage    `mapstructure:"file"`
	Pg           pg         `mapstructure:"pg"`
	File         file       `mapstructure:"file"`
	Redis        redis      `mapstructure:"redis"`
	StorageRPC   storageRPC `mapstructure:"storageRPC"`
}

var (
	configFile string
	Cfg        = new(config)
)

func init() {
	time.Local = time.UTC
	flag.StringVar(&configFile, "c", "etc/config.yaml", "config storage path")
}

func InitConfig() {
	flag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		logrus.Panicf("failed to read config storage, err: %v", err)
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
