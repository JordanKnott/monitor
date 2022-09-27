package monitorcli

import (
	"net/url"
	"strings"

	"github.com/jordanknott/monitor/internal/commands"
	"github.com/jordanknott/monitor/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	LogLevel        = "log.level"
	LogReportCaller = "log.report_caller"
	GeneralToken    = "general.token"
	GeneralUrl      = "general.url"
)

var defaults = map[string]interface{}{
	LogLevel:        "warn",
	LogReportCaller: false,
	GeneralToken:    "",
	GeneralUrl:      "http://localhost:3333",
}

var monitorCmd = &cobra.Command{
	Use:     "monitorcli",
	Long:    "a tool for monitoring driven digital wordpress installs",
	Version: commands.VersionTemplate(),
}

func initDefaults() {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

type AppConfig struct {
	Token string
	Url   string
}

func (c *AppConfig) GetURL(path string) (string, error) {
	base, err := url.Parse(c.Url)
	if err != nil {
		return "", err
	}
	ref, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(ref).String(), nil
}

func GetAppConfig() (AppConfig, error) {
	return AppConfig{
		Token: viper.GetString(GeneralToken),
		Url:   viper.GetString(GeneralUrl),
	}, nil
}

func initConfig() {
	viper.AddConfigPath("./conf")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/monitor")
	viper.SetConfigName("monitorcli")

	viper.SetEnvPrefix("MONITORCLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	initDefaults()

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
	cobra.OnInitialize(initConfig)
}

func Execute() {

	monitorCmd.AddCommand(newUserCmd(), newJobsCmd(), newInstallsCmd(), newApiCmd())
	monitorCmd.Execute()
}
