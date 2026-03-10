package main

import (
	"log"
	"os"

	"github.com/gwoff/terraform-change-risk-analyzer/internal/cli"
	"github.com/gwoff/terraform-change-risk-analyzer/internal/plan"
	"github.com/gwoff/terraform-change-risk-analyzer/internal/report"
	"github.com/gwoff/terraform-change-risk-analyzer/internal/risk"
)

func main() {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("failed to parse arguments: %v", err)
	}

	if opts.ShowVersion {
		cli.PrintVersion()
		return
	}

	planReader, closer, err := cli.OpenPlanSource(opts)
	if err != nil {
		log.Fatalf("failed to open plan file: %v", err)
	}
	if closer != nil {
		defer closer.Close()
	}

	pl, err := plan.Parse(planReader)
	if err != nil {
		log.Fatalf("failed to parse Terraform plan JSON: %v", err)
	}

	changeRisks, summary := risk.ScoreChanges(pl)

	if err := report.Render(os.Stdout, changeRisks, summary, opts.TopN); err != nil {
		log.Fatalf("failed to render risk report: %v", err)
	}

	exitCode := report.DecideExitCode(summary, opts.MaxAllowedRisk)
	os.Exit(exitCode)
}

