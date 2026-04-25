#!/usr/bin/env pwsh
# Build Prometheus MSI installer for Windows
# Usage: .\build-msi.ps1

$ErrorActionPreference = "Stop"

# Configuration
$Version = "1.0.5"
$SourceDir = "..\..\release"
$WxsFile = "prometheus.wxs"
$OutputDir = "..\..\release"

# Check if WiX is installed
$WixPath = "${env:ProgramFiles(x86)}\WiX Toolset v3.11\bin"
if (-not (Test-Path $WixPath)) {
    Write-Error "WiX Toolset not found. Install from: https://wixtoolset.org/releases/"
    exit 1
}

# Add WiX to PATH
$env:PATH = "$WixPath;$env:PATH"

# Check if binary exists
$BinaryPath = Join-Path $SourceDir "prometheus-windows-amd64.exe"
if (-not (Test-Path $BinaryPath)) {
    Write-Error "Binary not found at $BinaryPath. Build it first with build-release.sh"
    exit 1
}

# Create output directory
if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir -Force | Out-Null
}

# Copy binary to current directory for WiX
Copy-Item $BinaryPath -Destination "." -Force

# Generate GUID for UpgradeCode (static for upgrades)
$UpgradeCode = "A1B2C3D4-E5F6-7890-ABCD-EF1234567890"

# Update version in wxs file
$WxsContent = Get-Content $WxsFile -Raw
$WxsContent = $WxsContent -replace 'Version="1\.0\.5\.0"', "Version=`"$Version.0`""
$WxsContent = $WxsContent -replace 'UpgradeCode="YOUR-GUID-HERE"', "UpgradeCode=`"$UpgradeCode`""
Set-Content $WxsFile -Value $WxsContent

# Build MSI
Write-Host "Building MSI..."
candle $WxsFile -out "prometheus.wixobj" -ext WixUIExtension
if ($LASTEXITCODE -ne 0) {
    Write-Error "candle failed"
    exit 1
}

light prometheus.wixobj -out (Join-Path $OutputDir "prometheus-windows-amd64.msi") -ext WixUIExtension
if ($LASTEXITCODE -ne 0) {
    Write-Error "light failed"
    exit 1
}

# Cleanup
Remove-Item "prometheus-windows-amd64.exe" -Force -ErrorAction SilentlyContinue
Remove-Item "prometheus.wixobj" -Force -ErrorAction SilentlyContinue

Write-Host "MSI built successfully: $OutputDir\prometheus-windows-amd64.msi" -ForegroundColor Green
