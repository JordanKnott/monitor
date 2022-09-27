package monitor

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jordanknott/monitor/internal/commands"
	"github.com/jordanknott/monitor/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	LogLevel         = "log.level"
	LogReportCaller  = "log.report_caller"
	ServerHostname   = "server.hostname"
	DatabaseHost     = "database.host"
	DatabaseName     = "database.name"
	DatabaseUser     = "database.user"
	DatabasePassword = "database.password"
	DatabasePort     = "database.port"
	DatabaseSslMode  = "database.ssl_mode"

	SecurityTokenExpiration = "security.token_expiration"
	SecuritySecret          = "security.secret"

	JobEnabled   = "job.enabled"
	JobBroker    = "job.broker"
	JobStore     = "job.store"
	JobQueueName = "job.queue_name"

	MessageQueue = "message.queue"

	SmtpFrom       = "smtp.from"
	SmtpHost       = "smtp.host"
	SmtpPort       = "smtp.port"
	SmtpUsername   = "smtp.username"
	SmtpPassword   = "smtp.password"
	SmtpSkipVerify = "false"
)

var defaults = map[string]interface{}{
	LogLevel:                "warn",
	LogReportCaller:         false,
	ServerHostname:          "0.0.0.0:3333",
	DatabaseHost:            "127.0.0.1",
	DatabaseName:            "monitor",
	DatabaseUser:            "monitor",
	DatabasePassword:        "monitor",
	DatabasePort:            "8865",
	DatabaseSslMode:         "disable",
	SecurityTokenExpiration: "15m",
	SecuritySecret:          "",
	MessageQueue:            "localhost:6379",
	JobEnabled:              false,
	JobBroker:               "redis://localhost:6379",
	JobStore:                "redis://localhost:6379",
	JobQueueName:            "taskcafe_tasks",
	SmtpFrom:                "no-reply@example.com",
	SmtpHost:                "localhost",
	SmtpPort:                "587",
	SmtpUsername:            "",
	SmtpPassword:            "",
	SmtpSkipVerify:          false,
}

func getJobConfig() config.JobConfig {
	return config.JobConfig{
		Enabled:   viper.GetBool(JobEnabled),
		Broker:    viper.GetString(JobBroker),
		QueueName: viper.GetString(JobQueueName),
		Store:     viper.GetString(JobStore),
	}
}
func GetAppConfig() (config.AppConfig, error) {
	secret := viper.GetString(SecuritySecret)
	if strings.TrimSpace(secret) == "" {
		logrus.Warn("server.secret is not set, generating a random secret")
		secret = uuid.New().String()
	}
	securityCfg, err := getSecurityConfig(viper.GetString(SecurityTokenExpiration), []byte(secret))
	if err != nil {
		return config.AppConfig{}, err
	}
	jobCfg := getJobConfig()
	databaseCfg := getDatabaseConfig()
	emailCfg := getEmailConfig()
	messageCfg := config.MessageQueueConfig{URI: viper.GetString("message.queue")}
	return config.AppConfig{
		Email:        emailCfg,
		Security:     securityCfg,
		Database:     databaseCfg,
		Job:          jobCfg,
		MessageQueue: messageCfg,
	}, err
}

func getEmailConfig() config.EmailConfig {
	return config.EmailConfig{
		From:               viper.GetString(SmtpFrom),
		Host:               viper.GetString(SmtpHost),
		Port:               viper.GetInt(SmtpPort),
		Username:           viper.GetString(SmtpUsername),
		Password:           viper.GetString(SmtpPassword),
		InsecureSkipVerify: viper.GetBool(SmtpSkipVerify),
	}
}

func getSecurityConfig(accessTokenExp string, secret []byte) (config.SecurityConfig, error) {
	exp, err := time.ParseDuration(accessTokenExp)
	if err != nil {
		logrus.WithError(err).Error("issue parsing duration")
		return config.SecurityConfig{}, err
	}
	return config.SecurityConfig{AccessTokenExpiration: exp, Secret: secret}, nil
}

func getDatabaseConfig() config.DatabaseConfig {
	return config.DatabaseConfig{
		Username: viper.GetString(DatabaseUser),
		Password: viper.GetString(DatabasePassword),
		Port:     strconv.Itoa(viper.GetInt(DatabasePort)),
		SslMode:  viper.GetString(DatabaseSslMode),
		Name:     viper.GetString(DatabaseName),
		Host:     viper.GetString(DatabaseHost),
	}
}

func InitDefaults() {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

var monitorCmd = &cobra.Command{
	Use:     "monitor",
	Long:    "a tool for monitoring driven digital wordpress installs",
	Version: commands.VersionTemplate(),
}

var migration http.FileSystem

func InitConfig() {
	viper.AddConfigPath("./conf")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/monitor")
	viper.SetConfigName("monitor")

	viper.SetEnvPrefix("MONITOR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	InitDefaults()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
	level := config.GetLogLevel(viper.GetString(LogLevel))
	logrus.SetReportCaller(viper.GetBool(LogReportCaller))
	logrus.SetLevel(level)
}

func init() {
	cobra.OnInitialize(InitConfig)
}

func Execute() {
	monitorCmd.AddCommand(newWorkerCmd(), newMigrateCmd(), newWebCmd(), newTokenCmd())
	monitorCmd.Execute()
}
