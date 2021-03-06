package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	APIPort               string
	AdminPort             string
	APIInterface          string
	AdminInterface        string
	GracefulShutdownDelay int

	IngestionHandler string

	StreamName string

	KafkaTopic   string
	KafkaBrokers []string

	PrometheusEnabled  bool
	HealthCheckEnabled bool
	PProfEnabled       bool

	DebugMode bool

	LogFormat string

	EnableFailover  bool
	FailOverBackend string
	S3Bucket        string
	S3Timeout       int
	S3Region        string
}

var AppConfig Configuration

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if sandboxPath, ok := os.LookupEnv("MESOS_SANDBOX"); ok {
		viper.AddConfigPath(sandboxPath)
	}

	AppConfig = Configuration{}

	setDefaults()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed reading config file: %s\n", err)
	}

	viper.SetEnvPrefix("ingestion")
	viper.AutomaticEnv()

	if viper.Get("admin_port") != nil {
		viper.Set("AdminPort", viper.GetString("admin_port"))
	}
	if viper.Get("api_port") != nil {
		viper.Set("APIPort", viper.GetString("api_port"))
	}
	if viper.Get("enable_prometheus") != nil {
		viper.Set("PrometheusEnabled", viper.GetBool("enable_prometheus"))
	}
	if viper.Get("debug_mode") != nil {
		viper.Set("DebugMode", viper.GetBool("debug_mode"))
	}
	if viper.Get("api_interface") != nil {
		viper.Set("APIInterface", viper.GetString("api_interface"))
	}
	if viper.Get("admin_interface") != nil {
		viper.Set("AdminInterface", viper.GetString("admin_interface"))
	}
	if viper.Get("enable_health_check") != nil {
		viper.Set("HealthCheckEnabled", viper.GetBool("enable_health_check"))
	}
	if viper.Get("kafka_brokers") != nil {
		viper.Set("KafkaBrokers", viper.GetString("kafka_brokers"))
	}
	if viper.Get("kafka_topic") != nil {
		viper.Set("KafkaTopic", viper.GetString("kafka_topic"))
	}
	if viper.Get("graceful_shutdown") != nil {
		viper.Set("GracefulShutdownDelay", viper.GetInt("graceful_shutdown"))
	}
	if viper.Get("enable_pprof") != nil {
		viper.Set("PProfEnabled", viper.GetBool("enable_pprof"))
	}
	if viper.Get("enable_failover") != nil {
		viper.Set("EnableFailover", viper.GetBool("enable_failover"))
	}
	if viper.Get("ingestion_handler") != nil {
		viper.Set("IngestionHandler", viper.GetString("ingestion_handler"))
	}
	if viper.Get("EnableFailover") == true {
		if awsKeyID := os.Getenv("AWS_ACCESS_KEY_ID"); awsKeyID == "" {
			log.Fatalf("provide a valid AWS_ACCESS_KEY_ID")
		}
		if awsKey := os.Getenv("AWS_SECRET_KEY"); awsKey == "" {
			log.Fatalf("provide a valid AWS_SECRET_KEY")
		}
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}

func setDefaults() {
	viper.SetDefault("PProfEnabled", false)
	viper.SetDefault("PrometheusEnabled", true)
	viper.SetDefault("HealthCheckEnabled", true)
	viper.SetDefault("GracefulShutdownDelay", 10)
	viper.SetDefault("APIPort", "8080")
	viper.SetDefault("AdminPort", "8081")
	viper.SetDefault("APIInterface", "0.0.0.0")
	viper.SetDefault("AdminInterface", "0.0.0.0")
	viper.SetDefault("DebugMode", false)
	viper.SetDefault("LogFormat", "json")
	viper.SetDefault("FailOverBackend", "s3")
	viper.SetDefault("S3Timeout", 10)
	viper.SetDefault("S3Region", "us-east-1")
	viper.SetDefault("S3Bucket", "ingestion-api-messages")
	viper.SetDefault("EnableFailover", false)
	viper.SetDefault("IngestionHandler", "kafka")
}

// Get returns the config
func Get() *Configuration {
	return &AppConfig
}

// Print prints the configuration to stdout
func (c *Configuration) Print() {
	fmt.Println("Configurations:")
	fmt.Printf("IngestionHandler: %s\n", c.IngestionHandler)
	fmt.Printf("StreamName: %s\n", c.StreamName)
	fmt.Printf("KafkaTopic: %s\n", c.KafkaTopic)
	fmt.Printf("KafkaBrokers: %s\n", c.KafkaBrokers)
	fmt.Printf("PrometheusEnabled: %t\n", c.PrometheusEnabled)
	fmt.Printf("HealthCheckEnabled: %t\n", c.HealthCheckEnabled)
	fmt.Printf("PProfEnabled: %t\n", c.PProfEnabled)
	fmt.Printf("APIPort: %s\n", c.APIPort)
	fmt.Printf("AdminPort: %s\n", c.AdminPort)
	fmt.Printf("GracefulShutdownDelay: %d\n", c.GracefulShutdownDelay)
	fmt.Printf("AdminInterface: %s\n", c.AdminInterface)
	fmt.Printf("ApiInterface: %s\n", c.APIInterface)
	fmt.Printf("DebugMode: %s\n", c.DebugMode)
	fmt.Printf("FailoverBackend: %s", c.FailOverBackend)
	fmt.Printf("S3Timeout: %d", c.S3Timeout)
}
