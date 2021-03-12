package config

import "github.com/spf13/viper"

// Config - основной конфиг, содержащий в себе все остальные
type Config struct {
	DB  DB  `mapstructure:"db"`
	Bot Bot `mapstructure:"bot"`
	SSL SSL `mapstructure:"ssl"`
}

// DB - конфиг для работы с базой данных
type DB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string
	Name     string `mapstructure:"name"`
}

// Bot - конфиг для работы с телеграм ботом, нужен только для хранения токена
type Bot struct {
	APIToken string
}

// SSL - конфиг хранящий пути к сертификатам
type SSL struct {
	Fullchain string `mapstructure:"fullchain"`
	Privkey   string `mapstructure:"privkey"`
}

// Init - анмаршалим конфигурационный файл в нашу структуру
func Init() (*Config, error) {
	if err := SetUpViper(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := ParseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SetUpViper - находим файл конфигурации и пытаемся его прочитать
func SetUpViper() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

// ParseEnv - анмаршалим переменные окружения в нашу струткуру cfg
func ParseEnv(cfg *Config) error {
	if err := viper.BindEnv("API_TOKEN"); err != nil {
		return err
	}

	if err := viper.BindEnv("DB_PASSWORD"); err != nil {
		return err
	}

	cfg.Bot.APIToken = viper.GetString("API_TOKEN")

	cfg.DB.Password = viper.GetString("DB_PASSWORD")

	return nil
}
