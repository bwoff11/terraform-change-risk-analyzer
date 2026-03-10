package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Options holds CLI flag values after parsing and validation.
type Options struct {
	PlanFile       string
	TopN           int
	MaxAllowedRisk string
	ShowVersion    bool
}

const (
	defaultTopN           = 20
	defaultMaxAllowedRisk = "high"
)

// version is set at build time via -ldflags when possible.
var version = "dev"

// Parse parses command-line arguments into Options.
func Parse(args []string) (Options, error) {
	fs := flag.NewFlagSet("tf-risk-report", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var opts Options
	fs.StringVar(&opts.PlanFile, "plan-file", "-", "Path to Terraform plan JSON file ('-' for stdin)")
	fs.IntVar(&opts.TopN, "top-n", defaultTopN, "Maximum number of highest-risk changes to display")
	fs.StringVar(&opts.MaxAllowedRisk, "max-allowed-risk", defaultMaxAllowedRisk, "Maximum allowed risk level before non-zero exit code (none, low, medium, high, critical)")
	fs.BoolVar(&opts.ShowVersion, "version", false, "Print version and exit")

	if err := fs.Parse(args); err != nil {
		return Options{}, err
	}

	if opts.TopN <= 0 {
		return Options{}, fmt.Errorf("top-n must be positive")
	}

	return opts, nil
}

// OpenPlanSource returns an io.Reader for the Terraform plan JSON and a closer if needed.
func OpenPlanSource(opts Options) (io.Reader, io.Closer, error) {
	if opts.PlanFile == "" || opts.PlanFile == "-" {
		return os.Stdin, nil, nil
	}

	f, err := os.Open(opts.PlanFile)
	if err != nil {
		return nil, nil, err
	}

	return f, f, nil
}

// PrintVersion writes the CLI version to stdout.
func PrintVersion() {
	fmt.Fprintln(os.Stdout, version)
}

