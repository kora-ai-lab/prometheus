#!/usr/bin/env pwsh
# Prometheus Windows Installer
# Usage: irm https://raw.githubusercontent.com/kora-ai-lab/prometheus/main/scripts/install.ps1 | iex

$ErrorActionPreference = "Stop"

# Configuration
$Owner = "kora-ai-lab"
$Repo = "prometheus"
$InstallDir = $env:INSTALL_DIR
if (-not $InstallDir) {
    $InstallDir = "$env:LOCALAPPDATA\Programs"
}

# Detect architecture
$Arch = "amd64"
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
    $Arch = "arm64"
}

$BinaryName = "prometheus-windows-$Arch.exe"
$InstallPath = Join-Path $InstallDir "prometheus.exe"

Write-Host "Detecting platform... Windows/$Arch" -ForegroundColor Cyan

# Create install directory if needed
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
}

# Get latest release URL
$ApiUrl = "https://api.github.com/repos/$Owner/$Repo/releases/latest"
Write-Host "Fetching latest release..." -ForegroundColor Cyan

try {
    $Release = Invoke-RestMethod -Uri $ApiUrl -Headers @{ "Accept" = "application/vnd.github+json" }
    $Asset = $Release.assets | Where-Object { $_.name -eq $BinaryName }
    
    if (-not $Asset) {
        Write-Error "Could not find binary $BinaryName in latest release"
        exit 1
    }
    
    $DownloadUrl = $Asset.browser_download_url
    Write-Host "Downloading from $DownloadUrl..." -ForegroundColor Cyan
    
    # Download binary
    $TempPath = Join-Path $env:TEMP "prometheus.exe"
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempPath -UseBasicParsing
    
    # Verify checksum if available
    $ChecksumAsset = $Release.assets | Where-Object { $_.name -eq "$BinaryName.sha256" }
    if ($ChecksumAsset) {
        $ExpectedHash = (Invoke-WebRequest -Uri $ChecksumAsset.browser_download_url -UseBasicParsing).Content.Trim().Split()[0]
        $ActualHash = (Get-FileHash -Path $TempPath -Algorithm SHA256).Hash.ToLower()
        
        if ($ExpectedHash -ne $ActualHash) {
            Write-Error "Checksum mismatch! Expected: $ExpectedHash, Got: $ActualHash"
            Remove-Item $TempPath -Force
            exit 1
        }
        Write-Host "Checksum verified ✓" -ForegroundColor Green
    }
    
    # Install binary
    Move-Item -Path $TempPath -Destination $InstallPath -Force
    Write-Host "Installed to $InstallPath" -ForegroundColor Green
    
    # Add to PATH if not already there
    $CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($CurrentPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$InstallDir", "User")
        Write-Host "Added $InstallDir to PATH (restart terminal to use)" -ForegroundColor Yellow
    }
    
    Write-Host "`nPrometheus installed successfully!" -ForegroundColor Green
    Write-Host "Run 'prometheus --help' to get started" -ForegroundColor Cyan
    
} catch {
    Write-Error "Installation failed: $_"
    exit 1
}
