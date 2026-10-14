package config

import (
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"
)

type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	return "invalid config:\n " + strings.Join(e.Problems, "\n ")
}

func (c *Config) Validate() error {
	var problems []string
	add := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}

	if c.Server.MaxBodyBytes <= 0 {
		add("server.max_body_bytes: must be positive")
	}

	if !slices.Contains([]string{"debug", "info", "warn", "error"}, c.Log.Level) {
		add("log.level: must be debug, info, warn or error, got %q", c.Log.Level)
	}

	if !slices.Contains([]string{"text", "json"}, c.Log.Format) {
		add("log.format: must be text or json, got %q", c.Log.Format)
	}

	if len(c.Providers) == 0 {
		add("providers: at least one provider is required")
	}

	for _, name := range slices.Sorted(maps.Keys(c.Providers)) {
		p := c.Providers[name]
		if p.BaseURL == "" {
			add("providers.%s.base_url: is required", name)
		} else if u, err := url.Parse(p.BaseURL); err != nil || u.Scheme == "" || u.Host == "" {
			add("providers.%s.base_url: must be an http or https URL, got %q", name, p.BaseURL)
		}

		if p.ConnectTimeout < 0 || p.RequestTimeout < 0 {
			add("providers.%s: timeouts cannot be negative", name)
		}
	}

	if len(c.Models) == 0 {
		add("models: at least one model alias is required")
	}

	for _, alias := range slices.Sorted(maps.Keys(c.Models)) {
		targets := c.Models[alias]

		if len(targets) == 0 {
			add("models.%s: needs at least one target", alias)
		}

		for i, t := range targets {
			switch {
			case t.Provider == "":
				add("models.%s[%d].provider: is required", alias, i)
			case !hasKey(c.Providers, t.Provider):
				add("models.%s[%d].provider: %q is not defined under providers", alias, i, t.Provider)
			}

			if t.Model == "" {
				add("models.%s[%d].model: is required", alias, i)
			}
		}
	}

	if len(problems) > 0 {
		return &ValidationError{Problems: problems}
	}

	return nil
}

func hasKey[V any](m map[string]V, key string) bool {
	_, ok := m[key]
	return ok
}
