package config_manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

//
// ──────────────────────────────────────────────────────────────
//   SIMULATION CONFIG
// ──────────────────────────────────────────────────────────────
//

// SimulationConfig holds engine-level simulation settings.
type SimulationConfig struct {
	BatchPath             string   `yaml:"batch_path"`
	IterationLimit        int      `yaml:"iteration_limit"`
	DefaultAttackerLimit  int      `yaml:"default_attacker_limit"`
	StoreTopAttackerLimit int      `yaml:"store_top_attacker_limit"`
	SyncRepetitions       int      `yaml:"sync_repetitions"`
	AttackModes           []string `yaml:"attack_modes"`
	// DatabaseName optionally overrides DB_NAME from .env, so an experiment
	// can run against its own database. Empty means fall back to .env.
	DatabaseName string `yaml:"database_name"`
}

// DefaultSimulationConfig returns defaults equivalent
// to the former hardcoded constants.
func DefaultSimulationConfig() SimulationConfig {
	return SimulationConfig{
		BatchPath:             filepath.Join("config", "batches"),
		IterationLimit:        100000,
		DefaultAttackerLimit:  100,
		StoreTopAttackerLimit: 10,
		SyncRepetitions:       100,
		AttackModes:           []string{"NAIVE", "GEOMETRIC", "MAJORITY"},
	}
}

// LoadSimulationConfig loads simulation.yaml and merges it with defaults.
func LoadSimulationConfig(path string) (SimulationConfig, error) {
	cfg := DefaultSimulationConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read simulation config: %w", err)
	}

	// Partial config loaded from YAML
	var partial SimulationConfig
	if err := yaml.Unmarshal(data, &partial); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal simulation config: %w", err)
	}

	// Merge partial into defaults
	cfg = mergeSimulationConfig(cfg, partial)
	return cfg, nil
}

func mergeSimulationConfig(base, partial SimulationConfig) SimulationConfig {
	if strings.TrimSpace(partial.BatchPath) != "" {
		base.BatchPath = partial.BatchPath
	}
	if partial.IterationLimit != 0 {
		base.IterationLimit = partial.IterationLimit
	}
	if partial.DefaultAttackerLimit != 0 {
		base.DefaultAttackerLimit = partial.DefaultAttackerLimit
	}
	if partial.StoreTopAttackerLimit != 0 {
		base.StoreTopAttackerLimit = partial.StoreTopAttackerLimit
	}
	if partial.SyncRepetitions != 0 {
		base.SyncRepetitions = partial.SyncRepetitions
	}
	if partial.AttackModes != nil {
		base.AttackModes = partial.AttackModes
	}
	if strings.TrimSpace(partial.DatabaseName) != "" {
		base.DatabaseName = partial.DatabaseName
	}
	return base
}

//
// ──────────────────────────────────────────────────────────────
//   TRACKING CONFIG
// ──────────────────────────────────────────────────────────────
//

// TrackingConfig holds tracking-related parameters.
type TrackingConfig struct {
	SSETTL         string `yaml:"sse_ttl"`         // raw YAML string
	SkipIterations int    `yaml:"skip_iterations"` // integer in YAML

	ParsedTTL time.Duration `yaml:"-"` // derived after parsing
}

// DefaultTrackingConfig provides defaults equivalent to your constants.
func DefaultTrackingConfig() TrackingConfig {
	return TrackingConfig{
		SSETTL:         "1m",
		SkipIterations: 150,
	}
}

// LoadTrackingConfig loads tracking.yaml and merges it with defaults.
func LoadTrackingConfig(path string) (TrackingConfig, error) {
	cfg := DefaultTrackingConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read tracking config: %w", err)
	}

	var partial TrackingConfig
	if err := yaml.Unmarshal(data, &partial); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal tracking config: %w", err)
	}

	// merge
	cfg = mergeTrackingConfig(cfg, partial)

	// parse duration
	ttl, err := time.ParseDuration(cfg.SSETTL)
	if err != nil {
		return cfg, fmt.Errorf("invalid tracking.sse_ttl duration: %w", err)
	}
	cfg.ParsedTTL = ttl

	return cfg, nil
}

func mergeTrackingConfig(base, partial TrackingConfig) TrackingConfig {
	if partial.SSETTL != "" {
		base.SSETTL = partial.SSETTL
	}
	if partial.SkipIterations != 0 {
		base.SkipIterations = partial.SkipIterations
	}
	return base
}
