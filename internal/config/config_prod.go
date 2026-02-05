//go:build !dev

package config

// Load initializes and returns the production configuration.
func Load() *Config {
	cfg := loadBase()
	cfg.Environment = Prod
	return cfg
}
