package config

import (
	"github.com/spf13/viper"

	"order-service/utils/helper"

	"log"
	"os"
)

var Config AppConfig

type AppConfig struct {
	Port                               int             `json:"port" yaml:"port"`
	AppName                            string          `json:"appName" yaml:"appName"`
	AppEnv                             string          `json:"appEnv" yaml:"appEnv"`
	AppDebug                           bool            `json:"appDebug" yaml:"appDebug"`
	SignatureKey                       string          `json:"signatureKey" yaml:"signatureKey"`
	Database                           Database        `json:"database" yaml:"database"`
	InternalService                    InternalService `json:"internalService" yaml:"internalService"`
	KafkaHosts                         []string        `json:"kafkaHosts" yaml:"kafkaHosts"`
	KafkaTimeoutInMs                   int             `json:"kafkaTimeoutInMs" yaml:"kafkaTimeoutInMs"`
	KafkaMaxRetry                      int             `json:"kafkaMaxRetry" yaml:"kafkaMaxRetry"`
	KafkaProducerTopic                 string          `json:"kafkaProducerTopic" yaml:"kafkaProducerTopic"`
	KafkaConsumerFetchDefault          int32           `json:"kafkaConsumerFetchDefault" yaml:"kafkaConsumerFetchDefault"`
	KafkaConsumerFetchMin              int32           `json:"kafkaConsumerFetchMin" yaml:"kafkaConsumerFetchMin"`
	KafkaConsumerFetchMax              int32           `json:"kafkaConsumerFetchMax" yaml:"kafkaConsumerFetchMax"`
	KafkaConsumerMaxWaitTimeInMs       int32           `json:"kafkaConsumerMaxWaitTimeInMs" yaml:"kafkaConsumerMaxWaitTimeInMs"`     //nolint:lll
	KafkaConsumerMaxProcessingTimeInMs int32           `json:"kafkaConsumerMaxProcessingTimeInMs" yaml:"kafkaConsumerMaxProcTimeMs"` //nolint:lll
	KafkaConsumerBackoffTimeInMs       int32           `json:"kafkaConsumerBackoffTimeInMs" yaml:"kafkaConsumerBackoffTimeMs"`       //nolint:lll
	KafkaConsumerTopics                []string        `json:"kafkaConsumerStatusTopics" yaml:"kafkaConsumerTopics"`
	KafkaConsumerGroupID               string          `json:"kafkaConsumerGroupID" yaml:"kafkaConsumerGroupID"`
	SentryDsn                          string          `json:"sentryDsn" yaml:"sentryDsn"`
	SentrySampleRate                   float64         `json:"sentrySampleRate" yaml:"sentrySampleRate"`
	SentryEnableTracing                bool            `json:"SentryEnableTracing" yaml:"SentryEnableTracing"`
	CircuitBreakerMaxRequest           uint32          `json:"circuitBreakerMaxRequest" yaml:"circuitBreakerMaxRequest"`
	CircuitBreakerTimeoutInSecond      uint32          `json:"circuitBreakerTimeoutInSecond" yaml:"circuitBreakerTimeoutInSecond"` //nolint:lll
	RateLimiterMaxRequest              float64         `json:"rateLimiterMaxRequest" yaml:"rateLimiterMaxRequest"`
	RateLimiterTimeSecond              int             `json:"rateLimiterTimeSecond" yaml:"rateLimiterTimeSecond"`
}

type Database struct {
	Host                  string `json:"host" yaml:"host"`
	Port                  int    `json:"port" yaml:"port"`
	Name                  string `json:"name" yaml:"name"`
	Username              string `json:"username" yaml:"username"`
	Password              string `json:"password" yaml:"password"`
	MaxOpenConnection     int    `json:"maxOpenConnection" yaml:"maxOpenConnection"`
	MaxLifetimeConnection int    `json:"maxLifetimeConnection" yaml:"maxLifetimeConnection"`
	MaxIdleConnection     int    `json:"maxIdleConnection" yaml:"maxIdleConnection"`
	MaxIdleTime           int    `json:"maxIdleTime" yaml:"maxIdleTime"`
	AutoMigrate           bool   `json:"autoMigrate" yaml:"autoMigrate"`
}

type InternalService struct {
	Payment Payment `json:"payment" yaml:"payment"`
	RBAC    RBAC    `json:"rbac" yaml:"rbac"`
	Package Package `json:"package" yaml:"package"`
	Invoice Invoice `json:"invoice" yaml:"invoice"`
}

type Invoice struct {
	Host      string `json:"host" yaml:"host"`
	SecretKey string `json:"secret_key" yaml:"secretKey"`
}

type Payment struct {
	Host      string `json:"host" yaml:"host"`
	SecretKey string `json:"secret_key" yaml:"secretKey"`
}

type RBAC struct {
	Host      string `json:"host" yaml:"host"`
	SecretKey string `json:"secret_key" yaml:"secretKey"`
}

type Package struct {
	Host      string `json:"host" yaml:"host"`
	SecretKey string `json:"secret_key" yaml:"secretKey"`
}

func Init() {
	err := helper.BindFromJSON(&Config, "config.json", ".")
	if err != nil {
		log.Printf("failed load cold config from file: %s", viper.ConfigFileUsed())
		err = helper.BindFromConsul(&Config, os.Getenv("CONSUL_HTTP_URL"), os.Getenv("CONSUL_HTTP_KEY"))
		if err != nil {
			panic(err)
		}
	}
}
