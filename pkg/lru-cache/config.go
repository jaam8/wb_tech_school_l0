package lrucache

type Config struct {
	CleanupInterval int `env:"CLEANUP_INTERVAL" env-default:"5"    yaml:"cleanup_interval"`
	TTL             int `env:"TTL"              env-default:"15"   yaml:"ttl"`
	Capacity        int `env:"CAPACITY"         env-default:"1000" yaml:"capacity"`
}
