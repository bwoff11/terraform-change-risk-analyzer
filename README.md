# tf-risk-report

`tf-risk-report` is a small Go CLI that reads a Terraform plan JSON file (from `terraform show -json`), scores the risk of each resource change, and prints a human-readable report suitable for CI/CD pipelines and pull request review.

## Installation

Build from source:

```bash
go install github.com/gwoff/terraform-change-risk-analyzer/cmd/tf-risk-report@latest
```

## Usage

Generate a Terraform plan and pipe it into `tf-risk-report`:

```bash
terraform plan -out=plan.out
terraform show -json plan.out | tf-risk-report
```

Or point directly at a saved plan JSON:

```bash
tf-risk-report --plan-file=plan.json
```

Key flags:

- `--plan-file` (default `-`): Path to Terraform plan JSON (`-` reads from stdin).
- `--top-n` (default `20`): Maximum number of highest-risk changes to display.
- `--max-allowed-risk` (default `high`): Maximum allowed risk level before a non-zero exit code is returned (`none`, `low`, `medium`, `high`, `critical`).
- `--version`: Print the CLI version and exit.

## Exit codes

- `0` when the highest risk level is less than or equal to `--max-allowed-risk`.
- `1` when any change exceeds `--max-allowed-risk`.

This makes it easy to gate deployments on infrastructure risk.

## Development

To run tests and exercise the CLI against the bundled example plans on Windows (PowerShell):

```powershell
.\scripts\dev.ps1
```

This will:

- Run `go test ./...`.
- Run `tf-risk-report` against `testdata/plan_simple.json`.
- Run `tf-risk-report` against `testdata/plan_high_risk.json` with `--fail-threshold=MEDIUM`.

