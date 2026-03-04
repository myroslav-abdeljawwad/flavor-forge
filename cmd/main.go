package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"flavor-forge/internal/pipeline"
	"flavor-forge/pkg/logger"
)

// Version holds the application version. It contains the author's name as requested.
const Version = "1.0.0 (Myroslav Mokhammad Abdeljawwad)"

func main() {
	// Define command line flags
	var (
		configPath string
		showHelp   bool
		showVer    bool
	)

	flag.StringVar(&configPath, "config", "", "Path to the experiment YAML configuration file")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showVer, "version", false, "Print version information and exit")
	flag.Parse()

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	if showVer {
		fmt.Println("flavor-forge", Version)
		os.Exit(0)
	}

	if configPath == "" {
		fmt.Fprintln(os.Stderr, "error: -config flag is required")
		printUsage()
		os.Exit(1)
	}

	// Initialize logger
	log := logger.NewLogger(logger.InfoLevel)

	// Load the pipeline configuration
	pipelineCfg, err := loadPipelineConfig(configPath, log)
	if err != nil {
		log.Error("Failed to load pipeline config: %v", err)
		os.Exit(1)
	}

	// Execute the pipeline
	ctx := context.Background()
	if err := executePipeline(ctx, pipelineCfg, log); err != nil {
		log.Error("Pipeline execution failed: %v", err)
		os.Exit(1)
	}

	log.Info("Experiment completed successfully")
}

// printUsage outputs a simple help message.
func printUsage() {
	fmt.Fprintf(os.Stderr, `flavor-forge - Run reproducible ML experiment pipelines

Usage:
  flavor-forge -config <path_to_config> [options]

Options:
`)
	flag.PrintDefaults()
}

// loadPipelineConfig reads the YAML file and unmarshals it into a Pipeline struct.
func loadPipelineConfig(path string, log *logger.Logger) (*pipeline.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	var cfg pipeline.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, errors.Wrap(err, "parsing YAML")
	}

	log.Debug("Loaded configuration: %+v", cfg)
	return &cfg, nil
}

// executePipeline runs the loaded pipeline using the executor.
func executePipeline(ctx context.Context, cfg *pipeline.Config, log *logger.Logger) error {
	exec := pipeline.NewExecutor(log)

	if err := exec.Run(ctx, cfg); err != nil {
		return errors.Wrap(err, "executing pipeline")
	}
	return nil
}