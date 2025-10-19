package lru_cache

type Config struct {
	CleanupInterval int `yaml:"cleanup_interval" env:"CLEANUP_INTERVAL" env-default:"5"`
	TTL             int `yaml:"ttl" env:"TTL" env-default:"15"`
	Capacity        int `yaml:"capacity" env:"CAPACITY" env-default:"1000"`
}
