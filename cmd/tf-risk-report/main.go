package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gwoff/terraform-change-risk-analyzer/internal/plan"
	"github.com/gwoff/terraform-change-risk-analyzer/internal/report"
	"github.com/gwoff/terraform-change-risk-analyzer/internal/risk"
)

func main() {
	var (
		planFile      string
		failThreshold string
		topN          int
	)

	flag.StringVar(&planFile, "plan-file", "-", "Path to Terraform plan JSON file ('-' for stdin)")
	flag.StringVar(&failThreshold, "fail-threshold", "HIGH", "Severity threshold (LOW, MEDIUM, HIGH, CRITICAL) at or above which the command exits non-zero")
	flag.IntVar(&topN, "top-n", 20, "Maximum number of highest-risk resources to display")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "tf-risk-report: analyze Terraform plan risk.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  terraform show -json plan.out | tf-risk-report [flags]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  tf-risk-report --plan-file=plan.json [flags]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	var (
		r   *os.File
		err error
	)

	if planFile == "-" {
		r = os.Stdin
	} else {
		r, err = os.Open(planFile)
		if err != nil {
			log.Fatalf("failed to open plan file %q: %v", planFile, err)
		}
		defer r.Close()
	}

	pl, err := plan.Parse(r)
	if err != nil {
		log.Fatalf("failed to parse Terraform plan JSON: %v", err)
	}

	changeRisks, summary := risk.ScoreChanges(pl)

	if err := report.Render(os.Stdout, changeRisks, summary, topN); err != nil {
		log.Fatalf("failed to render risk report: %v", err)
	}

	exitCode := report.DecideExitCode(summary, strings.ToLower(failThreshold))
	os.Exit(exitCode)
}

