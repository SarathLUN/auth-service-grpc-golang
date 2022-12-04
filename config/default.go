package config

type Config struct {
	DBUri    string `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri string `mapstruture:"REDIS_URL"`
	Port     string `mapstructure:"PORT"`
}
