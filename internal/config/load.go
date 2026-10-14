package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"time"

	"go.yaml.in/yaml/v3"
)

// Load reads, expands and validates the config at file path
// Create a function called load that takes filepath as a string a returns a pointer to config or an error
func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg, err := Parse(raw, os.LookupEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

// Parse loads the file on bytes that is already in memory. The env lookup is passed in so tests can actually use a fake env
func Parse(raw []byte, lookupEnv func(string) (string, bool)) (*Config, error) {
	expanded := expandEnv(string(raw), lookupEnv)

	dec := yaml.NewDecoder(bytes.NewReader([]byte(expanded)))

	// Known files turns unknown keys into errors, so a type like request timeout fails loudly
	dec.KnownFields(true)

	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("file is empty")
		}
		return nil, fmt.Errorf("parse YAML: %w", err)
	}

	cfg.applyDefaults()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil

}

var envPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(?::-([^}]*))?\}`)

// Ensures that the env pattern matches the ${NAME} and
func expandEnv(s string, lookupEnv func(string) (string, bool)) string {
	return envPattern.ReplaceAllStringFunc(s, func(match string) string {
		parts := envPattern.FindStringSubmatch(match)
		if value, ok := lookupEnv(parts[1]); ok && value != "" {
			return value
		}

		return parts[2]
	})
}

func (c *Config) applyDefaults() {
	setIfZero(&c.Server.Addr, ":8080")
	setIfZero(&c.Server.ReadHeaderTimeout, 5*time.Second)
	setIfZero(&c.Server.ReadTimeout, 30*time.Second)
	setIfZero(&c.Server.IdleTimeout, 120*time.Second)
	setIfZero(&c.Server.ShutdownTimeout, 30*time.Second)
	setIfZero(&c.Server.MaxBodyBytes, 1<<20) // 1MiB
	setIfZero(&c.Log.Level, "info")
	setIfZero(&c.Log.Format, "text")

	// Map values are not addressable, so copy each provider out and write it back
	for name, p := range c.Providers {
		setIfZero(&p.ConnectTimeout, 2*time.Second)
		setIfZero(&p.RequestTimeout, 60*time.Second)
		c.Providers[name] = p
	}
}

func setIfZero[T comparable](field *T, value T) {
	var zero T
	if *field == zero {
		*field = value
	}
}
