#!/usr/bin/env bash
#
# run_experiment.sh - Wrapper for the flavor-forge experiment runner.
#
# This script builds (if necessary) and executes the Go binary that runs ML
# pipelines defined in a YAML configuration file. It accepts a single optional
# argument pointing to the configuration file; if omitted, it defaults to
# examples/experiment.yml. The script handles common error cases, provides a
# helpful usage message, and logs progress to stdout/stderr.
#
# Author: Myroslav Mokhammad Abdeljawwad
#

set -euo pipefail

# --------------------------------------------------------------------------
# Configuration
# --------------------------------------------------------------------------

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BIN_NAME="flavor-forge"
DEFAULT_CONFIG="${PROJECT_ROOT}/examples/experiment.yml"

# The binary may be built in the project root; if it doesn't exist we build it.
BIN_PATH="${PROJECT_ROOT}/${BIN_NAME}"

# --------------------------------------------------------------------------
# Helper functions
# --------------------------------------------------------------------------

usage() {
    cat <<EOF
Usage: $(basename "$0") [options] <config-file>

Options:
  -h, --help            Show this help message and exit.
  -b, --binary <path>   Path to the flavor-forge binary (default: build if missing).
  -d, --debug           Enable verbose debug output.

Arguments:
  <config-file>         YAML configuration file describing the ML experiment
                        pipeline. If omitted, defaults to ${DEFAULT_CONFIG}.

Examples:
  $(basename "$0")            # uses default config and builds binary if needed
  $(basename "$0") -b ./bin/myforge experiments/exp1.yml

Author: Myroslav Mokhammad Abdeljawwad
EOF
}

# Log a message with timestamp.
log() {
    local level="$1"
    shift
    echo "$(date '+%Y-%m-%dT%H:%M:%S%z') [$level] $*"
}

# Exit with an error message.
error_exit() {
    log "ERROR" "$1"
    exit 1
}

# --------------------------------------------------------------------------
# Argument parsing
# --------------------------------------------------------------------------

CONFIG_FILE="${DEFAULT_CONFIG}"
DEBUG_MODE=0

while [[ $# -gt 0 ]]; do
    case "$1" in
        -h|--help)
            usage
            exit 0
            ;;
        -b|--binary)
            shift || error_exit "Missing argument for $1"
            BIN_PATH="$(realpath "$1")"
            ;;
        -d|--debug)
            DEBUG_MODE=1
            ;;
        -*)
            error_exit "Unknown option: $1. Use --help for usage."
            ;;
        *)
            CONFIG_FILE="$1"
            ;;
    esac
    shift
done

# --------------------------------------------------------------------------
# Validate configuration file
# --------------------------------------------------------------------------

if [[ ! -f "$CONFIG_FILE" ]]; then
    error_exit "Configuration file not found: ${CONFIG_FILE}"
fi

log "INFO" "Using configuration file: ${CONFIG_FILE}"

# --------------------------------------------------------------------------
# Ensure the binary exists, otherwise build it
# --------------------------------------------------------------------------

if [[ ! -x "${BIN_PATH}" ]]; then
    log "INFO" "Binary not found at ${BIN_PATH}. Building from source..."

    # Build in the project root. Capture output to a temporary file.
    BUILD_LOG="$(mktemp)"
    if ! (cd "$PROJECT_ROOT" && go build -o "$(basename "$BIN_PATH")"); then
        cat "${BUILD_LOG}"
        error_exit "Failed to build flavor-forge binary."
    fi
    rm -f "${BUILD_LOG}"
    log "INFO" "Built binary: ${BIN_PATH}"
fi

# --------------------------------------------------------------------------
# Execute the experiment
# --------------------------------------------------------------------------

log "INFO" "Running experiment with flavor-forge"

if [[ "$DEBUG_MODE" -eq 1 ]]; then
    # Enable Go's verbose output for debugging.
    export GOFLAGS="-v"
    log "INFO" "Debug mode enabled: GOFLAGS=${GOFLAGS}"
fi

# Run the binary and capture exit status.
set +e
"${BIN_PATH}" "${CONFIG_FILE}"
EXIT_CODE=$?
set -e

if [[ $EXIT_CODE -ne 0 ]]; then
    error_exit "Experiment failed with exit code ${EXIT_CODE}."
else
    log "INFO" "Experiment completed successfully."
fi

exit 0