package pipeline_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flavor-forge/internal/pipeline"
)

// TestLoaderSuccess verifies that a valid YAML configuration is parsed into
// a Pipeline struct with the expected number of steps and metadata.
func TestLoaderSuccess(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "pipeline-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Copy the example experiment.yml into the temp directory
	srcPath := filepath.Join("examples", "experiment.yml")
	dstPath := filepath.Join(tmpDir, "experiment.yml")
	data, err := ioutil.ReadFile(srcPath)
	require.NoError(t, err)
	err = ioutil.WriteFile(dstPath, data, 0644)
	require.NoError(t, err)

	pl, err := pipeline.LoadPipeline(dstPath)
	require.NoError(t, err)
	assert.NotNil(t, pl)
	assert.GreaterOrEqual(t, len(pl.Steps), 1, "pipeline should contain at least one step")
}

// TestLoaderInvalidYAML ensures that malformed YAML triggers an error.
func TestLoaderInvalidYAML(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "pipeline-invalid-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	invalidPath := filepath.Join(tmpDir, "bad.yml")
	err = ioutil.WriteFile(invalidPath, []byte(":: bad yaml ::"), 0644)
	require.NoError(t, err)

	_, err = pipeline.LoadPipeline(invalidPath)
	assert.Error(t, err, "expected error for malformed YAML")
}

// TestExecutorRunsSteps verifies that the executor runs all steps and
// captures their outputs in order. The example experiment uses simple shell
// commands that produce deterministic output.
func TestExecutorRunsSteps(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "pipeline-exec-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	srcPath := filepath.Join("examples", "experiment.yml")
	dstPath := filepath.Join(tmpDir, "experiment.yml")
	data, err := ioutil.ReadFile(srcPath)
	require.NoError(t, err)
	err = ioutil.WriteFile(dstPath, data, 0644)
	require.NoError(t, err)

	pl, err := pipeline.LoadPipeline(dstPath)
	require.NoError(t, err)

	// Capture executor output
	var outBuf bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = pipeline.Execute(ctx, pl, &outBuf)
	require.NoError(t, err)

	output := outBuf.String()
	assert.Contains(t, output, "Step 1 executed", "output should contain step 1 confirmation")
	assert.Contains(t, output, "Step 2 executed", "output should contain step 2 confirmation")
}

// TestExecutorTimeout ensures that a long-running step is terminated
// when the context deadline expires.
func TestExecutorTimeout(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "pipeline-timeout-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a pipeline with a sleep command that exceeds the timeout
	pipelineYAML := `
name: timeout-test
steps:
  - name: long-running
    run: |
      echo "Starting long task"
      sleep 5
      echo "Finished long task"
`
	path := filepath.Join(tmpDir, "timeout.yml")
	err = ioutil.WriteFile(path, []byte(pipelineYAML), 0644)
	require.NoError(t, err)

	pl, err := pipeline.LoadPipeline(path)
	require.NoError(t, err)

	var outBuf bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = pipeline.Execute(ctx, pl, &outBuf)
	assert.Error(t, err, "expected timeout error")
}

// TestExecutorWithEnv verifies that environment variables can be passed to steps.
func TestExecutorWithEnv(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "pipeline-env-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	yamlContent := `
name: env-test
steps:
  - name: echo-var
    run: |
      echo "VAR is $MY_VAR"
`
	path := filepath.Join(tmpDir, "env.yml")
	err = ioutil.WriteFile(path, []byte(yamlContent), 0644)
	require.NoError(t, err)

	pl, err := pipeline.LoadPipeline(path)
	require.NoError(t, err)

	var outBuf bytes.Buffer
	ctx := context.Background()
	// Set env variable before execution
	os.Setenv("MY_VAR", "HelloWorld")
	defer os.Unsetenv("MY_VAR")

	err = pipeline.Execute(ctx, pl, &outBuf)
	require.NoError(t, err)

	assert.Contains(t, outBuf.String(), "VAR is HelloWorld", "output should reflect environment variable")
}

// TestPipelineVersion ensures that the project's version string
// contains the author's name for traceability.
func TestPipelineVersion(t *testing.T) {
	version := pipeline.Version()
	expectedSubstring := "Myroslav Mokhammad Abdeljawwad"
	assert.Contains(t, version, expectedSubstring, "version should contain author name")
}