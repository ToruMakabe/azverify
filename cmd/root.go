package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/logutils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	Tenant       Tenant       `mapstructure:"tenant"`
	Subscription Subscription `mapstructure:"subscription"`
	Environment  string       `mapstructure:"environment"`
	Client       AADClient    `mapstructure:"client"`
	Cert         AADCert      `mapstructure:"certificate"`
	Log          LogConfig    `mapstructure:"log"`
}

type Tenant struct {
	ID string `mapstructure:"id"`
}

type Subscription struct {
	ID string `mapstructure:"id"`
}

type AADClient struct {
	ID     string `mapstructure:"id"`
	Secret string `mapstructure:"secret"`
}

type AADCert struct {
	Path     string `mapstructure:"path"`
	Password string `mapstructure:"password"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

var (
	configFile string
	config     Config
	envPrefix  string
)

const (
	defaultAzureEnv = "public"
	defaultLogLevel = "INFO"
)

var rootCmd = &cobra.Command{
	Use:   "azverify",
	Short: "Azure Resource Verifier with Resource Graph",
	Long: `Azure Resource Verifier with Resource Graph.

You can verify if there is a difference between your desired properties
and actual with this CLI. This CLI read your desired properties as JSON files,
and query to Azure Resource Graph API, then check the difference.`,
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file path (default \"$HOME/.azverify/config.toml\")")
	rootCmd.PersistentFlags().StringVar(&envPrefix, "env-prefix", "AZV", "env prefix")

	rootCmd.PersistentFlags().String("tenant-id", "", "Azure AD tenant ID")
	rootCmd.PersistentFlags().String("subscription-id", "", "Azure subscription ID")
	rootCmd.PersistentFlags().String("environment", "", "Azure environment ([public]/usgovernment/german/china)")
	rootCmd.PersistentFlags().String("client-id", "", "Azure AD service principal App ID")
	rootCmd.PersistentFlags().String("client-secret", "", "Azure AD service principal App secret")
	rootCmd.PersistentFlags().String("cert-path", "", "PKCS12 (.pfx) cert file path")
	rootCmd.PersistentFlags().String("cert-password", "", "cert file password")
	rootCmd.PersistentFlags().String("log-level", "", "log level (DEBUG/[INFO]/ERROR)")

	viper.BindPFlag("tenant.id", rootCmd.PersistentFlags().Lookup("tenant-id"))
	viper.BindPFlag("subscription.id", rootCmd.PersistentFlags().Lookup("subscription-id"))
	viper.BindPFlag("environment", rootCmd.PersistentFlags().Lookup("environment"))
	viper.BindPFlag("client.id", rootCmd.PersistentFlags().Lookup("client-id"))
	viper.BindPFlag("client.secret", rootCmd.PersistentFlags().Lookup("client-secret"))
	viper.BindPFlag("certificate.path", rootCmd.PersistentFlags().Lookup("cert-path"))
	viper.BindPFlag("certificate.password", rootCmd.PersistentFlags().Lookup("cert-password"))
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))

	viper.SetDefault("environment", defaultAzureEnv)
	viper.SetDefault("log.level", defaultLogLevel)
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		checkErr(err)

		cfgDir := home + "/.azverify"
		viper.AddConfigPath(cfgDir)
		viper.SetConfigName("config")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	err := viper.Unmarshal(&config)
	checkErr(err)

	var ll string
	switch config.Log.Level {
	case "DEBUG":
		ll = "DEBUG"
	case "INFO":
		ll = "INFO"
	case "ERROR":
		ll = "ERROR"
	default:
		fmt.Printf("%s is invalid log level. set it INFO\n", config.Log.Level)
		ll = "INFO"
	}
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "ERROR"},
		MinLevel: logutils.LogLevel(ll),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
