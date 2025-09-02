package configs

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	AppName     string `env:"APP_NAME" envDefault:"Go Application Programming Interface (API)"`
	AppPort     string `env:"APP_PORT,notEmpty"`
	AppLocation string `env:"APP_LOCATION" envDefault:"Asia/Jakarta"`
	AppDebug    bool   `env:"APP_DEBUG" envDefault:"false"`
	AppVersion  string `env:"APP_VERSION" envDefault:"v1.0.0"`
	AppKey      string `env:"APP_KEY"`

	UseBodyDumpLog bool `env:"USE_BODY_DUMP_LOG" envDefault:"false"`

	DefaultTimeout int `env:"DEFAULT_TIMEOUT" envDefault:"1"`

	HeaderPrefix string `env:"HEADER_PREFIX"`

	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:","`
}

func New(toggle bool) (*Config, error) {
	var (
		tag string = "Configs.Main.New."
		cfg Config
	)

	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := godotenv.Load(); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "01",
			"error": err.Error(),
		}).Error("failed to load environment file")

		return &cfg, err
	}

	if err := env.Parse(&cfg); err != nil {
		logrus.WithFields(logrus.Fields{
			"tag":   tag + "02",
			"error": err.Error(),
		}).Error("failed to parse environment")

		return &cfg, err
	}

	if cfg.UseBodyDumpLog {
		if err := NewBodyDumpLog(); err != nil {
			logrus.WithFields(logrus.Fields{
				"tag":   tag + "03",
				"error": err.Error(),
			}).Error("failed to initiate a body dump for log")

			return &cfg, err
		}
	}

	LoadVersion(&cfg, toggle)

	return &cfg, nil
}
