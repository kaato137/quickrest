package conf

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultContentType    = "application/json"
	defaultStatusCode     = http.StatusOK
	defaultReloadInterval = 2 * time.Second
	defaultRecordDir      = "records"
)

var defaultPaths = [...]string{
	"quickrest.yml",
	"quickrest.yaml",
}

var wildcardRegexp = regexp.MustCompile(`\{([a-zA-Z][a-zA-Z0-9_]*)\}`)

type Config struct {
	Address        string        `yaml:"addr"`
	ReloadInterval time.Duration `yaml:"reload_interval"`
	RecordDir      string        `yaml:"record_dir"`

	Routes []RouteConfig `yaml:"routes"`

	Path string
}

type RouteConfig struct {
	Path        string            `yaml:"path"`
	Body        string            `yaml:"body"`
	BodyJS      string            `yaml:"body_js"`
	ContentType string            `yaml:"content_type"`
	Headers     map[string]string `yaml:"headers"`
	StatusCode  int               `yaml:"status"`
	Record      bool              `yaml:"record"`
	Latency     time.Duration     `yaml:"latency"`
	Jitter      time.Duration     `yaml:"jitter"`

	Wildcards []string
}

func LoadConfigFromFile(path string) (*Config, error) {
	if path == "" {
		defaultPath, err := findDefaultPaths()
		if err != nil {
			return nil, err
		}
		path = defaultPath
	}

	cfgFile, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigFileNotExist
		}
		return nil, fmt.Errorf("read file: %w", err)
	}

	var cfg Config
	if err := yaml.NewDecoder(cfgFile).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	enrichConfig(&cfg, path)

	return &cfg, nil
}

func enrichConfig(cfg *Config, path string) {
	cfg.Path = path

	setDefaults(cfg)

	for i := range cfg.Routes {
		resolvePlaceholders(&cfg.Routes[i])
		setRouteDefaults(&cfg.Routes[i])
	}
}

func setDefaults(cfg *Config) {
	if cfg.ReloadInterval == 0 {
		cfg.ReloadInterval = defaultReloadInterval
	}

	if cfg.RecordDir == "" {
		cfg.RecordDir = defaultRecordDir
	}
}

func setRouteDefaults(r *RouteConfig) {
	if r.ContentType == "" {
		r.ContentType = defaultContentType
	}

	if r.StatusCode == 0 {
		r.StatusCode = defaultStatusCode
	}
}

func resolvePlaceholders(r *RouteConfig) {
	results := wildcardRegexp.FindAllStringSubmatch(r.Path, -1)

	if len(results) == 0 {
		return
	}

	submatches := make([]string, 0, len(results))
	for i := range results {
		submatches = append(submatches, results[i][1])
	}

	r.Wildcards = submatches
}
