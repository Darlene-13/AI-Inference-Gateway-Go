package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func fakeEnv(vars map[string]string) func(string) (string, bool) {
	return func(name string) (string, bool) {
		v, ok := vars[name]
		return v, ok
	}
}

func loadTestdata(t *testing.T, name string, env map[string]string) (*Config, error) {
	t.Helper() // failures point at the calling test, not at this line
	raw, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read testdata: %v", err)
	}
	return Parse(raw, fakeEnv(env))
}

func TestParseValidConfigAppliesValuesAndDefaults(t *testing.T) {
	cfg, err := loadTestdata(t, "valid.yaml", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Addr != ":9090" {
		t.Errorf("Addr = %q, want :9090", cfg.Server.Addr)
	}
	if cfg.Server.ReadHeaderTimeout != 5*time.Second {
		t.Errorf("ReadHeaderTimeout default = %v, want 5s", cfg.Server.ReadHeaderTimeout)
	}
	ray := cfg.Providers["ray"]
	if ray.BaseURL != "http://localhost:8000" {
		t.Errorf("ray base_url = %q, want the ${...:-default} value", ray.BaseURL)
	}
	if ray.RequestTimeout != 90*time.Second || ray.ConnectTimeout != 2*time.Second {
		t.Errorf("ray timeouts = %v / %v, want 2s connect, 90s request", ray.ConnectTimeout, ray.RequestTimeout)
	}
	if got := cfg.Models["fast"]; len(got) != 2 || got[1].Model != "llama3.2:1b" {
		t.Errorf("fast targets = %+v", got)
	}
}

func TestParseUsesEnvironmentVariables(t *testing.T) {
	cfg, err := loadTestdata(t, "valid.yaml", map[string]string{"TEST_RAY_URL": "http://pc:8000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := cfg.Providers["ray"].BaseURL; got != "http://pc:8000" {
		t.Errorf("base_url = %q, want the environment value", got)
	}
}

func TestParseRejectsBadConfigs(t *testing.T) {
	tests := []struct {
		file string
		want string // text the error must contain
	}{
		{"missing_base_url.yaml", "providers.ray.base_url: is required"},
		{"unknown_provider.yaml", `"olama" is not defined under providers`},
		{"typo.yaml", "requst_timeout"},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			_, err := loadTestdata(t, tt.file, nil)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q\nwant it to contain %q", err, tt.want)
			}
		})
	}
}

func TestValidateReportsEveryProblemAtOnce(t *testing.T) {
	cfg := &Config{
		Providers: map[string]Provider{"ray": {BaseURL: "not a url"}},
		Models:    map[string][]Target{"fast": {{Provider: "nope"}}},
	}
	cfg.applyDefaults()

	var verr *ValidationError
	if err := cfg.Validate(); !errors.As(err, &verr) {
		t.Fatalf("want *ValidationError, got %v", err)
	}
	if len(verr.Problems) != 3 { // bad URL, unknown provider, missing model
		t.Errorf("got %d problems, want 3:\n%s", len(verr.Problems), verr)
	}
}

func TestExpandEnv(t *testing.T) {
	env := fakeEnv(map[string]string{"SET": "value", "EMPTY": ""})
	tests := []struct {
		in, want string
	}{
		{"${SET}", "value"},
		{"${SET:-fallback}", "value"},
		{"${UNSET:-fallback}", "fallback"},
		{"${EMPTY:-fallback}", "fallback"},
		{"${UNSET}", ""},
		{"${UNSET:-llama3.2:1b}", "llama3.2:1b"},           // defaults may contain colons
		{"price is $5 and $HOME", "price is $5 and $HOME"}, // other $ forms untouched
	}
	for _, tt := range tests {
		if got := expandEnv(tt.in, env); got != tt.want {
			t.Errorf("expandEnv(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
