package config

import "os"

// Config carries the runtime settings for the LIMS console server.
type Config struct {
	Addr    string
	DataDir string
}

// Default returns the built-in configuration used when no environment
// override is present.
func Default() Config {
	return Config{Addr: ":8080", DataDir: "data"}
}

// LoadFromEnv reads LIMS_ADDR and LIMS_DATA_DIR, falling back to Default.
func LoadFromEnv() Config {
	cfg := Default()
	if value := os.Getenv("LIMS_ADDR"); value != "" {
		cfg.Addr = value
	}
	if value := os.Getenv("LIMS_DATA_DIR"); value != "" {
		cfg.DataDir = value
	}
	return cfg
}
