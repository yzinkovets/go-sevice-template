package config

type AppConfig struct {
	AppName string `env:"APP_NAME"`
	LogLvl  string `env:"APP_LOG_LEVEL"`
}

type DbConfig struct {
	Host             string `env:"DB_HOST,required"`
	Port             int    `env:"DB_PORT" envDefault:"5432"`
	Db               string `env:"DB_NAME,required"`
	User             string `env:"DB_USER,required"`
	Password         string `env:"DB_PASSWORD,required"`
	MaxOpenConns     int    `env:"DB_MAX_OPEN_CONNS" envDefault:"10"` // max open connections. 0 means unlimited
	MaxIdleConns     int    `env:"DB_MAX_IDLE_CONNS" envDefault:"10"` // max idle connections. 0 means unlimited
	InsertTimeoutSec int    `env:"DB_INSERT_TIMEOUT_SEC" envDefault:"60"`
}

type KafkaConfig struct {
	Host            string `env:"KAFKA_HOST"`
	ClientId        string `env:"KAFKA_CLIENT_ID"`
	Topic           string `env:"KAFKA_TOPIC"`
	AutoOffsetReset string `env:"KAFKA_AUTO_OFFSET_RESET" envDefault:"latest"`
}

type MqttConfig struct {
	Host                string `env:"MQTT_HOST,required"`
	Username            string `env:"MQTT_USERNAME"`
	Password            string `env:"MQTT_PASSWORD"`
	CertPath            string `env:"MQTT_CERT_PATH"`
	ClientId            string `env:"MQTT_CLIENT_ID" envDefault:"go_iot_network_topology"`
	TopicSend           string `env:"MQTT_TOPIC_SEND"`
	TopicReceive        string `env:"MQTT_TOPIC_RECEIVE"`
	Qos                 byte   `env:"MQTT_QOS" envDefault:"0"`
	KeepAliveTimeoutSec int    `env:"MQTT_KEEP_ALIVE_TIMEOUT_SEC" envDefault:"60"`
	PingTimeoutSec      int    `env:"MQTT_PING_TIMEOUT_SEC" envDefault:"5"`
}

type JwtAuthConfig struct {
	JwksFilePath          string `env:"APP_JWKS_FILE,required"`
	JwtCheckUrl           string `env:"APP_JWT_CHECK_URL,required"`
	JwtCheckUrlTimeoutSec int    `env:"APP_JWT_CHECK_URL_TIMEOUT_SEC" envDefault:"3"`
}

type ServerConfig struct {
	Addr string `env:"APP_WEB_ADDR,required"`
	Tls  struct {
		IsEnable    bool   `env:"APP_WEB_TLS_IS_ENABLE"`
		TlsCertFile string `env:"APP_WEB_TLS_CERT_FILE"`
		TlsKeyFile  string `env:"APP_WEB_TLS_KEY_FILE"`
	}
	JwtAuthConfig
}

type MainConfig struct {
	AppConfig
	DbConfig
	KafkaConfig
	MqttConfig
	ServerConfig
}
