package config

import "time"

type Config struct {
	Server    Server              `yaml:"server"`
	Log       Log                 `yaml:"log"`
	Providers map[string]Provider `yaml:"providers"`
	Models    map[string][]Target `yaml:"models"`
}

type Server struct {
	//Address or port to listen to
	Addr              string        `yaml:"addr"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
	DrainDelay        time.Duration `yaml:"drain_delay"`
	MaxBodyBytes      int64         `yaml:"max_body_bytes"`
}

type Log struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Provider struct {
	BaseURL        string        `yaml:"base_url"`
	APIKey         string        `yaml:"api_key"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
}

type Target struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
}
