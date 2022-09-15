package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// SampleConfig is the sample configuration.
type SampleConfig struct {
	// Name is the sample name.
	Name string `yaml:"name"`

	// Description is the description of the sample.
	Description string `yaml:"description"`

	// Tags is the tags of the sample.
	Tags map[string]string `yaml:"tags"`

	// ScrapeInterval is the interval to scrape the sample.
	Interval time.Duration `yaml:"interval"`

	// Timeout is the timeout to scrape the sample.
	Timeout time.Duration `yaml:"timeout"`

	// Steps is the steps to scrape the sample.
	Steps []StepConfig `yaml:"steps"`

	// Variables is the variables to scrape the sample
	Variables []map[string]string `yaml:"variables"`
}

// StepConfig is the step configuration.
type StepConfig struct {
	// Plugin is the plugin to scrape the sample. If not value given, the latest used plugin will be used.
	Plugin string `yaml:"plugin"`
	// Action is the action to scrape the sample
	Action string `yaml:"action"`
	// Parameters is the parameters to scrape the sample
	Parameters map[string]string `yaml:"parameters"`
	// Negate is the negate to scrape the sample
	Negate bool `yaml:"negate"`
}

// LoadSampleConfig loads from byte array.
func LoadSampleConfig(data []byte) (*SampleConfig, error) {
	var config SampleConfig
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	if config.Variables == nil {
		config.Variables = make([]map[string]string, 0)
	}

	if len(config.Variables) == 0 {
		config.Variables = append(config.Variables, make(map[string]string))
	}

	return &config, nil
}

// LoadSampleConfigFromFile loads from file.
func LoadSampleConfigFromFile(path string) (*SampleConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadSampleConfig(data)
}
