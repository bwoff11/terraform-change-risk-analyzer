$ErrorActionPreference = "Stop"

Write-Host "Running go test ./..."
go test ./...

if ($LASTEXITCODE -ne 0) {
    Write-Error "go test failed with exit code $LASTEXITCODE"
    exit $LASTEXITCODE
}

Write-Host ""
Write-Host "Running tf-risk-report against simple plan..."
go run ./cmd/tf-risk-report --plan-file=testdata/plan_simple.json

Write-Host ""
Write-Host "Running tf-risk-report against high-risk plan (fail-threshold=MEDIUM)..."
go run ./cmd/tf-risk-report --plan-file=testdata/plan_high_risk.json --fail-threshold=MEDIUM

