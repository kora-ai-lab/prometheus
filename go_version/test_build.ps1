$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -o /tmp/prometheus-test ./cmd/prometheus
if ($LASTEXITCODE -eq 0) {
    Write-Host "T5: Linux ARM64 build PASSED"
} else {
    Write-Host "T5: Linux ARM64 build FAILED"
}