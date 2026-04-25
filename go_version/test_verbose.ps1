param([string]$goal="exit")
$env:PROMETHEUS_VERBOSE="1"
Set-Location "C:\Users\junio\OneDrive\AI AGENT HACKATON\Prometheus\go_version"
.\bin\prometheus.exe $goal 2>&1