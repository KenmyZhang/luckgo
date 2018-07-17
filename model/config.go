package model

import (
	"luckgo/log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	ReleaseMode     bool
	ServiceSettings ServiceSettings
	LogSettings     LogSettings
	SentrySettings  SentrySettings
	SqlSettings     SqlSettings
}

type ServiceSettings struct {
	ListenAddress string
	ReadTimeout   int
	WriteTimeout  int
}

type SqlSettings struct {
	DriverName               *string
	DataSource               *string
	DataSourceReplicas       []string
	DataSourceSearchReplicas []string
	MaxIdleConns             *int
	MaxOpenConns             *int
	Trace                    *bool
	QueryTimeout             *int
}

type LogSettings struct {
	EnableConsole bool
	ConsoleLevel  string
	ConsoleJson   *bool
	EnableFile    bool
	FileLevel     string
	FileJson      *bool
	FileLocation  string
}

type SentrySettings struct {
	SentryDsn string
	SentryLog string
}

const (
	SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS = ":8065"
	DATABASE_DRIVER_MYSQL                       = "mysql"
	SQL_SETTINGS_DEFAULT_DATA_SOURCE            = "mmuser:mostest@tcp(dockerhost:3306)/mattermost_test?charset=utf8mb4,utf8&readTimeout=30s&writeTimeout=30s"
)

var Cfg *Config = &Config{}

func (o *Config) SetDefaults() {
	o.LogSettings.SetDefaults()
	o.ServiceSettings.SetDefaults()
	o.SqlSettings.SetDefaults()
}

func (s *SqlSettings) SetDefaults() {
	if s.DriverName == nil {
		s.DriverName = NewString(DATABASE_DRIVER_MYSQL)
	}

	if s.DataSource == nil {
		s.DataSource = NewString(SQL_SETTINGS_DEFAULT_DATA_SOURCE)
	}

	if s.MaxIdleConns == nil {
		s.MaxIdleConns = NewInt(20)
	}

	if s.MaxOpenConns == nil {
		s.MaxOpenConns = NewInt(300)
	}

	if s.QueryTimeout == nil {
		s.QueryTimeout = NewInt(30)
	}
}

func (s *ServiceSettings) SetDefaults() {
	if s.ListenAddress == "" {
		s.ListenAddress = SERVICE_SETTINGS_DEFAULT_LISTEN_AND_ADDRESS
	}
}

func (s *LogSettings) SetDefaults() {
	if s.ConsoleJson == nil {
		s.ConsoleJson = NewBool(true)
	}

	if s.FileJson == nil {
		s.FileJson = NewBool(true)
	}
}

func NewBool(b bool) *bool       { return &b }
func NewInt(n int) *int          { return &n }
func NewString(s string) *string { return &s }

func (s *LogSettings) PrintLogConfig() {
	log.Info(fmt.Sprintf("EnableConsole:%v", s.EnableConsole))
	log.Info(fmt.Sprintf("ConsoleLevel:%v", s.ConsoleLevel))
	log.Info(fmt.Sprintf("ConsoleJson:%v", *s.ConsoleJson))
	log.Info(fmt.Sprintf("EnableFile:%v", s.EnableFile))
	log.Info(fmt.Sprintf("FileLevel:%v", s.FileLevel))
	log.Info(fmt.Sprintf("FileJson:%v", *s.FileJson))
	log.Info(fmt.Sprintf("FileLocation:%v", s.FileLocation))
}

func (o *Config) PrintConfig() {
	o.LogSettings.PrintLogConfig()
	log.Info(fmt.Sprintf("ReleaseMode:%v", o.ReleaseMode))
}

func (o *Config) LoggerConfigFromLoggerConfig() *log.LoggerConfiguration {
	return &log.LoggerConfiguration{
		EnableConsole: o.LogSettings.EnableConsole,
		ConsoleJson:   *o.LogSettings.ConsoleJson,
		ConsoleLevel:  strings.ToLower(o.LogSettings.ConsoleLevel),
		EnableFile:    o.LogSettings.EnableFile,
		FileJson:      *o.LogSettings.FileJson,
		FileLevel:     strings.ToLower(o.LogSettings.FileLevel),
		FileLocation:  GetLogFileLocation(o.LogSettings.FileLocation),
	}
}

func GetConfig() *Config {
	return Cfg
}

func LoadConfig(filename string) {
	if filename == "" {
		filename = "config.json"
	}

	filename = FindConfigFile(filename)
	log.Info("config file:" + filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(data, Cfg)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	Cfg.PrintConfig()
	return
}
